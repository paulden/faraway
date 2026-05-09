package models

import "gorm.io/gorm"

type Round struct {
	gorm.Model
	Number uint `json:"number" gorm:"not null"`
	GameID uint `json:"game_id" gorm:"not null"`
}
