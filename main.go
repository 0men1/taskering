package main

import (
	//"tsk/bubble"
	tsktasks "tsk/tasks"
	"google.golang.org/api/tasks/v1"
	//"log"
)


func main() {
	srv := tsktasks.Init()


	myTask := &tasks.Task {
		Title: "TEST TASK INSERT",
		Due: "2023-03-07T14:45:00Z",
		Notes: "TESTINSERTTASK",
	}



	tsktasks.InsertTask(srv.CatList[0].Id, myTask)
	
	//log.Println(srv.CatList[0].Id)	


	
	//log.Println(myTask)
	//bubble.Run()	
}

