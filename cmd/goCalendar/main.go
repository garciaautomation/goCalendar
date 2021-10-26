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

	"github.com/garciaautomation/goCalendar/cal"
	"github.com/garciaautomation/goCalendar/help"
	"github.com/garciaautomation/goCalendar/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	homedir := utils.GetHomeDir()
	// help.General()
	// help.

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

// func GetHomeDir() string {
// 	h, err := os.UserHomeDir()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return h
// }

func main() {

	homedir := utils.GetHomeDir()
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

func argparse(srv *calendar.Service) {
	flag.Parse()
	// a := flag.Args()
	if len(flag.Args()) < 1 {
		help.General()
	}

	cmd := flag.Arg(0)
	opt := flag.Arg(1)
	opt2 := flag.Arg(2)

	switch cmd {
	case "list":
		cal.List(srv, opt, opt2)
	case "add":
		cal.AddEvent(srv, opt, opt2)
	case "delete":
		cal.DeleteEvent(srv, opt, opt2)

	}
}
