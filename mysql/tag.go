package mysql

import "strconv"

type Tag struct {
	ID    int
	Name  string
	Games int
	Votes int
}

func (tag Tag) GetPath() string {
	return "/apps?tag=" + strconv.Itoa(tag.ID)
}
