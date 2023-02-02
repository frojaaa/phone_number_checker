package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"phone_numbers_checker/checker"
	"phone_numbers_checker/errors"
	"strconv"
)

type CheckerInput struct {
	InputFileDir  string `json:"inputFileDir" binding:"required"`
	OutputFileDir string `json:"outputFileDir" binding:"required"`
}

func RunChecker(c *gin.Context) {
	headless := true
	numWorkers, err := strconv.ParseInt(c.Query("numWorkers"), 10, 32)
	errors.HandleError("Error while parsing query param numWorkers: ", &err)
	numbersChecker := checker.Checker{
		Headless:      &headless,
		NumWorkers:    int(numWorkers),
		InputFileDir:  "./input/",
		OutputFileDir: "./output/",
	}
	go numbersChecker.Run()
	c.JSON(http.StatusOK, gin.H{
		"message": "Checker has began his work",
	})
}
