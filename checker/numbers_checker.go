package checker

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/playwright-community/playwright-go"
	tele "gopkg.in/telebot.v3"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"phone_numbers_checker/bot"
	"phone_numbers_checker/browser"
	"phone_numbers_checker/errors"
	"strconv"
	"strings"
	"sync"
	"time"
)

var handled = 0

type Checker struct {
	headless      *bool
	NumWorkers    int64  `json:"numWorkers" form:"numWorkers" binding:"required"`
	LkLogin       string `json:"lkLogin"`
	LkPassword    string `json:"lkPassword"`
	BotToken      string `json:"botToken"`
	TgUserID      int64  `json:"tgUserID"`
	InputFileDir  string `json:"inputFileDir"`
	OutputFileDir string `json:"outputFileDir"`
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

func generateFormData(num string, account string) url.Values {
	form := url.Values{}
	form.Add("payerAccount", account)
	form.Add("paymentMethod", "BY_PHONE")
	form.Add("destinationPhoneNumber", num)
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

func askForPhone(url *string, headers *map[string]string, account string, client *http.Client,
	tasks chan Number, results chan Number, wg *sync.WaitGroup, id int64, tgBot *tele.Bot, userID int64, m *sync.Mutex) {
	template := fmt.Sprintf("Обработано %s", strconv.Itoa(handled))
	msg := bot.SendMessage(tgBot, userID, template)
	for num := range tasks {
		if handled != 0 && handled%50 == 0 {
			bot.EditMessage(tgBot, msg, fmt.Sprintf("Обработано %s", strconv.Itoa(handled)))
		}
		fmt.Println("[askForPhone] Checking ", num)
		form := generateFormData(num.Value, account)
		request, err := makeRequest(url, form, headers)
		errors.HandleError("error doing request: ", &err)
		resp, err := getResponse(request, client)
		errors.HandleError("error while request: ", &err)
		bodyBytes, err := io.ReadAll(resp.Body)
		errors.HandleError("error getting response: ", &err)
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
		if resp.StatusCode == 200 {
			//err = getJSON(resp, target)
			//bankName := target.Transfer["payeeBankName"].(string)
			//bodyBytes, err := io.ReadAll(resp.Body)
			//HandleError("Error while reading response bytes: ", &err)
			//bodyString := string(bodyBytes)
			//fmt.Printf("[worker %d] Worker Sending result of task %s\n", id, num)
			results <- num
		} else {
			fmt.Printf("[worker %d] Number irrelevant\n", id)
		}
		m.Lock()
		handled++
		m.Unlock()
	}
	wg.Done()
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
		if err = scanner.Err(); err != nil {
			log.Fatalf("Error while reading file: %s", err)
		}
		file.Close()
	}
}

func (c Checker) SaveRelevantNumbers(numbers chan Number, tgBot *tele.Bot, userId int64) {
	items, err := os.ReadDir(c.InputFileDir)
	files := make(map[string]*os.File)
	for _, file := range items {
		files[file.Name()], err = os.Create(c.OutputFileDir + string(os.PathSeparator) + fileNameWithoutExt(file.Name()) + " СБЕР" + ".txt")
	}
	defer func() {
		for _, file := range files {
			stat, err := file.Stat()
			errors.HandleError("error getting file stats: ", &err)
			size := stat.Size()
			if size != 0 {
				bot.SendDocument(tgBot, userId, file.Name())
			} else {
				bot.SendMessage(tgBot, userId, "Релевантные номера не найдены в файле "+file.Name())
			}
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
	transport := http.Transport{DisableKeepAlives: true, MaxIdleConns: 200}
	client := &http.Client{Timeout: 30 * time.Second, Transport: &transport}
	requestUrl := "https://ib.rencredit.ru/rencredit.server.portal.app/rest/private/transfers/internal/register"
	pw, err := playwright.Run()
	errors.HandleError("Unable to run playwright", &err)
	headless := true
	driver, err := pw.Firefox.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless})
	errors.HandleError("Can't launch Chromium", &err)
	page := browser.GetBrowserPage(driver)

	browser.LoginToAccount(page, c.LkLogin, c.LkPassword)

	timeout := 10000.0
	browser.GetTransferPage(page, &timeout)

	headers, account := browser.SendFirstPhoneRequest(page)

	quit := make(chan bool, 1)
	tgBot := bot.GetBot(c.BotToken)
	go func() {
		for range quit {
			browser.KeepSession(page)
		}
	}()
	go tgBot.Start()
	respChan := make(chan Number, 500000)
	tasksChan := make(chan Number, 500000)
	var wg sync.WaitGroup
	//transfer := new(Transfer)
	fmt.Println(c.NumWorkers)
	var i int64
	var m sync.Mutex
	for i = 0; i <= c.NumWorkers; i++ {
		wg.Add(1)
		go askForPhone(&requestUrl, &headers, account, client, tasksChan, respChan, &wg, i, tgBot, c.TgUserID, &m)
	}

	c.GetNumbers(tasksChan)
	fmt.Printf("[checker] Wrote tasks\n")
	close(tasksChan)

	wg.Wait()

	quit <- true
	close(quit)
	fmt.Println("Quit channel closed")
	c.SaveRelevantNumbers(respChan, tgBot, c.TgUserID)
	c.DeleteInputFiles()
	tgBot.Stop()
	fmt.Println("[checker] Checker.Run() stopped")

	if err = driver.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
}
