package model

import "time"

type Actor struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	FirstName string    `json:"first_name" gorm:"size:32;not null"`
	LastName  string    `json:"last_name" gorm:"size:32;not null"`
	Role      string    `json:"role" gorm:"size:32;default:'actor';not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type ActorList struct {
	Actors []Actor `json:"actors"`
	Total  int64   `json:"total"`
}

type ActorID struct {
	ID uint `json:"id"`
}
