package dao

import (
	"maisonsport/log"
	"maisonsport/models"
	"math"
	"time"

	"go.uber.org/zap"
)

// 插入新的活动数据，返回插入后的 那个活动ID
func CreateActivity(activity *models.ActivityInfo) (int, error) {
	// Insert the new activity
	result := db.Create(activity)
	if result.Error != nil {
		return 0, result.Error
	}

	return activity.ID, nil
}

func GetActivityByActivityID(activityID int) (models.ActivityInfo, error) {
	var activityInfo models.ActivityInfo
	result := db.First(&activityInfo, activityID)
	if result.Error != nil {
		return models.ActivityInfo{}, result.Error
	}

	return activityInfo, nil
}

func ApplyActivity(activity *models.ActivityMember) error {
	// Insert the new activity member
	result := db.Create(activity)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetMyActivities(userId string) (myActivity []models.ActivityInfo, othersActivity []models.ActivityInfo, err error) {
	var activityMembers []models.ActivityMember

	// 从t_activity_member;表里取出所有相关的
	result := db.Where("user_id = ?", userId).Find(&activityMembers)
	if result.Error != nil {
		err = result.Error
		return nil, nil, err
	}

	// 拿出其中的 activity_id
	//先分出来自己的的活动 ID 和别人发起的 ID
	var myActivityIDs []int
	var othersActivityIDs []int

	for _, member := range activityMembers {
		if member.UserID == userId {
			myActivityIDs = append(myActivityIDs, member.ActivityID)
		} else {
			othersActivityIDs = append(othersActivityIDs, member.ActivityID)
		}
	}

	// 查询自己创建的活动
	result = db.Where("id IN (?)", myActivityIDs).Find(&myActivity)
	if result.Error != nil {
		err = result.Error
		return nil, nil, err
	}
	// 查询别人创建的活动
	result = db.Where("id IN (?)", othersActivityIDs).Find(&othersActivity)
	if result.Error != nil {
		err = result.Error
		return nil, nil, err
	}

	return myActivity, othersActivity, nil

}

// 更改申请的状态，通过还是被拒绝
func ApplyUpdate(ActivityID string, UserID string, Update int) error {
	err := db.Model(&models.ActivityMember{}).Where("activity_id =? and user_id =?", ActivityID, UserID).Update("status", Update).Error

	return err

}

// =================================================================

func FilterActivity(filter models.Filter) ([]models.ActivityInfo, error) {
	// Execute the query with the constructed conditions
	var activities []models.ActivityInfo

	if filter.City != "" { // 过滤城市
		db = db.Where("city", filter.City)
	}

	if filter.Gender != "" { // 过滤性别
		db = db.Where("tar_gender", filter.Gender)
	}

	if filter.Date != "" { // 过滤日期
		dateString := filter.Date
		// 解析日期字符串为 Time 类型
		dateTime, err := time.Parse("2006-01-02", dateString)
		if err != nil {
			log.Logger.Error("FilterActivity 时间解析错误")
		}
		db = db.Where("time = ?", dateTime)
	}

	if filter.LevelMax != 0 { // 过滤最高水平
		db = db.Where("tar_level", "<=", filter.LevelMax)
	}
	if filter.LevelMin != 0 { // 过滤最低水平
		db = db.Where("tar_level", ">=", filter.LevelMin)
	}

	if filter.MemberMax != 0 { //过滤最高人数
		db = db.Where("member_num", "<=", filter.MemberMax)
	}

	if filter.MemberMin != 0 { //过滤最低人数
		db = db.Where("member_num", ">=", filter.MemberMin)
	}

	offset := (filter.PageIndex - 1) * filter.Pagesize

	db = db.Limit(filter.Pagesize).Offset(offset) //加上分页查询

	err := db.Find(&activities).Error
	if err != nil {
		log.Logger.Error("FilterActivity 查询数据库出错", zap.Any("err", err))
	}
	return activities, nil
}

type FilterActivityInfo struct {
	ActivityInfo models.ActivityInfo // 原生的活动信息
	Distance     float64             // 距离
	CreatorInfo  models.UserInfo     // 创建者信息
	MemberInfo   []models.UserInfo   // 成员信息
}

// 拿到广场过滤的所有信息
func GetFilterActivityAllInfo(activity_list []models.ActivityInfo, usr_X, usr_Y float64) ([]FilterActivityInfo, error) {
	var filterActivityInfoList []FilterActivityInfo
	for _, activity := range activity_list {
		filterActivityInfo, err := getOneFilterActivityAllInfo(activity, usr_X, usr_Y)
		if err != nil {
			return filterActivityInfoList, err
		}
		filterActivityInfoList = append(filterActivityInfoList, filterActivityInfo)
	}

	return filterActivityInfoList, nil
}

// GetFilterActivityAllInfo 里面的处理单个activity的函数
func getOneFilterActivityAllInfo(activity models.ActivityInfo, usr_X, usr_Y float64) (FilterActivityInfo, error) {
	creatorInfo, err := GetUserInfoByUserID(activity.Creator)
	if err != nil {
		return FilterActivityInfo{}, err
	}

	memberIDs, err := GetUserIDListByActivityID(activity.ID) // 拿到所有申请通过的用户列表
	if err != nil {
		return FilterActivityInfo{}, err
	}

	membersInfoList, err := GetUserInfoByUserIDs(memberIDs) // 批量查出所有的申请通过的用户信息
	if err != nil {
		return FilterActivityInfo{}, err
	}

	distance := calculateDistance(usr_X, usr_Y, activity.PositionX, activity.PositionY) // 计算出这个活动的距离

	return FilterActivityInfo{ //把结果返回去
		ActivityInfo: activity,
		Distance:     distance,
		CreatorInfo:  creatorInfo,
		MemberInfo:   membersInfoList,
	}, nil

}

// =================================================================
// 计算两地距离

const (
	earthRadius = 6371.0 // 地球半径（单位：公里）
)

func toRadians(deg float64) float64 {
	return deg * (math.Pi / 180.0)
}

// calculateDistance 计算两地距离（单位：公里）
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// 将经度和纬度转换为弧度
	lat1 = toRadians(lat1)
	lon1 = toRadians(lon1)
	lat2 = toRadians(lat2)
	lon2 = toRadians(lon2)

	// Haversine formula
	dlon := lon2 - lon1
	dlat := lat2 - lat1
	a := math.Pow(math.Sin(dlat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(dlon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// 距离（单位：公里）
	distance := earthRadius * c

	return distance
}

// =================================================================

// 找到所有申请成功的用户 ID
func GetUserIDListByActivityID(activityID int) ([]string, error) {
	var userIDs []string
	result := db.Model(&models.ActivityMember{}).Where("activity_id =? and status = 1", activityID).Select("user_id").Find(&userIDs)
	if result.Error != nil {
		return nil, result.Error
	}

	return userIDs, nil
}
