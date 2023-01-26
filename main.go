package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Headers struct {
	Headers []struct {
		Name  string
		Value string
	}
}

type Cookies struct {
	Cookies []struct {
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

func parseCookies(cookies *Cookies) error {
	file, err := os.ReadFile("./cookies.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &cookies)
	if err != nil {
		log.Fatalln("Error while decoding JSON: ", err)
	}
	return err
}

func main() {
	client := &http.Client{}
	headers := Headers{}
	//cookies := Cookies{}
	err := parseHeaders(&headers)
	if err != nil {
		log.Fatal(err)
	}
	//err = parseCookies(&cookies)
	url := "https://ib.rencredit.ru/rencredit.server.portal.app/rest/private/transfers/internal/register"
	var data []byte
	copy(data[:], "payerAccount=505419565&paymentMethod=BY_PHONE&destinationPhoneNumber=79005200846&amount=1&paymentPurposeCode=GIFT")
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	for _, header := range headers.Headers {
		request.Header.Set(header.Name, header.Value)
	}
	//for _, cookie := range cookies.Cookies {
	//	request.AddCookie(&http.Cookie{Name: cookie.Name, Value: cookie.Value})
	//}
	response, err := client.Do(request)
	defer response.Body.Close()
	if err != nil {
		log.Fatalln(err)
	}
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println(response.Status)
	log.Println(bodyString)

}
