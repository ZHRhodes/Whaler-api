package models

import (
	"time"
)

type DBModel struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

type Account struct {
	DBModel
	Name        string
	Industry    string
	Description string
}
