package main

import (
	"bytes"
	"fmt"
	"github.com/playwright-community/playwright-go"
	"io"
	"log"
	"net/http"
)

func handleError(helpText string, err *error) {
	if *err != nil {
		log.Fatalf("%s: %v", helpText, *err)
	}
}

func askForPhone(url *string, data *[]byte, headers *map[string]string, client *http.Client) {
	request, err := http.NewRequest("POST", *url, bytes.NewBuffer(*data))
	for key, value := range *headers {
		request.Header.Set(key, value)
	}
	fmt.Println(request.Header)
	resp, err := client.Do(request)
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println(resp.Status)
	log.Println(bodyString)
}

func main() {
	client := &http.Client{}
	//headers := Headers{}
	////cookies := Cookies{}
	//err := parseHeaders(&headers)
	//if err != nil {
	//	log.Fatal(err)
	//}
	////err = parseCookies(&cookies)
	url := "https://ib.rencredit.ru/rencredit.server.portal.app/rest/private/transfers/internal/register"
	var data []byte
	copy(data[:], "payerAccount=505419565&paymentMethod=BY_PHONE&destinationPhoneNumber=79005200846&amount=1&paymentPurposeCode=GIFT")
	//request, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	//for _, header := range headers.Headers {
	//	request.Header.Set(header.Name, header.Value)
	//}
	////for _, cookie := range cookies.Cookies {
	////	request.AddCookie(&http.Cookie{Name: cookie.Name, Value: cookie.Value})
	////}
	//response, err := client.Do(request)
	//defer response.Body.Close()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//bodyBytes, err := io.ReadAll(response.Body)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//bodyString := string(bodyBytes)
	//fmt.Println(response.Status)
	//log.Println(bodyString)

	pw, err := playwright.Run()
	handleError("Unable to run playwright", &err)
	headless := true
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless})
	handleError("Can't launch Chromium", &err)
	page, err := browser.NewPage()
	handleError("Can't  open page", &err)
	page.On("response", func(response playwright.Response) {
		if response.URL() == "https://ib.rencredit.ru/rencredit.server.portal.app/rest/private/transfers/internal/register" {
			fmt.Printf("<< %v %s\n", response.Status(), response.URL())
		}
	})

	if _, err = page.Goto("https://ib.rencredit.ru/#/login", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	}); err != nil {
		log.Fatalf("could not goto: %v", err)
	}
	usernameInput, err := page.QuerySelector("#username")
	handleError("Can't select username input", &err)
	passwordInput, err := page.QuerySelector("#password")
	login := "encredit7989"
	handleError("Can't select password input", &err)
	password := "6P-.Q8$_i6$z_2A"
	err = usernameInput.Fill(login)
	handleError("error while filling username", &err)
	err = passwordInput.Fill(password)
	handleError("error while filling password", &err)
	loginBtn, err := page.QuerySelector("#button-button")
	handleError("Can't select login button", &err)
	timeout := 10000.0
	err = loginBtn.Click(playwright.ElementHandleClickOptions{Timeout: &timeout})
	handleError("Can't click login button", &err)
	stateVisible := playwright.WaitForSelectorState("visible")
	selector, err := page.WaitForSelector(".header-user-menu__profile", playwright.PageWaitForSelectorOptions{
		State:   &stateVisible,
		Timeout: &timeout,
	})
	handleError(selector.String(), &err)
	fmt.Println("Visit transfers")
	_, err = page.Goto("https://ib.rencredit.ru/#/transfers/internal", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	handleError("Error while getting page", &err)
	phoneBtn, err := page.WaitForSelector("#paymentMethod2", playwright.PageWaitForSelectorOptions{State: &stateVisible})
	handleError("Can't select phone button", &err)
	err = phoneBtn.Click()
	handleError("Can't click phone button", &err)
	phoneInput, err := page.WaitForSelector("#destinationPhoneNumber")
	handleError("Can't select phone input", &err)
	err = phoneInput.Fill("9005200846")
	handleError("Error while filling phone", &err)
	transferSumInput, err := page.WaitForSelector("#amount")
	handleError("Can't select transfer input", &err)
	err = transferSumInput.Fill("1")
	handleError("Error while filling transfer sum", &err)
	destinationInput, err := page.WaitForSelector("#select2-paymentPurposeCode-container")
	handleError("Can't select destination input", &err)
	err = destinationInput.Focus()
	err = destinationInput.Click()
	handleError("Error while click on destination input", &err)
	gift, err := page.WaitForSelector("ul.select2-results__options:nth-child(2) > li:nth-child(1)")
	handleError("Error while selecting gift button", &err)
	err = gift.Click()
	handleError("Error while clicking on gift button", &err)
	primaryBtn, err := page.WaitForSelector("#primary-button")
	handleError("Error while selecting primary button", &err)
	err = primaryBtn.Click()
	handleError("Error while clicking on primary button", &err)
	response := page.WaitForResponse("https://ib.rencredit.ru/rencredit.server.portal.app/rest/private/transfers/internal/register")
	headers, err := response.Request().AllHeaders()
	handleError("Error while getting headers", &err)
	go askForPhone(&url, &data, &headers, client)
	go askForPhone(&url, &data, &headers, client)

	if err = browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
}
