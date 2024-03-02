package calendars


import (
	"google.golang.org/api/calendar/v3"
	"encoding/json"
	"os"
	"fmt"
	"log"
	"time"
)





func saveCalendars(path string, calendars[] itemCalendar) {
	fmt.Printf("Saving calendars to file: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to save calendars to file: %s\n", err)
	}

	defer f.Close()

	b,_ := json.MarshalIndent(calendars, "", "	")
	f.Write(b)
}


func FindEvents(srv *calendar.Service) (*[]itemCalendar) {
	t := time.Now().Format(time.RFC3339)

	// All calendars associated with user
	_, err := os.Stat("userdata")
	if err != nil {
		err := os.Mkdir("userdata", 0777)
		if err != nil {
			log.Fatalf("Couldn't make the userdata directory: %v", err)
		}
	}


	calendarFile := "userdata/calendars.json"
	allCalendars := []itemCalendar{}
	calendars, err := srv.CalendarList.List().ShowDeleted(false).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve all existing calendars: %v", err)
	}
	if len(calendars.Items) == 0 {
		fmt.Println("Could not find an calendars")

	} else {
		for _, calendar := range calendars.Items {

			allEvents := []itemEvent{}
			events, err := srv.Events.List(calendar.Id).ShowDeleted(false).
			SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
			if err != nil {
				log.Fatalf("Unable to retrieve the next 10 user events of calendar: %s. %v",
				calendar.Summary, err)
			}
			if len(events.Items) == 0 {
				fmt.Printf("No upcoming events found in calendar: %s\n", calendar.Summary)
			} else {
				for _, item := range events.Items {
					if (item.Summary != "") {
						allEvents = append(allEvents, makeItemEventStruct(item))
					}
				}
			}

			allCalendars = append(allCalendars, makeItemCalendarStruct(calendar, allEvents))
		}
	}
	
	saveCalendars(calendarFile, allCalendars)
	return &allCalendars
}

