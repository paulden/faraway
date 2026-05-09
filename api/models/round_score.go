package models

import "gorm.io/gorm"

type RoundScore struct {
	gorm.Model
	Score    int  `json:"score"`
	RoundID  uint `json:"round_id"  gorm:"not null"`
	PlayerID uint `json:"player_id" gorm:"not null"`
}
