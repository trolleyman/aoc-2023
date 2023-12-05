package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func makeDigitsMap(useWords bool) map[string]int {
	m := make(map[string]int)
	if useWords {
		m["one"] = 1
		m["two"] = 2
		m["three"] = 3
		m["four"] = 4
		m["five"] = 5
		m["six"] = 6
		m["seven"] = 7
		m["eight"] = 8
		m["nine"] = 9
	}
	for i := 1; i < 10; i++ {
		m[fmt.Sprintf("%v", i)] = i
	}
	return m
}

var digits = makeDigitsMap(false)
var digitsWithWords = makeDigitsMap(true)

func getInput(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// fmt.Println(line)
		if line != "" {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func parseArgs() (string, bool, error) {
	if len(os.Args) != 3 {
		return "", false, fmt.Errorf("invalid number of arguments (expected 3, got %v)", len(os.Args))
	}
	useWords, err := strconv.ParseBool(os.Args[2])
	if err != nil {
		return "", false, err
	}
	return os.Args[1], useWords, nil
}

func getDigits(useWords bool) map[string]int {
	if useWords {
		return digits
	}
	return digitsWithWords
}

func getFirstDigit(line string, useWords bool) (int, error) {
	for i := 0; i < len(line); i++ {
		lineSlice := line[i:]
		for possibleDigit, digit := range getDigits(useWords) {
			if strings.HasPrefix(lineSlice, possibleDigit) {
				return digit, nil
			}
		}
	}
	return 0, fmt.Errorf("first digit not found in %#v", line)
}

func getLastDigit(line string, useWords bool) (int, error) {
	for i := len(line) - 1; i >= 0; i-- {
		lineSlice := line[i:]
		for possibleDigit, digit := range getDigits(useWords) {
			if strings.HasPrefix(lineSlice, possibleDigit) {
				return digit, nil
			}
		}
	}
	return 0, fmt.Errorf("last digit not found in %#v", line)
}

func getFirstLastDigits(line string, useWords bool) (int, int, error) {
	first, err := getFirstDigit(line, useWords)
	if err != nil {
		return 0, 0, err
	}
	last, err := getLastDigit(line, useWords)
	if err != nil {
		return 0, 0, err
	}
	return first, last, nil
}

func run() error {
	filepath, useWords, err := parseArgs()
	if err != nil {
		return err
	}

	lines, err := getInput(filepath)
	if err != nil {
		return err
	}

	sum := 0
	for _, line := range lines {
		first, last, err := getFirstLastDigits(line, useWords)
		if err != nil {
			return err
		}
		lineSum := first*10 + last
		sum += lineSum
		// fmt.Printf("line:%#v, first:%v, last:%v, lineSum:%v, sum:%v\n", line, first, last, lineSum, sum)
	}
	fmt.Println(sum)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
