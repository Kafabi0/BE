package models

type User struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Username   string `gorm:"uniqueIndex" json:"username"`
	Password   string `json:"password"`
}
