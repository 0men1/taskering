package utils

import (
	"google.golang.org/api/tasks/v1"
)


func DeleteElement(slice []*tasks.Task, index int) []*tasks.Task {
	return append(slice[:index], slice[index+1:]...)

}
