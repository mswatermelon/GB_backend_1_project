package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Hash struct {
	gorm.Model
	ID        uint      `json:"id"`
	Hash      string    `json:"title"`
	Url       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
