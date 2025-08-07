package types

import (
	"github.com/dwburke/weather/db"
)

type Test struct {
	Amount   float32 `json:"amount" gorm:"column:amount" mapstructure:"amount"`
	DateTime string  `json:"date_time" gorm:"column:date_time" mapstructure:"date_time"`
}

func (Test) TableName() string {
	return "test"
}

func (self *Test) Create() error {
	gdbh, err := db.GetDB().DB()
	if err != nil {
		return err
	}

	if err := gdbh.Create(&self).Error; err != nil {
		return err
	}

	return nil
}

func (self *Test) Save() error {
	gdbh, err := db.GetDB().DB()
	if err != nil {
		return err
	}

	if err := gdbh.Save(&self).Error; err != nil {
		return err
	}

	return nil
}

func (self *Test) Delete() error {
	if gdbh, err := db.GetDB().DB(); err != nil {
		return err
	} else {
		return gdbh.Delete(&self).Error
	}
}

/*
func AccntInfoFind(uniq_id string) (*Test, error) {

	//if err := validation.CheckOrError("uniq_id|uuid|resize_uniq_id", uniq_id); err != nil {
	//return nil, err
	//}

	dbh, err := db.GetDB().DB()
	if err != nil {
		return nil, err
	}

	var ai Test

	if err := dbh.Where("uniq_id = ?", uniq_id).Find(&ai).Error; err != nil {
		return nil, err
	}
	return &ai, nil
}
*/
