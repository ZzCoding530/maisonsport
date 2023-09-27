package models

import "time"

type ActivityInfo struct {
	ID        int       `gorm:"column:id;primaryKey" json:"id"`      // 对应 MySQL 表字段: id
	Title     string    `gorm:"column:title" json:"title"`           // 对应 MySQL 表字段: title
	TarGender string    `gorm:"column:tar_gender" json:"tar_gender"` // 对应 MySQL 表字段: tar_gender
	TarLevel  float64   `gorm:"column:tar_level" json:"tar_level"`   // 对应 MySQL 表字段: tar_level
	Time      time.Time `gorm:"column:time" json:"time"`             // 对应 MySQL 表字段: time
	Note      string    `gorm:"column:note" json:"note"`             // 对应 MySQL 表字段: note
	MaxMember int       `gorm:"column:max_member" json:"max_member"` // 对应 MySQL 表字段: max_member
	PositionX float64   `gorm:"column:position_x" json:"position_x"`
	PositionY float64   `gorm:"column:position_y" json:"position_y"`
	City      string    `gorm:"column:city" json:"city"`
	Creator   string    `gorm:"column:creator" json:"creator"`
}

// TableName specifies the table name for the Activity model
func (ActivityInfo) TableName() string {
	return "t_activity_info"
}

// 按照条件筛选活动
// Filter 定义查询条件参数结构体
type Filter struct {
	City      string  `json:"city"`       // 城市
	Gender    string  `json:"gender"`     // 性别
	LevelMin  float64 `json:"level_min"`  // 级别下限
	LevelMax  float64 `json:"level_max"`  // 级别上限
	MemberMin int     `json:"member_min"` // 人数下限
	MemberMax int     `json:"member_max"` // 人数上限
	Date      string  `json:"date"`       // 日期
	Time      int     `json:"time"`       // 时间段
	PositionX float64 `json:"positon_x"`  // x坐标
	PositionY float64 `json:"position_y"` // y坐标
	Pagesize  int     `json:"pagesize"`   // 每页条数
	PageIndex int     `json:"pageindex"`
}
