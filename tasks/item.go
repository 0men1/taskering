package tasks

import (
	"google.golang.org/api/tasks/v1"
)

type Categories struct {
	CatList      []Category
	CurrentIndex int
	Size         int
}

type Category struct {
	Title string
	Id    string
	Items []*tasks.Task
}
