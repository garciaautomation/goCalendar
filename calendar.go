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

// Dependencies
// go get -u google.golang.org/api/calendar/v3
// go get -u golang.org/x/oauth2/google

// func bindCommandWithAliases(key, description string, cmd command.Cmd, requiredFlags []string) {
// 	command.On(key, description, cmd, requiredFlags)
// 	aliases, ok := drive.Aliases[key]
// 	if ok {
// 		for _, alias := range aliases {
// 			command.On(alias, description, cmd, requiredFlags)
// 		}
// 	}
// }

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
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

func main() {
	ctx := context.Background()
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	// config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	argparse(srv)
	// command.ParseAndRun()
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

type eventStruct struct {
	summary     string
	location    string `default: "here"`
	description string
	start       string
	end         string
	// recurrence  int
}

func (e *eventStruct) eventDefaults() {
	// now := time.Now()
	e.summary = "Default Summary"
	e.description = "Default Description"
	e.location = "Here"
	e.start = time.Now().Add(30 * time.Minute).Format("2006-01-02T15:04:05-0700")
	e.end = time.Now().Add(90 * time.Minute).Format("2006-01-02T15:04:05-0700")

}

func addEvent(srv *calendar.Service, calId string, event string) {
	// e := calendar.NewEventsService(srv)
	// calendar.EventCreator
	p := new(eventStruct)
	p.eventDefaults()
	fmt.Printf("p.location: %v\n", p)
	e := &calendar.Event{
		Summary:     event,
		Location:    p.location,
		Description: p.description,
		Start: &calendar.EventDateTime{
			DateTime: p.start,
			// TimeZone: "America/Chicago",
		},
		End: &calendar.EventDateTime{
			DateTime: p.end,
			// TimeZone: "America/Chicago",
		},
		// Recurrence: []string{"RRULE:FREQ=WEEKLY;COUNT=2"},
		// Attendees: []*calendar.EventAttendee{
		// 	&calendar.EventAttendee{Email: "lpage@example.com"},
		// 	&calendar.EventAttendee{Email: "sbrin@example.com"},
		// },
	}

	// e, err := srv.Events.QuickAdd(calId, event).Do()
	e, err := srv.Events.Insert(calId, e).Do()
	if err != nil {
		log.Fatalf("Unable to create event. %v\n", err)
	}
	fmt.Printf("Event created: %s\n", e.HtmlLink)
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

// type listCmd struct{}

// func (cmd *listCmd) Flags(fs *flag.FlagSet) *flag.FlagSet {
// 	return fs
// }

// func (cmd *listCmd) Run(srv calendar.Service) {
// 	fmt.Printf("\"asdf\": %v\n", "asdf")
// }
