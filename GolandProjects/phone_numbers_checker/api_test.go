package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
)

type Headers struct {
	Headers []struct {
		Name  string
		Value string
	}
}

func parseHeaders(headers *Headers) error {
	file, err := os.ReadFile("./headers.json")
	if err != nil {
		log.Println("Error while reading headers: ", err)
	}
	err = json.Unmarshal(file, &headers)
	if err != nil {
		log.Println("Error while decoding JSON: ", err)
	}
	return err
}

func TestAPICall(t *testing.T) {
	client := &http.Client{}
	headers := Headers{}
	err := parseHeaders(&headers)
	if err != nil {
		log.Fatal(err)
	}
	request, err := http.NewRequest("GET", "https://ib.rencredit.ru/rencredit.server.portal.app/rest/private/transfers/internal/register", nil)
	for _, header := range headers.Headers {
		request.Header.Set(header.Name, header.Value)
	}
	response, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	body := response.StatusCode
	fmt.Println(body)
}
