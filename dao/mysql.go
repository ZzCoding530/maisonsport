package dao

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() {
	// 读取 MySQL 配置

	mysqlHost := Viper.GetString("mysql.host")
	mysqlPort := Viper.GetInt("mysql.port")
	mysqlUsername := Viper.GetString("mysql.username")
	mysqlPassword := Viper.GetString("mysql.password")
	mysqlDBName := Viper.GetString("mysql.dbname")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysqlUsername, mysqlPassword, mysqlHost, mysqlPort, mysqlDBName)

	fmt.Printf("dsn: %v\n", dsn)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

}
