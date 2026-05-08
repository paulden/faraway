package models

import "gorm.io/gorm"

type Player struct {
	gorm.Model
	Name   string `json:"name"    gorm:"not null"`
	GameID uint   `json:"game_id" gorm:"not null"`
}
