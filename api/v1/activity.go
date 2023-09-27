package v1

import (
	"maisonsport/dao"
	"maisonsport/log"
	"maisonsport/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CreateActivity(c *gin.Context) {

	user_id := c.GetString("user_id") // 先取出 user_id，用 useri_id 去证明身份
	var data models.ActivityInfo
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Logger.Error("CreateActivity 反序列化失败", zap.Any("err", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data.Creator = user_id // 把创建者 ID  带上

	activityID, err := dao.CreateActivity(&data)
	if err != nil {
		log.Logger.Error("dao.CreateActivity 出问题了", zap.Any("err", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"activity_id": activityID,
	})
}

type FrontendRequest struct {
	ID int `json:"id"` // 活动 ID
}

// 通过活动 ID获取互动的信息和活动的相关参与人
func GetActivityByActivityID(c *gin.Context) {
	var data FrontendRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Logger.Error("GetActivityByActivityID 反序列化失败", zap.Any("err", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	activityInfo, err := dao.GetActivityByActivityID(data.ID) // 先获取活动的本身信息，然后下面获取参与的人
	if err != nil {
		log.Logger.Error("dao.GetActivityByActivityID 出问题了", zap.Any("err", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	memberInfoList, err := dao.GetActivityMembersByActivityID(data.ID) //获取相关人员
	if err != nil {
		log.Logger.Error("dao.GetActivityMembersByActivityID 出问题了", zap.Any("err", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	memberList := []MemberInfo{}
	for _, info := range memberInfoList {

		// 要根据 info  里面的  user_id 去查出每一个人的 userInfo
		thisUserInfo, _ := dao.GetUserInfoByUserID(info.UserID)

		thisMemberInfo := MemberInfo{
			Status:   info.Status,
			UserInfo: thisUserInfo,
		}
		memberList = append(memberList, thisMemberInfo)
	}

	c.JSON(http.StatusOK, gin.H{
		"title":        activityInfo.Title,
		"tar_gender":   activityInfo.TarGender,
		"tar_level":    activityInfo.TarLevel, //float,
		"time":         activityInfo.Time,     //string
		"note":         activityInfo.Note,
		"max_member":   activityInfo.MaxMember,
		"creator":      activityInfo.Creator,
		"memeber_list": memberList,
	})

}

type MemberInfo struct {
	Status   int             `json:"status"` // 成员状态
	UserInfo models.UserInfo `json:"user_info"`
}

// =================================================================

type ApplyRequest struct {
	ID int `json:"id"` // 活动 ID
}

// 申请加入某个活动
func ApplyActivity(c *gin.Context) {
	user_id := c.GetString("user_id") // 先取出 user_id，用 useri_id 去证明身份

	var data ApplyRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Logger.Error("ApplyActivity 反序列化失败", zap.Any("err", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	activityApply := models.ActivityMember{
		ActivityID: data.ID,
		UserID:     user_id,
		Status:     0, // 申请是审核中 - 0，其他状态为 成功 -  1， 被拒绝-  2
	}

	err := dao.ApplyActivity(&activityApply)
	if err != nil {
		log.Logger.Error("dao.ApplyActivity 出问题了", zap.Any("err", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{})

}

// 获取关于自己的所有活动详情
func GetMyActivities(c *gin.Context) {
	user_id := c.GetString("user_id") // 先取出 user_id，用 useri_id 去证明身份

	myActivity, othersActivity, err := dao.GetMyActivities(user_id)
	if err != nil {
		log.Logger.Error("dao.GetMyActivities 出问题了", zap.Any("err", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"mine":   myActivity,
		"others": othersActivity,
	})

}

// ********************************
// 核心函数， 广场展示活动信息 ，按照过滤条件
// ********************************
func FilterActivity(c *gin.Context) {

	var filter models.Filter
	if err := c.ShouldBindJSON(&filter); err != nil {
		log.Logger.Error("ApplyActivity 反序列化失败", zap.Any("err", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	activity_list, err := dao.FilterActivity(filter) // 查出所有符合条件的活动
	if err != nil {
		log.Logger.Error("dao.FilterActivity 出问题了", zap.Any("err", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// 接下来要查出每个活动的  creator  和参与者的 list，把他们的  userinfo 全都丢进去 ,还有距离
	filterActivityInfoList, err := dao.GetFilterActivityAllInfo(activity_list, filter.PositionX, filter.PositionY)
	if err != nil {
		log.Logger.Error("dao.GetFilterActivityAllInfo 出问题了", zap.Any("err", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"activity_list": filterActivityInfoList,
	})

}

type ApplyParam struct {
	ActivityID string `json:"activity_id"`
	UserID     string `json:"user_id"`
	Update     int    `json:"update"`
}

// 更新请求活动的状态
func ApplyUpdate(c *gin.Context) {
	reqParams := &ApplyParam{}
	if err := c.ShouldBindJSON(&reqParams); err != nil {
		log.Logger.Error("ApplyUpdate 反序列化失败", zap.Any("err", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := dao.ApplyUpdate(reqParams.ActivityID, reqParams.UserID, reqParams.Update)
	if err != nil {
		log.Logger.Error("dao.ApplyUpdate 出问题了", zap.Any("err", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{})
}
