package models

type Round struct {
	Base
	Number uint `json:"number" gorm:"not null"`
	GameID uint `json:"game_id" gorm:"not null"`
}
