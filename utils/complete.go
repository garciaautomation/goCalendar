package utils

import (
	"flag"
	"os"
	"path/filepath"
	"strconv"

	"github.com/posener/complete/v2"
	"github.com/posener/complete/v2/compflag"
)

var (
	// Add variables to the program.
	startDate = compflag.String("name", "", "Give your name")
	startTime = compflag.String("something", "", "Expect somthing, but we don't know what, so no other completion options will be provided.")
	// nothing   = compflag.String("nothing", "", "Expect nothing after flag, so other completion can be provided.")
)

var (
	// Add variables to the program. Since we are using the compflag library, we can pass options to
	// enable bash completion to the flag values.
	list   = flag.String("list", "", "Give your name")
	add    = flag.String("add", "", "Expect somthing, but we don't know what, so no other completion options will be provided.")
	delete = flag.String("delete", "", "Expect nothing after flag, so other completion can be provided.")
)

type getEnvFn = func(string) string

var promptEnv = func(contents string) getEnvFn {
	return func(key string) string {
		switch key {
		case "COMP_LINE":
			return contents
		case "COMP_POINT":
			return strconv.Itoa(len(contents))
		}
		return ""
	}
}

func AddCompletion() {

	cmd := &complete.Command{
		Sub: map[string]*complete.Command{
			"list": {
				Sub: map[string]*complete.Command{
					"calendars": {},
					"events": {
						Sub: map[string]*complete.Command{"primary": {}},
					},
				},
			},
			"add": {
				Sub: map[string]*complete.Command{
					"<id>":    {Sub: map[string]*complete.Command{"startDate": {}}},
					"primary": {Sub: map[string]*complete.Command{"startDate": {}}}},
			},
			"delete": {
				Sub: map[string]*complete.Command{"<id>": {}, "primary": {}},
			},
		},
	}

	complete.Complete(os.Args[0], cmd)
	compflag.Parse()

}

func getBinaryPath() (string, error) {
	bin, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Abs(bin)
}
