package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/playwright-community/playwright-go"
	"io"
	"log"
	"net/http"
	url2 "net/url"
	"strings"
)

func handleError(helpText string, err *error) {
	if *err != nil {
		log.Fatalf("%s: %v", helpText, *err)
	}
}

func askForPhone(url *string, headers *map[string]string, client *http.Client, c chan string) {
	form := url2.Values{}
	form.Add("payerAccount", "505419565")
	form.Add("paymentMethod", "BY_PHONE")
	form.Add("destinationPhoneNumber", "79005200846")
	form.Add("amount", "1")
	form.Add("paymentPurposeCode", "GIFT")
	request, err := http.NewRequest("POST", *url, strings.NewReader(form.Encode()))
	for key, value := range *headers {
		request.Header.Set(key, value)
	}
	handleError("error doing request: ", &err)
	resp, err := client.Do(request)
	handleError("error while request: ", &err)
	bodyBytes, err := io.ReadAll(resp.Body)
	handleError("Error while reading response bytes: ", &err)
	bodyString := string(bodyBytes)
	fmt.Println(resp.Status)
	c <- bodyString
}

func main() {
	env, err := godotenv.Read()
	handleError("Error while loading .env: ", &err)
	client := &http.Client{}
	url := "https://ib.rencredit.ru/rencredit.server.portal.app/rest/private/transfers/internal/register"

	pw, err := playwright.Run()
	handleError("Unable to run playwright", &err)
	headless := false
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless})
	handleError("Can't launch Chromium", &err)
	fmt.Println(env["BANK_LOGIN"], env["BANK_PASSWORD"])
	page := GetBrowserPage(browser)

	LoginToAccount(page, env["BANK_LOGIN"], env["BANK_PASSWORD"])

	timeout := 10000.0
	GetTransferPage(page, &timeout)

	headers := sendFirstPhoneRequest(page)

	respChan := make(chan string)
	for i := 0; i < 2; i++ {
		go askForPhone(&url, &headers, client, respChan)
		respText := <-respChan
		handleError("error while reading POST response: ", &err)
		fmt.Println(respText)
	}
	close(respChan)

	if err = browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
}
