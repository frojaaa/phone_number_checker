package checker

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/playwright-community/playwright-go"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"phone_numbers_checker/browser"
	"phone_numbers_checker/errors"
	"strings"
	"sync"
	"time"
)

type Checker struct {
	Headless      *bool
	NumWorkers    int
	InputFileDir  string
	OutputFileDir string
}

type Transfer struct {
	Transfer map[string]interface{} `json:"transfer"`
}

type Number struct {
	Value    string
	FileName string
}

func fileNameWithoutExt(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}

func generateFormData() url.Values {
	form := url.Values{}
	form.Add("payerAccount", "589079858")
	form.Add("paymentMethod", "BY_PHONE")
	form.Add("destinationPhoneNumber", "79181185688")
	form.Add("amount", "1")
	form.Add("paymentPurposeCode", "GIFT")
	return form
}

func makeRequest(url *string, form url.Values, headers *map[string]string) (*http.Request, error) {
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

func askForPhone(url *string, headers *map[string]string, client *http.Client, tasks chan Number, results chan Number, wg *sync.WaitGroup, id int, target *Transfer) {
	defer wg.Done()
	for num := range tasks {
		form := generateFormData()
		//form.Add("payerAccount", "589079858")
		//form.Add("paymentMethod", "BY_PHONE")
		//form.Add("destinationPhoneNumber", "79181185688")
		//form.Add("amount", "1")
		//form.Add("paymentPurposeCode", "GIFT")
		request, err := makeRequest(url, form, headers)
		errors.HandleError("error doing request: ", &err)
		resp, err := getResponse(request, client)
		errors.HandleError("error while request: ", &err)
		if resp.StatusCode == 200 {
			err = getJSON(resp, target)
			//bankName := target.Transfer["payeeBankName"].(string)
			//bodyBytes, err := io.ReadAll(resp.Body)
			//HandleError("Error while reading response bytes: ", &err)
			//bodyString := string(bodyBytes)
			fmt.Printf("[worker %d] Worker Sending result of task %s\n", id, num)
			results <- num
		} else {
			fmt.Println()
		}
	}
}

func (c Checker) GetNumbers(tasks chan Number) {
	items, err := os.ReadDir(c.InputFileDir)
	errors.HandleError("Error while reading input directory: ", &err)
	for _, item := range items {
		file, err := os.Open(c.InputFileDir + string(os.PathSeparator) + item.Name())
		errors.HandleError("Error while opening file: ", &err)

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			//fmt.Println(scanner.Text())
			tasks <- Number{
				Value:    scanner.Text(),
				FileName: item.Name(),
			}
		}
		fmt.Println("Wrote number")
		if err = scanner.Err(); err != nil {
			log.Fatalf("Error while reading file: %s", err)
		}
		file.Close()
	}
}

func (c Checker) SaveRelevantNumbers(numbers chan Number) {
	items, err := os.ReadDir(c.InputFileDir)
	files := make(map[string]*os.File)
	for _, file := range items {
		files[file.Name()], err = os.Create(c.OutputFileDir + string(os.PathSeparator) + fileNameWithoutExt(file.Name()) + " СБЕР" + ".txt")
	}
	defer func() {
		for _, file := range files {
			file.Close()
		}
	}()
	errors.HandleError("Error while reading input directory: ", &err)
	go close(numbers)
	for num := range numbers {
		file := files[num.FileName]
		fmt.Println(num.Value)
		errors.HandleError("Error while opening file: ", &err)
		_, err = fmt.Fprintln(file, num.Value)
		errors.HandleError("Error while writing line to file", &err)
		file.Close()
	}
	fmt.Println("All numbers are saved")
	fmt.Println("respChan closed")
}

func (c Checker) DeleteInputFiles() {
	items, err := os.ReadDir(c.InputFileDir)
	errors.HandleError("Error while reading input directory: ", &err)
	for _, item := range items {
		err = os.Remove(c.InputFileDir + string(os.PathSeparator) + item.Name())
		errors.HandleError("Can't delete input files: ", &err)
	}
	fmt.Println("Input files deleted")
}

func (c Checker) Run() {
	env, err := godotenv.Read()
	errors.HandleError("Error while loading .env: ", &err)
	transport := http.Transport{DisableKeepAlives: true, MaxIdleConns: 200}
	client := &http.Client{Timeout: 30 * time.Second, Transport: &transport}
	requestUrl := "https://ib.rencredit.ru/rencredit.server.portal.app/rest/private/transfers/internal/register"
	pw, err := playwright.Run()
	errors.HandleError("Unable to run playwright", &err)
	driver, err := pw.Firefox.Launch(playwright.BrowserTypeLaunchOptions{Headless: c.Headless})
	errors.HandleError("Can't launch Chromium", &err)
	fmt.Println(env["BANK_LOGIN"], env["BANK_PASSWORD"])
	page := browser.GetBrowserPage(driver)

	browser.LoginToAccount(page, env["BANK_LOGIN"], env["BANK_PASSWORD"])

	timeout := 10000.0
	browser.GetTransferPage(page, &timeout)

	headers := browser.SendFirstPhoneRequest(page)

	quit := make(chan bool, 1)
	go func() {
		for range quit {
			browser.KeepSession(page)
		}
	}()

	respChan := make(chan Number, 50000)
	tasksChan := make(chan Number, 50000)
	var wg sync.WaitGroup
	for i := 0; i < c.NumWorkers; i++ {
		transfer := new(Transfer)
		wg.Add(1)
		go askForPhone(&requestUrl, &headers, client, tasksChan, respChan, &wg, i, transfer)
		errors.HandleError("error while reading POST response: ", &err)
	}

	c.GetNumbers(tasksChan)
	fmt.Printf("[checker] Wrote tasks\n")
	close(tasksChan)

	wg.Wait()

	quit <- true
	close(quit)
	fmt.Println("Quit channel closed")
	c.SaveRelevantNumbers(respChan)
	c.DeleteInputFiles()
	fmt.Println("[checker] Checker.Run() stopped")

	if err = driver.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
}
