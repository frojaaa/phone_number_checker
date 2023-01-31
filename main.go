package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/playwright-community/playwright-go"
	"io"
	"log"
	"net/http"
	url2 "net/url"
	"strings"
	"sync"
	"time"
)

type Transfer struct {
	Transfer map[string]interface{} `json:"transfer"`
}

func handleError(helpText string, err *error) {
	if *err != nil {
		log.Fatalf("%s: %v", helpText, *err)
	}
}

func generateFormData() url2.Values {
	form := url2.Values{}
	form.Add("payerAccount", "589079858")
	form.Add("paymentMethod", "BY_PHONE")
	form.Add("destinationPhoneNumber", "79181185688")
	form.Add("amount", "1")
	form.Add("paymentPurposeCode", "GIFT")
	return form
}

func makeRequest(url *string, form url2.Values, headers *map[string]string) (*http.Request, error) {
	request, err := http.NewRequest("POST", *url, strings.NewReader(form.Encode()))
	for key, value := range *headers {
		request.Header.Set(key, value)
	}
	return request, err
}

func getResponse(request *http.Request, client *http.Client) (*http.Response, error) {
	resp, err := client.Do(request)
	return resp, err
}

func getJSON(response *http.Response, target interface{}) error {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln("error while closing response body: ", err)
		}
	}(response.Body)
	return json.NewDecoder(response.Body).Decode(target)
}

func askForPhone(url *string, headers *map[string]string, client *http.Client, tasks chan int, results chan string, wg *sync.WaitGroup, id int, target *Transfer) {
	defer wg.Done()
	for num := range tasks {
		form := generateFormData()
		//form.Add("payerAccount", "589079858")
		//form.Add("paymentMethod", "BY_PHONE")
		//form.Add("destinationPhoneNumber", "79181185688")
		//form.Add("amount", "1")
		//form.Add("paymentPurposeCode", "GIFT")
		request, err := makeRequest(url, form, headers)
		handleError("error doing request: ", &err)
		resp, err := getResponse(request, client)
		handleError("error while request: ", &err)
		err = getJSON(resp, target)
		bankName := target.Transfer["payeeBankName"].(string)
		//bodyBytes, err := io.ReadAll(resp.Body)
		//handleError("Error while reading response bytes: ", &err)
		//bodyString := string(bodyBytes)
		fmt.Printf("[worker %d] Worker Sending result of task %d\n", id, num)
		results <- bankName
	}
}

func main() {
	env, err := godotenv.Read()
	handleError("Error while loading .env: ", &err)
	transport := http.Transport{DisableKeepAlives: true, MaxIdleConns: 200}
	client := &http.Client{Timeout: 30 * time.Second, Transport: &transport}
	url := "https://ib.rencredit.ru/rencredit.server.portal.app/rest/private/transfers/internal/register"
	pw, err := playwright.Run()
	handleError("Unable to run playwright", &err)
	headless := false
	browser, err := pw.Firefox.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless})
	handleError("Can't launch Chromium", &err)
	fmt.Println(env["BANK_LOGIN"], env["BANK_PASSWORD"])
	page := GetBrowserPage(browser)

	LoginToAccount(page, env["BANK_LOGIN"], env["BANK_PASSWORD"])

	timeout := 10000.0
	GetTransferPage(page, &timeout)

	headers := sendFirstPhoneRequest(page)

	quit := make(chan bool, 2)
	go func() {
		for range quit {
			keepSession(page)
		}
	}()

	quit <- false

	respChan := make(chan string, 50)
	tasksChan := make(chan int, 50)
	wg := sync.WaitGroup{}
	numbers := make([]string, 1)
	for i := 0; i < 1; i++ {
		transfer := new(Transfer)
		wg.Add(1)
		go askForPhone(&url, &headers, client, tasksChan, respChan, &wg, i, transfer)
		handleError("error while reading POST response: ", &err)
	}

	for num := range numbers {
		tasksChan <- num
	}
	fmt.Println("[main] Wrote 1 task")

	close(tasksChan)

	go func() {
		wg.Wait()
		close(respChan)
	}()

	for respBody := range respChan {
		fmt.Println(respBody)
	}
	//fmt.Println("Sleeping")
	//time.Sleep(5 * time.Minute)
	quit <- true
	close(quit)
	fmt.Println("Quit channel closed")
	fmt.Println("[main] Main stopped")

	if err = browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
}
