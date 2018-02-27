package mysql

import (
	"time"
)

type Genre struct {
	ID        int        `gorm:"not null;column:id;primary_key;AUTO_INCREMENT"` //
	Name      string     `gorm:"not null;column:name"`
	CreatedAt *time.Time `gorm:"not null;column:created_at"`
	UpdatedAt *time.Time `gorm:"not null;column:updated_at"`
	Apps      int        `gorm:"not null;column:apps"`
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
		return genres, db.Error
	}

	return genres, nil
}

func SaveOrUpdateGenre(id int, name string, apps int) (err error) {

	db, err := getDB()
	if err != nil {
		return err
	}

	genre := new(Genre)
	db.Attrs(Genre{Name: name}).Assign(Genre{Apps: apps}).FirstOrCreate(genre, Genre{ID: id})
	if db.Error != nil {
		return db.Error
	}

	return nil
}
