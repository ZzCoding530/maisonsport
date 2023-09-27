package dao

import (
	"log"

	"github.com/spf13/viper"
)

var Viper *viper.Viper

func InitViper() {
	// 初始化 Viper
	v := viper.New()
	v.SetConfigName("config") // 配置文件名 (without extension)
	v.SetConfigType("yaml")   // 配置文件类型
	v.AddConfigPath(".")      // 配置文件路径（当前目录）

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	Viper = v
}
