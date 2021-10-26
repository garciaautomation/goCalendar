package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"

	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	homedir := GetHomeDir()

	tokenFile := homedir + "/.config/goCalendar/token.json"
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokenFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func GetHomeDir() string {
	h, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return h
}

func main() {

	homedir := GetHomeDir()
	ctx := context.Background()
	read, err := ioutil.ReadFile(homedir + "/.config/goCalendar/credentials.json")
	// readwrite, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	// config, err := google.ConfigFromJSON(readwrite, calendar.CalendarScope)
	config, err := google.ConfigFromJSON(read, calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	argparse(srv)
}

func upcomingEvents(srv *calendar.Service, cal string) {
	t := time.Now().Format(time.RFC3339)
	// l := srv.Events.List()
	// for _, v := range l {
	// fmt.Printf("l: %v\n", l)
	// }
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

func list(srv *calendar.Service, opt string, opt2 string) {
	switch opt {
	case "calendars":
		listCalendars(srv)
	case "events":
		fmt.Printf("opt2: %v\n", opt2)
		upcomingEvents(srv, opt2)
	}
}

type Event struct {
	calendar.Event
}

func (a Event) eventDefaults() *calendar.Event {
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

func addEvent(srv *calendar.Service, calId string, name string) {
	q := new(Event)
	event := q.eventDefaults()
	// spew.Dump(event)
	event, err := srv.Events.Insert(calId, event).Do()
	if err != nil {
		log.Fatalf("Unable to create event. %v\n", err)
	}
	fmt.Printf("Event created: %s\n", event.HtmlLink)
}

func deleteEvent(srv *calendar.Service, calId string, event string) {
	e := srv.Events.Delete(calId, event).Do()
	if e != nil {
		log.Fatal(e.Error())
	}
	fmt.Printf("Event delted: %s\n", event)
}

func argparse(srv *calendar.Service) {
	flag.Parse()
	// a := flag.Args()
	if len(flag.Args()) < 1 {
		General()
	}

	cmd := flag.Arg(0)
	opt := flag.Arg(1)
	opt2 := flag.Arg(2)

	switch cmd {
	case "list":
		list(srv, opt, opt2)
	case "add":
		addEvent(srv, opt, opt2)
	case "delete":
		deleteEvent(srv, opt, opt2)

	}

}

func General() {
	fmt.Println("Basic Help")
}
