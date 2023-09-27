package models

type ActivityMember struct {
	ID         int    `gorm:"column:id;primaryKey" json:"id"`        // 对应 MySQL 表字段: id
	ActivityID int    `gorm:"column:activity_id" json:"activity_id"` // 对应 MySQL 表字段: activity_id
	UserID     string `gorm:"column:user_id" json:"user_id"`         // 对应 MySQL 表字段: user_id
	Status     int    `gorm:"column:status" json:"status"`           // 对应 MySQL 表字段: status
}

// TableName specifies the table name for the ActivityMember model
func (ActivityMember) TableName() string {
	return "t_activity_member"
}
