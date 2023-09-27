package v1

import (
	"maisonsport/log"
	"maisonsport/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 接收从前端发来的 code  用的
type ReqCos struct {
	FileName string `json:"fileName"`
}

func GetTempUrl(c *gin.Context) {

	var req ReqCos
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Logger.Error("GetTempUrl 反序列化失败")
		c.JSON(http.StatusBadRequest, gin.H{"error": "反序列化失败"})
		return
	}
	user_id := c.GetString("user_id") // 先取出 user_id，用 useri_id 去证明身份

	temp_url, name, public_url := utils.GetTencentPreSignedUrl(user_id, req.FileName)

	c.JSON(http.StatusOK, gin.H{
		"temp_url":   temp_url,
		"name":       name,
		"public_url": public_url,
	})

}
