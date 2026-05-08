package models

import "gorm.io/gorm"

type Game struct {
	gorm.Model
	Title      string `json:"title"       gorm:"not null"`
	IsFinished bool   `json:"is_finished" gorm:"default:false"`
}
