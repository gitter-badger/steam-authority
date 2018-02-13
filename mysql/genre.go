package mysql

import (
	"strconv"

	"github.com/Jleagle/go-helpers/logger"
)

type Genre struct {
	ID   int    `gorm:"not null;column:id;primary_key;AUTO_INCREMENT"`
	Name string `gorm:"not null;column:name"`
}

func (genre Genre) GetPath() string {
	return "/apps?genre=" + strconv.Itoa(genre.ID)
}

func GetAllGenres() (tags []Tag, err error) {

	db, err := getDB()
	if err != nil {
		logger.Error(err)
		return tags, err
	}

	db.Limit(1000).Order("name DESC").Find(&tags)
	if db.Error != nil {
		return tags, err
	}

	return tags, nil
}
