package main_test

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"phone_numbers_checker/checker"
	"testing"
)

type transfer struct {
	Transfer map[string]interface{} `json:"transfer"`
}

var (
	headless    = true
	numsChecker = checker.Checker{
		Headless:      &headless,
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
	numbers := make(chan checker.Number, 500000)
	numsChecker.GetNumbers(numbers)
	t.Log(len(numbers))
	go close(numbers)
	numsChecker.SaveRelevantNumbers(numbers)
	numsChecker.DeleteInputFiles()
}
