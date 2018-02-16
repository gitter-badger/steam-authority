package mysql

import (
	"strconv"
	"time"

	"github.com/Jleagle/go-helpers/logger"
)

type Genre struct {
	ID        int        `gorm:"not null;column:id;primary_key;AUTO_INCREMENT"`
	CreatedAt *time.Time `gorm:"not null;column:created_at"`
	UpdatedAt *time.Time `gorm:"not null;column:updated_at"`
	Name      string     `gorm:"not null;column:name"`
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
