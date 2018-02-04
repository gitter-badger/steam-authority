package mysql

import (
	"strconv"

	"github.com/Jleagle/go-helpers/logger"
)

type Tag struct {
	ID    int
	Name  string
	Games int
	Votes int
}

func (tag Tag) GetPath() string {
	return "/apps?tag=" + strconv.Itoa(tag.ID)
}

func GetAllTags() (tags []Tag, err error) {

	conn, err := getDB()
	if err != nil {
		logger.Error(err)
		return tags, err
	}

	err = conn.Select(&tags, "SELECT * FROM tags ORDER BY games DESC LIMIT 1000")
	if err != nil {
		logger.Error(err)
		return tags, err
	}

	return tags, nil
}
