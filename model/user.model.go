package model

type User struct {
	ID               uint   `gorm:"primary_key" json:"id"`
	Email            string `sql:"type: VARCHAR(255); not null; unique" json:"email"`
	Password         string `sql:"type: VARCHAR(255); not null;" json:"password"`
	DateOfBirth      string `sql:"type: VARCHAR(255); not null" json:"dateOfBirth"`
	FirstName        string `sql:"type: VARCHAR(255); not null" json:"firstName"`
	LastName         string `sql:"type: VARCHAR(255); not null" json:"lastName"`
	PhoneNumber      string `sql:"type: VARCHAR(255); not null" json:"phoneNumber"`
	Gender           string `sql:"type: VARCHAR(255); not null" json:"gender"`
	ChineseLastName  string `sql:"type: VARCHAR(255);" json:"chineseLastName"`
	ChineseFirstName string `sql:"type: VARCHAR(255);" json:"chineseFirstName"`
	PhotoPath        string `sql:"type: VARCHAR(255);" json:"photoPath"`
}
