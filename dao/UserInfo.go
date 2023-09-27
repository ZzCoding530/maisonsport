package dao

import "maisonsport/models"

// 向基础信息表中插入部分字段的数据
func InsertUserInfoPart(userID, openID, token string) error {
	userInfo := models.UserInfo{
		UserID: userID,
		OpenID: openID,
		Token:  token,
	}

	result := db.Create(&userInfo)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// 从基础信息表中拿出IsInfoVerified字段，用 user_id
//
//	0——未填写
//	1——审核通过
//	2——审核中
//	3——审核未通过
func GetUserInfoStatusByUserID(user_id string) (int, error) {
	var userInfo models.UserInfo

	// Find the user by user_id
	result := db.Where("user_id = ?", user_id).First(&userInfo)
	if result.Error != nil {
		return 0, result.Error
	}

	return userInfo.IsInfoVerified, nil
}

// 根据 userid 去更新基础信息表，用户填写信息时候用的
func UpdateUserInfoByUserID(userID string, updatedUserInfo models.UserInfo) error {
	// Update the user info for the specified user_id
	result := db.Model(&models.UserInfo{}).Where("user_id = ?", userID).Updates(updatedUserInfo)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// GetUserByUserID 根据 userID 查询用户信息
func GetUserInfoByUserID(userID string) (models.UserInfo, error) {
	var userInfo models.UserInfo
	result := db.Where("user_id = ?", userID).First(&userInfo)
	if result.Error != nil {
		return models.UserInfo{}, result.Error
	}

	return userInfo, nil
}

// GetUserInfoByUserIDs 根据 userID 批量查询用户信息
func GetUserInfoByUserIDs(userIDs []string) ([]models.UserInfo, error) {
	var userInfos []models.UserInfo

	if len(userIDs) == 0 {
		return userInfos, nil
	}

	result := db.Where("user_id IN (?)", userIDs).Find(&userInfos)
	if result.Error != nil {
		return nil, result.Error
	}

	return userInfos, nil
}
