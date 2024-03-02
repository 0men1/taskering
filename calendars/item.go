package	calendars

import(
	"google.golang.org/api/calendar/v3"
	"encoding/json"
	"log"
)

type itemEvent struct {
	Name		string 				`json:"Summary"`
	StartTime 	*calendar.EventDateTime		`json:"Start"`
	EndTime		*calendar.EventDateTime		`json:"EndTime"`
}	

type itemCalendar struct {
	Name 	    string           `json:"Summary"`
	Id   	    string	     `json:"Id"`
	Description string	     `json:"Description"`
	Events      []itemEvent  
}


func makeItemEventStruct(event *calendar.Event) itemEvent {
	var myEvent itemEvent 
	eventBytes, err := json.Marshal(event); if err != nil {
		log.Fatalf("Coudln't marshal event object: %v", err)	
	}
	if err := json.Unmarshal([]byte(eventBytes), &myEvent); err != nil {
		log.Fatalf("Coudln't unmarshal event Object: %v", err)
	}
	return myEvent
}


func makeItemCalendarStruct(calendar *calendar.CalendarListEntry, events []itemEvent)  itemCalendar {
	myCalendar := itemCalendar {
		Events:  events,
	}
	calendarBytes, err := json.Marshal(calendar); if err != nil {
		log.Fatalf("Couldn't marshal calendar object: %v", err)
	}
	if err := json.Unmarshal([]byte(calendarBytes), &myCalendar); err != nil {
		log.Fatalf("Couldn't unmarshal calendar object: %v", err)
	}
	return myCalendar
}
