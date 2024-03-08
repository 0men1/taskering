package tasks

import (
	"tsk/api"
	"google.golang.org/api/tasks/v1"
	"encoding/json"
	"log"
)


var AllCategories = Categories{
	CatList:      []Category{},
	CurrentIndex: 0,
}


var task_srv *tasks.Service


func Init() (*Categories){
	task_srv = api.GetSrvs()
	return FindTasks(task_srv)
}


func RefreshFull() (*Categories) {
	if task_srv == nil {
		log.Fatalf("Task service was not initialized")
	}

	return FindTasks(task_srv)
}


func MakeTask(title string, due string, notes string) (*tasks.Task) {
	myTask := tasks.Task {
		Title: title,
		Due: due,
		Notes: notes,
	}
	return &myTask	
}


func InsertTask(taskListId string, task *tasks.Task) *tasks.Task {
	newTask, err := task_srv.Tasks.Insert(taskListId, task).Do()

	if err != nil {
		log.Fatalf("There was an error inserting the task: %v", err)
	}

	return newTask
}


func makeCategory(tasklist *tasks.TaskList, allTasks *tasks.Tasks) Categories {
	taskItems := allTasks.Items
	c := Category {
		Title: tasklist.Title,
		Id:    tasklist.Id,    
		Items: taskItems,
	}
	AllCategories.CatList = append(AllCategories.CatList, c)
	return AllCategories
}


func MakeItem(task *tasks.Task) Item {
	var t Item
	tB, err := json.Marshal(task); if err != nil {
		log.Fatalf("There was an error Marshalling task: %v", err)
	}

	if err2 := json.Unmarshal([]byte(tB), &t); err != nil {
		log.Fatalf("There was an error Unmarshalling task: %v", err2)
	}
	return t
}


func FindTasks(srv *tasks.Service) (*Categories) {
	tasklists, err := srv.Tasklists.List().MaxResults(5).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve the tasklists: %v\n", err)
	}
	if len(tasklists.Items) == 0 {
		log.Fatalf("Could not find tasklists")
	} else {
		for _, tasklist := range tasklists.Items {
			allTasks, err := srv.Tasks.List(tasklist.Id).MaxResults(30).Do() 
			if err != nil {
				log.Fatalf("Unable to retrive the next 30 tasks from tasklist %s: %v\n", 
				tasklist.Title, 
				err)
			}
			makeCategory(tasklist, allTasks)
		}
	}
	return &AllCategories
}
