package mysql

import (
	"strconv"
	"time"
)

type Tag struct {
	ID        int        `gorm:"not null;column:id;primary_key;AUTO_INCREMENT"`
	CreatedAt *time.Time `gorm:"not null;column:created_at"`
	UpdatedAt *time.Time `gorm:"not null;column:updated_at"`
	Name      string     `gorm:"not null;column:name"`
	Apps      int        `gorm:"not null;column:apps"`
	Votes     int        `gorm:"not null;column:votes"`
}

func (tag Tag) GetPath() string {
	return "/apps?tag=" + strconv.Itoa(tag.ID)
}

func GetAllTags() (tags []Tag, err error) {

	db, err := getDB()
	if err != nil {
		return tags, err
	}

	db = db.Limit(1000).Order("name ASC").Find(&tags)
	if db.Error != nil {
		return tags, db.Error
	}

	return tags, nil
}

func SaveOrUpdateTag(id int, apps int) (err error) {

	db, err := getDB()
	if err != nil {
		return err
	}

	tag := new(Tag)
	db.Assign(Tag{Apps: apps}).FirstOrCreate(tag, Tag{ID: id})
	if db.Error != nil {
		return db.Error
	}

	return nil
}
