package mysql

import (
	"strconv"

	"github.com/Jleagle/go-helpers/logger"
)

type Genre struct {
	ID   int
	Name string
}

func (genre Genre) GetPath() string {
	return "/apps?genre=" + strconv.Itoa(genre.ID)
}

func GetAllGenres() (tags []Tag, err error) {

	conn, err := getDB()
	if err != nil {
		logger.Error(err)
		return tags, err
	}

	err = conn.Select(&tags, "SELECT * FROM genres ORDER BY name DESC LIMIT 1000")
	if err != nil {
		logger.Error(err)
		return tags, err
	}

	return tags, nil
}
