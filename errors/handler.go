package errors

import (
	"log"
)

func HandleError(helpText string, err *error) {
	if *err != nil {
		log.Fatalf("%s: %v", helpText, *err)
	}
}
