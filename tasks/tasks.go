package tasks

import (
	"tsk/api"
	"google.golang.org/api/tasks/v1"
	"encoding/json"
	"os"
	"fmt"
	"log"
)


var AllCategories = Categories{
	CatMap:       make(map[int][]Item),
	CatList:      []Category{},
	CurrentIndex: 0,
	Size:	      0,
}


func makeCategory(tasklist *tasks.TaskList, tasks []Item) Categories {
	c := Category {
		Title: tasklist.Title,
		Id:    tasklist.Id,    
	}
	AllCategories.CatMap[AllCategories.Size] = tasks
	AllCategories.CatList = append(AllCategories.CatList, c)
	AllCategories.Size++

	return AllCategories
}


func makeItem(task *tasks.Task) Item {
	var t Item
	tB, err := json.Marshal(task); if err != nil {
		log.Fatalf("There was an error Marshalling task: %v", err)
	}
	if err = json.Unmarshal([]byte(tB), &t); err != nil {
		log.Fatalf("There was an error Unmarshalling task: %v", err)
	}
	return t
}

func Init() (*Categories){
	srv,_  := api.GetSrvs()
	return FindTasks(srv)
}


func saveTasks(path string, taskLists []Category) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("There was an error opening file %s: %v\n",
		path, err)
	}
	defer f.Close()
	fmt.Printf("Saving tasks to %s\n", path)
	b, _:= json.MarshalIndent(taskLists, "", "	")
	f.Write(b)
}


func FindTasks(srv *tasks.Service) (*Categories) {
	/*
	_, err := os.Stat("userdata")
	if err != nil {
		err := os.Mkdir("userdata", 0777)
		if err != nil {
			log.Fatalf("Couldn't make the userdata directory: %v", err)
		}
	}

	taskFile := "userdata/tasks.json"
	*/

	tasklists, err := srv.Tasklists.List().MaxResults(5).Do()
	if err != nil {
		fmt.Println("Unable to retrieve the tasklists: %v\n", err)
	}
	if len(tasklists.Items) == 0 {
		fmt.Println("Could not find tasklists")
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

//	saveTasks(taskFile, allTaskLists)
	return &AllCategories
}

