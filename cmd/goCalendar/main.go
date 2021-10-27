package main

import (
	"flag"

	"log"

	"github.com/garciaautomation/goCalendar/cal"
	"github.com/garciaautomation/goCalendar/help"
	"github.com/garciaautomation/goCalendar/utils"
	"google.golang.org/api/calendar/v3"
)

func main() {

	srv, err := utils.GetSrv()
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	utils.AddCompletion()
	argparse(srv)
}

func argparse(srv *calendar.Service) {
	flag.Parse()
	if len(flag.Args()) < 1 {
		help.General("goCalendar")
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
