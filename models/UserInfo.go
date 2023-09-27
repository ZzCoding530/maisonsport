package models

type UserInfo struct {
	ID             uint    `gorm:"column:id;primaryKey"`
	UserID         string  `gorm:"column:user_id;not null;unique"`
	OpenID         string  `gorm:"column:open_id;not null"`
	Token          string  `gorm:"column:token;not null"`
	NickName       string  `gorm:"column:nick_name;not null"`
	Gender         string  `gorm:"column:gender;not null"`
	Age            int     `gorm:"column:age;not null"`
	SkillLevel     float64 `gorm:"column:skill_level"`
	ProvinceName   string  `gorm:"column:province_name"`
	CityName       string  `gorm:"column:city_name"`
	DistrictName   string  `gorm:"column:district_name"`
	AvatarURL      string  `gorm:"column:avatar_url"`
	VideoUrl       string  `gorm:"column:video_url"`
	CellPhone      string  `gorm:"column:cell_phone"`
	IsInfoVerified int     `gorm:"column:is_info_verified;not null"`
}

// Set the table name for the UserInfo model
func (UserInfo) TableName() string {
	return "t_user_info"
}
