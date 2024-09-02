package utils

import (
	"log"
)

func HandleFatalError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func HandleGracefulError(message string, err error) {
	if err != nil {
		log.Printf("%s: %v\n", message, err)
	}
}
