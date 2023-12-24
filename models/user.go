package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	GoogleID string
	Email    string
	Name     string
}
