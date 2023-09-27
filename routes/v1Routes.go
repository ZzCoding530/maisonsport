package routes

import (
	v1 "maisonsport/api/v1"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 填写用户资料的路由
func SetupUserInfoRoutesV1(r *gin.RouterGroup) {
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "看到这个说明 V1版本的 userInfo相关的路由都没问题")
	})
	r.POST("/updateUserInfo", v1.UpdateUserInfo)

	r.POST("/getUserInfo", v1.GetUserInfo)

}

// 活动相关的路由
func SetupActivityRoutesV1(r *gin.RouterGroup) {
	r.GET("/activity", func(c *gin.Context) {
		c.String(http.StatusOK, "看到这个说明 V1版本的 activity接口都没问题")
	})
	r.POST("/activity/create", v1.CreateActivity)
	r.POST("/activity/getInfo", v1.GetActivityByActivityID)
	r.POST("/activity/apply,", v1.ApplyActivity)
	r.POST("/activity/aboutMe", v1.GetMyActivities)
	r.POST("/activity/filter", v1.FilterActivity)
	r.POST("/activity/applyUpdate", v1.ApplyUpdate)

	r.POST("/getTempUrl", v1.GetTempUrl)

}
