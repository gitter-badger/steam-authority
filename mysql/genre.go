package mysql

import (
	"time"
)

type Genre struct {
	Name      string     `gorm:"not null;column:name;primary_key"`
	CreatedAt *time.Time `gorm:"not null;column:created_at"`
	UpdatedAt *time.Time `gorm:"not null;column:updated_at"`
	Games     int        `gorm:"not null;column:games"`
}

func (genre Genre) GetPath() string {
	return "/apps?genre=" + genre.Name
}

func GetAllGenres() (genres []Genre, err error) {

	db, err := getDB()
	if err != nil {
		return genres, err
	}

	db.Limit(1000).Order("name ASC").Find(&genres)
	if db.Error != nil {
		return genres, err
	}

	return genres, nil
}
