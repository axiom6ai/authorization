package model

type Student struct {
	ID       uint   `gorm:"primary_key" json:"id"`
	Name     string `sql:"type: VARCHAR(255); not null" json:"name"`
	Email    string `sql:"type: VARCHAR(255); not null; unique" json:"email"`
	Phone    string `sql:"type: VARCHAR(255); not null" json:"phone"`
	Company  string `sql:"type: VARCHAR(255)" json:"company"`
	CampusID uint   `json:"campusId"`
}
