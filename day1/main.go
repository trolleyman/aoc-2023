package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
)

func getinput(path string) ([]string, error) {
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

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getfilename() (string, error) {
	switch len(os.Args) {
	case 1:
		return "./input.txt", nil
	case 2:
		return os.Args[1], nil
	default:
		return "", errors.New(fmt.Sprintf("Invalid number of arguments: %v", len(os.Args)))
	}
}

func getfirstlastdigits(line string) (int, int, error) {
	var first *int = nil
	var last *int = nil

	for _, c := range line {
		i, err := strconv.Atoi(string(c))
		if err != nil {
			continue
		}
		last = &i
		if first == nil {
			first = &i
			continue
		}
	}

	if first == nil {
		return 0, 0, errors.New(fmt.Sprintf("First digit not found in %#v", line))
	}
	if last == nil {
		return 0, 0, errors.New(fmt.Sprintf("Last digit not found in %#v", line))
	}
	return *first, *last, nil
}

func run() error {
	filename, err := getfilename()
	if err != nil {
		return err
	}

	lines, err := getinput(filename)
	if err != nil {
		return err
	}

	sum := 0
	for _, line := range lines {
		first, last, err := getfirstlastdigits(line)
		if err != nil {
			return err
		}
		sum += first*10 + last
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
