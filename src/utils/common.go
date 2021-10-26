package utils

import (
	"log"
	"os"
)

func GetHomeDir() string {
	h, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return h
}
