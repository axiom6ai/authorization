package model

type User struct {
	ID    uint   `gorm:"primary_key" json:"id"`
	Email string `sql:"type: VARCHAR(255); not null; unique" json:"email"`
}
