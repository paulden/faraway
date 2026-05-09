package models

type Player struct {
	Base
	Name   string `json:"name"    gorm:"not null"`
	GameID uint   `json:"game_id" gorm:"not null"`
}
