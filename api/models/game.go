package models

type Game struct {
	Base
	Title      string `json:"title"       gorm:"not null"`
	IsFinished bool   `json:"is_finished" gorm:"default:false"`
}
