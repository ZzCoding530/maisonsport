package dao

import "maisonsport/models"

func GetActivityMembersByActivityID(activityID int) ([]models.ActivityMember, error) {
	var members []models.ActivityMember

	// 查询符合条件的记录
	result := db.Where("activity_id = ?", activityID).Find(&members)
	if result.Error != nil {
		return nil, result.Error
	}

	return members, nil
}
