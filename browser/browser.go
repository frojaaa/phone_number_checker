package browser

import (
	"fmt"
	"github.com/playwright-community/playwright-go"
	"log"
	"phone_numbers_checker/errors"
	"strings"
)

type loginFields struct {
	login    playwright.ElementHandle
	password playwright.ElementHandle
}

type transferFields struct {
	phone       playwright.ElementHandle
	transferSum playwright.ElementHandle
	destination playwright.ElementHandle
}

type requestData struct {
	PayerAccount string `json:"payerAccount"`
}

func GetBrowserPage(browser playwright.Browser) playwright.Page {
	page, err := browser.NewPage()
	errors.HandleError("Can't  open page", &err)
	page.On("response", func(response playwright.Response) {
		if response.URL() == "https://ib.rencredit.ru/rencredit.server.portal.app/rest/private/transfers/internal/register" {
			fmt.Printf("<< %v %s\n", response.Status(), response.URL())
		}
	})
	return page
}

func findLoginFields(page playwright.Page) loginFields {
	usernameInput, err := page.QuerySelector("#username")
	errors.HandleError("Can't select username input", &err)

	passwordInput, err := page.QuerySelector("#password")
	errors.HandleError("Can't select password input", &err)

	loginFields := loginFields{
		login:    usernameInput,
		password: passwordInput,
	}

	return loginFields
}

func findTransferFields(page playwright.Page) transferFields {
	phoneInput, err := page.WaitForSelector("#destinationPhoneNumber")
	errors.HandleError("Can't select phone input", &err)

	transferSumInput, err := page.WaitForSelector("#amount")
	errors.HandleError("Can't select transfer input", &err)
	err = transferSumInput.Fill("1")
	errors.HandleError("Error while filling transfer sum", &err)
	destinationInput, err := page.WaitForSelector("#select2-paymentPurposeCode-container")
	errors.HandleError("Can't select destination input", &err)

	fields := transferFields{
		phone:       phoneInput,
		transferSum: transferSumInput,
		destination: destinationInput,
	}

	return fields
}

func fillLoginFields(loginFields loginFields, login string, password string) {
	log.Println(len(login), password)
	err := loginFields.login.Fill(login)
	errors.HandleError("error while filling username", &err)

	err = loginFields.password.Fill(password)
	errors.HandleError("error while filling password", &err)
}

func fillTransferFields(page playwright.Page, fields transferFields) {
	err := fields.phone.Fill("9000000000")
	errors.HandleError("Error while filling phone", &err)

	err = fields.transferSum.Fill("")
	err = fields.transferSum.Fill("1")
	errors.HandleError("Error while filling transfer sum", &err)

	err = fields.destination.Focus()
	err = fields.destination.Click()
	errors.HandleError("Error while click on destination input", &err)

	gift, err := page.WaitForSelector("ul.select2-results__options:nth-child(2) > li:nth-child(1)")
	errors.HandleError("Error while selecting gift button", &err)
	err = gift.Click()
	errors.HandleError("Error while clicking on gift button", &err)
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
	errors.HandleError("Can't select login button", &err)

	timeout := 10000.0
	err = loginBtn.Click(playwright.ElementHandleClickOptions{Timeout: &timeout})
	errors.HandleError("Error while login to account: ", &err)
}

func GetTransferPage(page playwright.Page, timeout *float64) {
	stateVisible := playwright.WaitForSelectorState("visible")
	selector, err := page.WaitForSelector(".header-user-menu__profile", playwright.PageWaitForSelectorOptions{
		State:   &stateVisible,
		Timeout: timeout,
	})
	errors.HandleError(selector.String(), &err)
	fmt.Println("Visit transfers")
	_, err = page.Goto("https://ib.rencredit.ru/#/transfers/internal", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	errors.HandleError("Error while getting page", &err)
	phoneBtn, err := page.WaitForSelector("#paymentMethod2", playwright.PageWaitForSelectorOptions{State: &stateVisible})
	errors.HandleError("Can't select phone button", &err)
	err = phoneBtn.Click()
	errors.HandleError("Can't click phone button", &err)
}

func SendFirstPhoneRequest(page playwright.Page) (map[string]string, string) {
	fields := findTransferFields(page)
	fillTransferFields(page, fields)

	primaryBtn, err := page.WaitForSelector("#primary-button")
	errors.HandleError("Error while selecting primary button", &err)

	err = primaryBtn.Click()
	errors.HandleError("Error while clicking on primary button", &err)

	response := page.WaitForResponse("https://ib.rencredit.ru/rencredit.server.portal.app/rest/private/transfers/internal/register")
	headers, err := response.Request().AllHeaders()
	reqData, err := response.Request().PostData()
	errors.HandleError("error getting request data: ", &err)
	account := strings.Split(strings.Split(reqData, "&")[0], "=")[1]
	errors.HandleError("Error while getting headers", &err)
	return headers, account
}

func KeepSession(page playwright.Page) {
	for {
		fmt.Println("Keep session")
		timeout := 0.0
		button, err := page.WaitForSelector("#button-button", playwright.PageWaitForSelectorOptions{
			Timeout: &timeout,
		})
		if err != nil {
			log.Println("Error while waiting for session button: ", err)
			return
		}
		err = button.Click()
		errors.HandleError("Error while clicking session button: ", &err)
		fmt.Println("Button clicked")
	}
}
