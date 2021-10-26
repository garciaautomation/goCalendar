package goCalendar

import (
	"fmt"
	"log"

	"google.golang.org/api/calendar/v3"
)

func listCalendars(srv *calendar.Service) {
	calendars, err := srv.CalendarList.List().Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's calendars: %v", err)
	}
	fmt.Println("Upcoming calendars:")
	if len(calendars.Items) == 0 {
		fmt.Println("No upcoming calendars found.")
	} else {
		for _, item := range calendars.Items {
			fmt.Printf("\t%v :: %v\n", item.Summary, item.Id)
		}
	}
}
