package cal

import (
	"fmt"
	"log"
	"time"

	"google.golang.org/api/calendar/v3"
)

func ListCalendars(srv *calendar.Service) []*calendar.CalendarListEntry {
	calendars, err := srv.CalendarList.List().Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's calendars: %v", err)
	}
	fmt.Println("Upcoming calendars:")
	if len(calendars.Items) == 0 {
		fmt.Println("No calendars found.")
	} else {
		// for _, item := range calendars.Items {
		// 	fmt.Printf("\t%v :: %v\n", item.Summary, item.Id)
		// }
	}
	return calendars.Items
}

func UpcomingEvents(srv *calendar.Service, cal string) {
	t := time.Now().Format(time.RFC3339)
	events, err := srv.Events.List(cal).ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}
	fmt.Println("Upcoming events:")
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		for _, item := range events.Items {
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}
			fmt.Printf("\t%v :: %v :: (%v)\n", item.Summary, item.Id, date)
		}
	}
}

func List(srv *calendar.Service, opt string, opt2 string) {
	switch opt {
	case "calendars":
		ListCalendars(srv)
	case "events":
		UpcomingEvents(srv, opt2)
	}
}

func eventDefaults() *calendar.Event {
	e := &calendar.Event{}

	e.Summary = "Default Summary"
	e.Description = "Default Description"
	e.Location = "Here"
	e.Start = &calendar.EventDateTime{
		DateTime: time.Now().Add(30 * time.Minute).Format("2006-01-02T15:04:05-0700"),
		TimeZone: "America/Chicago",
	}
	e.End = &calendar.EventDateTime{
		DateTime: time.Now().Add(90 * time.Minute).Format("2006-01-02T15:04:05-0700"),
		TimeZone: "America/Chicago",
	}
	e.Visibility = "public" // default private public
	e.GuestsCanModify = true
	e.Attendees = append(e.Attendees, &calendar.EventAttendee{Email: "secret@gmail.com"})
	// e.Recurrence = []string{"RRULE:FREQ=WEEKLY;COUNT=2"}
	// Recurrence: []string{"RRULE:FREQ=WEEKLY;COUNT=2"},
	// Attendees: []*calendar.EventAttendee{
	// 	&calendar.EventAttendee{Email: "lpage@example.com"},
	// 	&calendar.EventAttendee{Email: "sbrin@example.com"},

	return e
}

func AddEvent(srv *calendar.Service, calendarId string, name string) {
	// e := new(Event)
	event := eventDefaults()
	event.Summary = name
	fmt.Printf("name: %v\n", name)
	// spew.Dump(event)
	event, err := srv.Events.Insert(calendarId, event).Do()
	if err != nil {
		log.Fatalf("Unable to create event. %v\n", err)
	}
	fmt.Printf("Event created: %s\n\t%s\n", event.HtmlLink, event.Id)
}

func DeleteEvent(srv *calendar.Service, calId string, event string) {
	e := srv.Events.Delete(calId, event).Do()
	if e != nil {
		log.Fatal(e.Error())
	}
	fmt.Printf("Event delted: %s\n", event)
}
