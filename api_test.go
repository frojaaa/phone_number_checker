package main_test

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"phone_numbers_checker/bot"
	"phone_numbers_checker/checker"
	"testing"
	"time"
)

type transfer struct {
	Transfer map[string]interface{} `json:"transfer"`
}

var (
	numsChecker = checker.Checker{
		NumWorkers:    0,
		InputFileDir:  "./input",
		OutputFileDir: "./output",
	}
)

func BenchmarkIONumbers(b *testing.B) {
	f, err := os.Open("test.txt")
	var numbers []string
	if err != nil {
		b.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
		numbers = append(numbers, scanner.Text())
	}

	b.Log(len(numbers))
	if err = scanner.Err(); err != nil {
		b.Fatal(err)
	}
}

func TestJSONParsing(t *testing.T) {
	jsonFile, err := os.Open("transferTest.json")
	if err != nil {
		t.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	var target transfer
	err = json.Unmarshal(byteValue, &target)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(target.Transfer["payeeBankName"].(string))
}

func TestNumbersIO(t *testing.T) {
	tgBot := bot.GetBot("5874789025:AAHNGeRfMY2bOmGvXpK3ZWI58GQc8NMmDF0")
	go tgBot.Start()
	numbers := make(chan checker.Number, 500000)
	numsChecker.GetNumbers(numbers)
	t.Log(len(numbers))
	go close(numbers)
	numsChecker.SaveRelevantNumbers(numbers, tgBot, 994854069)
	numsChecker.DeleteInputFiles()
	tgBot.Stop()
}

func TestBot(t *testing.T) {
	tgBot := bot.GetBot("5874789025:AAHNGeRfMY2bOmGvXpK3ZWI58GQc8NMmDF0")
	go tgBot.Start()
	fmt.Println("sleep")
	time.Sleep(10 * time.Second)
	tgBot.Stop()
}
