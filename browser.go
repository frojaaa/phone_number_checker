package main

import (
	"fmt"
	"github.com/playwright-community/playwright-go"
	"log"
)

type loginFields struct {
	login    playwright.ElementHandle
	password playwright.ElementHandle
}

type tranferFields struct {
	phone       playwright.ElementHandle
	transferSum playwright.ElementHandle
	destination playwright.ElementHandle
}

func GetBrowserPage(browser playwright.Browser) playwright.Page {
	page, err := browser.NewPage()
	handleError("Can't  open page", &err)
	page.On("response", func(response playwright.Response) {
		if response.URL() == "https://ib.rencredit.ru/rencredit.server.portal.app/rest/private/transfers/internal/register" {
			fmt.Printf("<< %v %s\n", response.Status(), response.URL())
		}
	})
	return page
}

func findLoginFields(page playwright.Page) loginFields {
	usernameInput, err := page.QuerySelector("#username")
	handleError("Can't select username input", &err)

	passwordInput, err := page.QuerySelector("#password")
	handleError("Can't select password input", &err)

	loginFields := loginFields{
		login:    usernameInput,
		password: passwordInput,
	}

	return loginFields
}

func findTransferFields(page playwright.Page) tranferFields {
	phoneInput, err := page.WaitForSelector("#destinationPhoneNumber")
	handleError("Can't select phone input", &err)

	transferSumInput, err := page.WaitForSelector("#amount")
	handleError("Can't select transfer input", &err)
	err = transferSumInput.Fill("1")
	handleError("Error while filling transfer sum", &err)
	destinationInput, err := page.WaitForSelector("#select2-paymentPurposeCode-container")
	handleError("Can't select destination input", &err)

	fields := tranferFields{
		phone:       phoneInput,
		transferSum: transferSumInput,
		destination: destinationInput,
	}

	return fields
}

func fillLoginFields(loginFields loginFields, login string, password string) {
	log.Println(len(login), password)
	err := loginFields.login.Fill(login)
	handleError("error while filling username", &err)

	err = loginFields.password.Fill(password)
	handleError("error while filling password", &err)
}

func fillTransferFields(page playwright.Page, fields tranferFields) {
	err := fields.phone.Fill("9000000000")
	handleError("Error while filling phone", &err)

	err = fields.transferSum.Fill("1")
	handleError("Error while filling transfer sum", &err)

	err = fields.destination.Focus()
	err = fields.destination.Click()
	handleError("Error while click on destination input", &err)

	gift, err := page.WaitForSelector("ul.select2-results__options:nth-child(2) > li:nth-child(1)")
	handleError("Error while selecting gift button", &err)
	err = gift.Click()
	handleError("Error while clicking on gift button", &err)
}

func LoginToAccount(page playwright.Page, login string, password string) {
	if _, err := page.Goto("https://ib.rencredit.ru/#/login", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	}); err != nil {
		log.Fatalf("could not goto: %v", err)
	}
	loginFields := findLoginFields(page)
	fillLoginFields(loginFields, login, password)

	loginBtn, err := page.QuerySelector("#button-button")
	handleError("Can't select login button", &err)

	timeout := 10000.0
	err = loginBtn.Click(playwright.ElementHandleClickOptions{Timeout: &timeout})
	handleError("Error while login to account: ", &err)
}

func GetTransferPage(page playwright.Page, timeout *float64) {
	stateVisible := playwright.WaitForSelectorState("visible")
	selector, err := page.WaitForSelector(".header-user-menu__profile", playwright.PageWaitForSelectorOptions{
		State: &stateVisible,
		//Timeout: timeout,
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
}

func sendFirstPhoneRequest(page playwright.Page) map[string]string {
	fields := findTransferFields(page)
	fillTransferFields(page, fields)

	primaryBtn, err := page.WaitForSelector("#primary-button")
	handleError("Error while selecting primary button", &err)

	err = primaryBtn.Click()
	handleError("Error while clicking on primary button", &err)

	response := page.WaitForResponse("https://ib.rencredit.ru/rencredit.server.portal.app/rest/private/transfers/internal/register")
	headers, err := response.Request().AllHeaders()
	handleError("Error while getting headers", &err)
	return headers
}
