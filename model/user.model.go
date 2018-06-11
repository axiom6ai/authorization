package model

type User struct {
	ID          uint   `gorm:"primary_key" json:"id"`
	Email       string `sql:"type: VARCHAR(255); not null; unique" json:"email"`
	ChineseName string `sql:"type: VARCHAR(255);" json:"chineseName"`
	EnglishName string `sql:"type: VARCHAR(255);" json:"englishName"`
	DateOfBirth string `sql:"type: VARCHAR(255);" json:"dateOfBirth"`
	PhoneNumber string `sql:"type: VARCHAR(255);" json:"phoneNumber"`
	PhotoPath   string `sql:"type: VARCHAR(255);" json:"photoPath"`
}
