package model

type Campus struct {
	ID   uint   `gorm:"primary_key" json:"id"`
	Name string `sql:"type: VARCHAR(255); not null; unique" json:"campusName"`
	User User
}
