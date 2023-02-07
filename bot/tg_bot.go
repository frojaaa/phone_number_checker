package bot

import (
	tele "gopkg.in/telebot.v3"
	"phone_numbers_checker/errors"
	"strconv"
	"time"
)

type TgUser struct {
	ID int64
}

func (u TgUser) Recipient() string {
	return strconv.FormatInt(u.ID, 10)
}

func GetBot(token string) *tele.Bot {
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := tele.NewBot(pref)
	errors.HandleError("Error getting bot: ", &err)
	return b
}

func SendMessage(bot *tele.Bot, userID int64, message any) *tele.Message {
	u := TgUser{ID: userID}
	msg, err := bot.Send(u, message)
	errors.HandleError("Error while sending message to user: ", &err)
	return msg
}

func EditMessage(bot *tele.Bot, msg *tele.Message, text any) *tele.Message {
	msg, err := bot.Edit(msg, text)
	errors.HandleError("error editing message: ", &err)
	return msg
}

func SendDocument(bot *tele.Bot, userID int64, documentPath string) {
	u := TgUser{ID: userID}
	file := &tele.Document{
		File: tele.FromDisk(documentPath),
	}
	_, err := bot.Send(u, file)
	errors.HandleError("Error sending document to user: ", &err)
}
