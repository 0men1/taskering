package tasks
import (
	"google.golang.org/api/tasks/v1"
)



type Categories struct {
	CatList      []Category
	CurrentIndex int
	Size	     int
}


type Category struct {
	Title string
	Id    string
	Items []*tasks.Task
}


type Item struct {
	Id    		string `json:"Id"`
	Title 		string `json:"Title"`
	Due   		string `json:"Due"`
	Link  		string `json:"SelfLink"`
	Status		string `json:"Status"`// "needsAction" or "completed"
	Notes 		string `json:"Notes"`
}
