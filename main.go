package main

import (
	v1 "maisonsport/api/v1"
	"maisonsport/dao"
	"maisonsport/log"
	"maisonsport/middleware"
	"maisonsport/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	log.InitZap() // 先初始化日志系统
	log.Logger.Info("Zap日志框架初始化成功")
	dao.InitViper()
	log.Logger.Info("Viper日志框架初始化成功")
	dao.InitDB() //然后初始化 DB
	log.Logger.Info("DB初始化成功")
	dao.InitRedis()
	log.Logger.Info("Redis初始化成功")

	r := gin.Default()

	r.POST("/logIn", v1.SilentLogIn) //这个是自动登录时候调用的，不要加鉴权
	r.POST("/upload", v1.HandleVideoUpload)
	userInfoV1 := r.Group("/v1")
	{
		userInfoV1.Use(middleware.AuthMiddleware()) //鉴权用的中间件
		routes.SetupUserInfoRoutesV1(userInfoV1)
		routes.SetupActivityRoutesV1(userInfoV1)

	}

	r.Run(":8080")
}
