package main_test

import (
	"bufio"
	"fmt"
	"os"
	"testing"
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
