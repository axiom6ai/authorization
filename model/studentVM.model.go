package model

type StudentVM struct {
	Name       string `gorm:"column:name" json:"name"`
	Email      string `gorm:"column:email" json:"email"`
	Phone      string `gorm:"column:phone" json:"phone"`
	Company    string `gorm:"column:company" json:"company"`
	CampusName string `gorm:"column:campusname" json:"campusName"`
}
