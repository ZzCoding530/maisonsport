package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"maisonsport/dao"
	"maisonsport/log"
	"maisonsport/models"
	"maisonsport/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 接收从前端发来的 code  用的
type loginReq struct {
	Code string `json:"code"`
}

// 用户登录小程序初次，自动静默登录，
// 1.要返回给前端  token
// 2.同时生成一个 user_id
// 3.然后往 mysql 里面的 t_user_info 存入 user_id 和 token
// 4. Redis 里面也要存， key 为  token ， value 为 user_id
func SilentLogIn(c *gin.Context) {

	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Logger.Error("UserLoginHandler 反序列化失败")
		c.JSON(http.StatusBadRequest, gin.H{"error": "反序列化失败"})
		return
	} //拿到了微信小程序发来的 code。要向微信接口换取 openid  和 sessionkey
	log.Logger.Info("收到了前端发来的code", zap.Any("code", req.Code))

	// 然后向微信端要 openID 和 sessionkey
	loginResponse, err := getOpenIDAndSessionKey(req.Code)
	if err != nil {
		log.Logger.Error("getOpenIDAndSessionKey执行出错", zap.Any("Error", err))
		return
	}

	if loginResponse.ErrCode != 0 {
		log.Logger.Info("微信公开接口返回了错误", zap.Any("ErrorCode", loginResponse.ErrCode), zap.Any("ErrMsg", loginResponse.ErrMsg))
		fmt.Printf("微信公开接口返回了错误，Error: %d - %s\n", loginResponse.ErrCode, loginResponse.ErrMsg)
		return
	}

	//如果没错，那就是拿到了两个东西
	fmt.Println("成功拿到 OpenID 和 sessionKey！", "|||", loginResponse.OpenID, "|||", loginResponse.SessionKey)
	log.Logger.Info("成功拿到 OpenID 和 sessionKey！", zap.Any("loginResponse", loginResponse))

	token, err := utils.GenerateToken(loginResponse.OpenID, loginResponse.SessionKey) // 用 openID 和 sessionkey 生成 token
	if err != nil {
		log.Logger.Error("GenerateToken 出问题")
	}
	user_id, err := utils.GenerateUUID(loginResponse.OpenID) // 生成唯一标识符 user_id
	if err != nil {
		log.Logger.Error("GenerateUUID 出问题")
	}

	err = dao.InsertUserInfoPart(user_id, loginResponse.OpenID, token) // 把  userid token openid 都存mysql
	if err != nil {
		log.Logger.Error("InsertUserInfoPart 出问题", zap.Any("err", err))
	}

	dao.RedisDB.Set(token, user_id, 0) // 存入 Redis ，kV 为  token  和  user_id
	if err != nil {
		log.Logger.Error("向 Redis 里存token 和 userid 出问题了", zap.Any("err", err))
	}

	// f返回给前端 token
	c.JSON(http.StatusOK, gin.H{
		"msg":        "新注册成功",
		"code":       req.Code,
		"OpenID":     loginResponse.OpenID,
		"SessionKey": loginResponse.SessionKey,
		"token":      token,
		"user_id":    user_id,
	})

}

// =================================================================
// 接收 小程序公开接口返回  openID  和 sessionKey 用的
type WechatLoginResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// 向小程序发  code   返回 openID  和 sessionID
func getOpenIDAndSessionKey(code string) (*WechatLoginResponse, error) {
	appID := "wx0a92d5ee853f7654"
	appSecret := "5f48fd688e4ddf05938c0778872f6e58"

	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", appID, appSecret, code)

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var loginResponse WechatLoginResponse
	err = json.Unmarshal(body, &loginResponse)
	if err != nil {
		return nil, err
	}

	return &loginResponse, nil
}

// ================================================================
// ================================================================
// 通过 token 获取个人资料审核状态
func GetUserInfoStatusByToken(c *gin.Context) {
	user_id := c.GetString("user_id") // 先取出 user_id，用 useri_id 去证明身份
	fmt.Println("这里取出的 userid 为：", user_id)

	userInfoStatus, err := dao.GetUserInfoStatusByUserID(user_id)
	if err != nil {
		log.Logger.Error("GetUserInfoStatusByUserID 出问题", zap.Any("err", err))
	}

	c.JSON(http.StatusOK, gin.H{
		"userInfoStatus": userInfoStatus,
	})
}

// ================================================================
type FrontendData struct {
	Nickname     string  `json:"nickname"`
	Gender       string  `json:"gender"`
	Age          int     `json:"age"`
	Level        float64 `json:"level"`
	ProvinceName string  `json:"province_name"`
	CityName     string  `json:"city_name"`
	DistrictName string  `json:"district_name"`
	AvatarURL    string  `json:"avatar_url"`
	VideoURL     string  `json:"video_url"`
	CellPhone    string  `json:"cell_phone"`
}

// 修改个人资料表
func UpdateUserInfo(c *gin.Context) {
	var data FrontendData
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Logger.Error("UserLoginHandler 反序列化失败", zap.Any("err", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	showReq(data) //	打印展示一下

	userInfo := transType(data)       //转换一下类型，然后存入 mysql
	user_id := c.GetString("user_id") // 先取出 user_id，用 useri_id 去证明身份
	dao.UpdateUserInfoByUserID(user_id, userInfo)

	c.JSON(http.StatusOK, gin.H{
		"userInfo": userInfo,
	})

}

// 打印工具函数
func showReq(req FrontendData) {
	fmt.Println("后端接收到了内容为：===========")
	fmt.Printf("req.Nickname: %v\n", req.Nickname)
	fmt.Printf("req.Gender: %v\n", req.Gender)
	fmt.Printf("req.Level: %v\n", req.Level)
	fmt.Printf("req.CityName: %v\n", req.CityName)
	fmt.Printf("req.Age: %v\n", req.Age)
	fmt.Println("==============================")
}

// 转换一些字段的类型
func transType(req FrontendData) models.UserInfo {

	userInfo := models.UserInfo{
		NickName:       req.Nickname,
		Gender:         req.Gender,
		Age:            req.Age,
		SkillLevel:     req.Level,
		ProvinceName:   req.ProvinceName,
		CityName:       req.CityName,
		DistrictName:   req.DistrictName,
		AvatarURL:      req.AvatarURL,
		VideoUrl:       req.VideoURL,
		CellPhone:      req.CellPhone,
		IsInfoVerified: 1,
	}
	return userInfo
}

// 获取个人所有基础信息根据 userid
func GetUserInfo(c *gin.Context) {
	user_id := c.GetString("user_id") // 先取出 user_id，用 useri_id 去证明身份
	userInfo, err := dao.GetUserInfoByUserID(user_id)
	if err != nil {
		log.Logger.Error("GetUserInfoByUserID 出问题了", zap.Any("error", err))
	}

	c.JSON(http.StatusOK, gin.H{
		"userInfo": userInfo,
	})
}

func HandleVideoUpload(c *gin.Context) {
	videoURL := c.PostForm("videoUrl")
	// 在这里可以将 videoURL 存入数据库或进行其他操作
	fmt.Printf("Received video URL: %s\n", videoURL)
	c.String(http.StatusOK, "Video URL received and processed successfully")
}
