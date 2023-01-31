package main_test

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
)

type transfer struct {
	Transfer map[string]interface{} `json:"transfer"`
}

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
