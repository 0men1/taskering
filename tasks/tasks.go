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


func MakeTask(title string, due string, notes string) (*tasks.Task) {
	myTask := tasks.Task {
		Title: title,
		Due: due,
		Notes: notes,
	}
	return &myTask	
}


func InsertTask(taskId string, task *tasks.Task) {
	ok := task_srv.Tasks.Insert(taskId, task)
	if ok == nil {
		log.Fatalf("There was an error creating task\n")	
	}
}


func makeCategory(tasklist *tasks.TaskList, tasks []Item) Categories {
	c := Category {
		Title: tasklist.Title,
		Id:    tasklist.Id,    
		Items: tasks,
	}
	AllCategories.CatList = append(AllCategories.CatList, c)

	return AllCategories
}


func makeItem(task *tasks.Task) Item {
	var t Item

	tB, err := json.Marshal(task); if err != nil {
		log.Fatalf("There was an error Marshalling task: %v", err)
	}

	log.Println(task)


	if err = json.Unmarshal([]byte(tB), &t); err != nil {
		log.Fatalf("There was an error Unmarshalling task: %v", err)
	}

	return t
}


func Init() (*Categories){
	srv  := api.GetSrvs()
	return FindTasks(srv)
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
			allTasks := []Item{}
			tasks, err := srv.Tasks.List(tasklist.Id).MaxResults(30).Do() 
			if err != nil {
				log.Fatalf("Unable to retrive the next 30 tasks from tasklist %s: %v\n", 
				
				tasklist.Title, 
				err)
			}
			for _, task := range tasks.Items {
				allTasks = append(allTasks, makeItem(task))
			}
			makeCategory(tasklist, allTasks)
		}
	}
	return &AllCategories
}

