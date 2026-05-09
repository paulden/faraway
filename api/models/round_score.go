package models

type RoundScore struct {
	Base
	Score    int  `json:"score"`
	RoundID  uint `json:"round_id"  gorm:"not null"`
	PlayerID uint `json:"player_id" gorm:"not null"`
}
