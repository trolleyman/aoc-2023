package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

type Engine []string

func (e Engine) String() string {
	var s string
	for i, line := range e {
		if i != 0 {
			s += "\n"
		}
		s += line
	}
	return s
}

func getInput(path string) (Engine, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var engine Engine
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		engine = append(engine, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(engine) == 0 {
		return nil, errors.New("Zero-length engine")
	}

	for i, line := range engine {
		if len(line) != len(engine[0]) {
			return nil, errors.New(fmt.Sprintf("Jagged engine (line %v is length %v, line 1 is length %v)\n", i, len(line), len(engine[0])))
		}
	}

	return engine, nil
}

type Args struct {
	InputPath string
}

func parseArgs() (Args, error) {
	switch len(os.Args) {
	case 2:
		break
	default:
		return Args{}, errors.New(fmt.Sprintf("Invalid arguments. Expected %v <inputPath> [inputSet]", os.Args[0]))
	}
	return Args{InputPath: os.Args[1]}, nil
}

func isDigit(c rune) bool {
	_, ok := parseDigit(c)
	return ok
}

func parseDigit(c rune) (int, bool) {
	switch c {
	case '0':
		return 0, true
	case '1':
		return 1, true
	case '2':
		return 2, true
	case '3':
		return 3, true
	case '4':
		return 4, true
	case '5':
		return 5, true
	case '6':
		return 6, true
	case '7':
		return 7, true
	case '8':
		return 8, true
	case '9':
		return 9, true
	}
	return 0, false
}

func getRune(engine Engine, x int, y int) rune {
	return rune(engine[y][x])
}

func isSymbol(engine Engine, x int, y int) bool {
	if x < 0 || y < 0 || y >= len(engine) || x >= len(engine[y]) {
		return false
	}
	c := getRune(engine, x, y)
	return !isDigit(c) && c != '.'
}

func isAdjacentToSymbol(engine Engine, x int, y int) bool {
	for i := 0; i < 9; i++ {
		if i == 4 {
			continue
		}
		offsetX := (i % 3) - 1
		offsetY := (i / 3) - 1
		if isSymbol(engine, x+offsetX, y+offsetY) {
			return true
		}
	}
	return false
}

func getPartNumbers(engine Engine) []int {
	var numbers []int
	for y, line := range engine {
		fmt.Printf("= ")
		var numberBuffer int
		var numberBufferX int
		for x := 0; x <= len(line); x++ {
			var digit int
			var isDigit bool
			if x < len(line) {
				c := getRune(engine, x, y)
				digit, isDigit = parseDigit(c)
			}
			if isDigit {
				if numberBuffer == 0 {
					numberBufferX = x
				}
				numberBuffer *= 10
				numberBuffer += digit
			} else if numberBuffer != 0 {
				isPartNumber := false
				for numberX := numberBufferX; numberX < x && numberX < len(line); numberX++ {
					if isAdjacentToSymbol(engine, numberX, y) {
						isPartNumber = true
						break
					}
				}

				if !isPartNumber {
					fmt.Printf("~")
				}
				fmt.Printf("%v", numberBuffer)
				if !isPartNumber {
					fmt.Printf("~")
				}
				fmt.Printf(" ")
				if isPartNumber {
					numbers = append(numbers, numberBuffer)
				}

				numberBuffer = 0
				numberBufferX = 0
			}
		}
		fmt.Printf("=\n")
	}
	fmt.Printf("\n")
	return numbers
}

func run() error {
	args, err := parseArgs()
	if err != nil {
		return err
	}
	// fmt.Printf("Args: %+v\n", args)

	engine, err := getInput(args.InputPath)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n\n", engine)

	numbers := getPartNumbers(engine)
	sum := 0
	for _, number := range numbers {
		sum += number
	}
	fmt.Printf("%v\n", sum)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
