package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"

	"github.com/dwburke/weather/db/validate"
)

type MyDb struct {
	conn *gorm.DB
}

func init() {
	viper.SetDefault("db.maxidleconnections", 2)
	viper.SetDefault("db.maxopenconnections", 12)
	viper.SetDefault("db.connect_timeout", 90)
	viper.SetDefault("db.port", 3306)
	viper.SetDefault("db.user", "")
	viper.SetDefault("db.name", "")
	viper.SetDefault("db.pass", "")

}

func NewDB() *MyDb {
	db := &MyDb{}
	return db
}

// get gorm db handle for default context
func (db *MyDb) DB() (*gorm.DB, error) {
	return db.dbh()
}

func (db *MyDb) dbh() (*gorm.DB, error) {
	if db.conn != nil {
		return db.conn, nil
	}

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%ds",
		viper.GetString("db.user"),
		viper.GetString("db.pass"),
		viper.GetString("db.host"),
		viper.GetInt("db.port"),
		viper.GetString("db.name"),
		viper.GetInt("db.connect_timeout"),
	)

	conn, err := gorm.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}

	conn = conn.Set("gorm:auto_preload", true)

	conn.DB().SetMaxIdleConns(viper.GetInt("db.maxidleconnections"))
	conn.DB().SetMaxOpenConns(viper.GetInt("db.maxopenconnections"))

	validate.RegisterCallbacks(conn)

	db.conn = conn

	return db.conn, nil
}
