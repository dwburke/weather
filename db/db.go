package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"

	"github.com/dwburke/weather/db/validate"
)

type MyDb struct {
	conn *gorm.DB
}

func init() {
	viper.SetDefault("db.maxidleconnections", 2)
	viper.SetDefault("db.maxopenconnections", 12)
	viper.SetDefault("db.sslmode", "disable")
	viper.SetDefault("db.connect_timeout", 90)
	viper.SetDefault("db.port", 5432)
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

	connStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s connect_timeout=%d",
		viper.GetString("db.host"),
		viper.GetInt("db.port"),
		viper.GetString("db.user"),
		viper.GetString("db.name"),
		viper.GetString("db.pass"),
		viper.GetString("db.sslmode"),
		viper.GetInt("db.connect_timeout"),
	)

	conn, err := gorm.Open("postgres", connStr)
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
