package model

import "time"

type Movie struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Title     string    `json:"title" gorm:"size:255;not null"`
	Director  string    `json:"director" gorm:"size:255;not null"`
	Year      int       `json:"year" gorm:"not null"`
	Plot      string    `json:"plot" gorm:"type:text"`
	Cast      []Actor   `json:"cast" gorm:"many2many:movie_actors"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type MovieList struct {
	Movies []Movie `json:"movies"`
	Count  int     `json:"count"`
}
