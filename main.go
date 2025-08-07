package main

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/dwburke/weather/types"
)

func init() {
	// Initialize the database connection and register callbacks

	viper.SetDefault("db.user", "addict")
	viper.SetDefault("db.name", "test")
	viper.SetDefault("db.pass", "")
}

func main() {

	t := &types.Test{}
	t.Amount = 12.34
	t.DateTime = "2024-10-01 12:34:56"
	if err := t.Create(); err != nil {
		fmt.Println("Error creating Test record:", err)
		return
	}
}
