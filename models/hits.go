package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Hit struct {
	gorm.Model
	ID       uint      `json:"id"`
	HashID   uint      `json:"hash_id"`
	Ip       string    `json:"ip"`
	AccessAt time.Time `json:"access_at"`
}
