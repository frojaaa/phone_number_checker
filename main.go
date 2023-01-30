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
	"sync"
	"time"
)

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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln("error while closing response body: ", err)
		}
	}(resp.Body)
	return resp, err
}

func askForPhone(url *string, headers *map[string]string, client *http.Client, tasks chan int, results chan string, wg *sync.WaitGroup, id int) {
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
		//bodyBytes, err := io.ReadAll(resp.Body)
		//handleError("Error while reading response bytes: ", &err)
		//bodyString := string(bodyBytes)
		fmt.Printf("[worker %d] Worker Sending result of task %d\n", id, num)
		results <- resp.Status
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
	headless := true
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless})
	handleError("Can't launch Chromium", &err)
	fmt.Println(env["BANK_LOGIN"], env["BANK_PASSWORD"])
	page := GetBrowserPage(browser)

	LoginToAccount(page, env["BANK_LOGIN"], env["BANK_PASSWORD"])

	timeout := 10000.0
	GetTransferPage(page, &timeout)

	headers := sendFirstPhoneRequest(page)

	respChan := make(chan string, 50)
	tasksChan := make(chan int, 50)
	wg := sync.WaitGroup{}
	numbers := make([]string, 50)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go askForPhone(&url, &headers, client, tasksChan, respChan, &wg, i)
		handleError("error while reading POST response: ", &err)
	}

	for num := range numbers {
		tasksChan <- num
	}
	fmt.Println("[main] Wrote 50 tasks")

	close(tasksChan)

	go func() {
		wg.Wait()
		close(respChan)
	}()

	for respBody := range respChan {
		fmt.Println(respBody)
	}
	fmt.Println("[main] Main stopped")

	if err = browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
}
