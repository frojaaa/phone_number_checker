package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"phone_numbers_checker/checker"
	//"phone_numbers_checker/errors"
	//"strconv"
)

func RunChecker(c *gin.Context) {
	var checkerInput checker.Checker
	if err := c.ShouldBindJSON(&checkerInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(checkerInput.NumWorkers)
	numbersChecker := checker.Checker{
		NumWorkers:    checkerInput.NumWorkers,
		LkLogin:       checkerInput.LkLogin,
		LkPassword:    checkerInput.LkPassword,
		BotToken:      checkerInput.BotToken,
		TgUserID:      checkerInput.TgUserID,
		InputFileDir:  checkerInput.InputFileDir,
		OutputFileDir: checkerInput.OutputFileDir,
	}
	go numbersChecker.Run()
	c.JSON(http.StatusOK, gin.H{
		"message": "Checker has began his work",
	})
}
