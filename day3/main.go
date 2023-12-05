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
		return nil, errors.New("zero-length engine")
	}

	for i, line := range engine {
		if len(line) != len(engine[0]) {
			return nil, fmt.Errorf("jagged engine (line %v is length %v, line 1 is length %v)", i, len(line), len(engine[0]))
		}
	}

	return engine, nil
}

type Args struct {
	Part      int
	InputPath string
}

func parseArgs() (Args, error) {
	switch len(os.Args) {
	case 3:
		break
	default:
		return Args{}, fmt.Errorf("invalid arguments. Expected %v <part> <inputPath>", os.Args[0])
	}
	var part int
	switch os.Args[1] {
	case "1":
		part = 1
	case "2":
		part = 2
	default:
		return Args{}, fmt.Errorf("invalid part. Expected 1/2, got %#v", os.Args[1])
	}
	return Args{Part: part, InputPath: os.Args[2]}, nil
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

func isValidCoordinate(engine Engine, x int, y int) bool {
	return x >= 0 && y >= 0 && y < len(engine) && x < len(engine[y])
}

func isDigitAt(engine Engine, x int, y int) bool {
	if !isValidCoordinate(engine, x, y) {
		return false
	}
	c := getRune(engine, x, y)
	return isDigit(c)
}

func isSymbol(engine Engine, x int, y int) bool {
	if !isValidCoordinate(engine, x, y) {
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

func getGearNumberAt(engine Engine, x int, y int) int {
	// fmt.Printf("@%v,%v # ", x, y)
	for ; x < len(engine[y]) && isDigitAt(engine, x, y); x++ {
		// fmt.Printf("%v ", x)
	}
	x -= 1
	// fmt.Printf("# %v ", x)
	gearNumber := 0
	multiplier := 1
	for ; x >= 0; x-- {
		c := getRune(engine, x, y)
		// fmt.Printf("| %v,%v,%v ", x, y, c)
		digit, ok := parseDigit(c)
		if !ok {
			break
		}
		gearNumber += digit * multiplier
		multiplier *= 10
	}
	// fmt.Printf("@ ")
	return gearNumber
}

func getGearRatio(engine Engine, x int, y int) (int, bool) {
	var gearNumbers []int
	for offsetY := -1; offsetY <= 1; offsetY++ {
		for offsetX := -1; offsetX <= 1; offsetX++ {
			if offsetX == 0 && offsetY == 0 {
				continue
			}
			if isDigitAt(engine, x+offsetX, y+offsetY) {
				gearNumbers = append(gearNumbers, getGearNumberAt(engine, x+offsetX, y+offsetY))
				if offsetY == 0 {
					continue
				}
				if (offsetX == -1 && isDigitAt(engine, x, y+offsetY)) || (offsetX == 0 && isDigitAt(engine, x+1, y+offsetY)) {
					break
				}
			}
		}
	}
	switch len(gearNumbers) {
	case 0:
		return 0, false
	case 1:
		fmt.Printf("%v=", gearNumbers[0])
		return 0, false
	case 2:
		fmt.Printf("%v*%v=", gearNumbers[0], gearNumbers[1])
		return gearNumbers[0] * gearNumbers[1], true
	}
	panic(fmt.Sprintf("> 2 gear numbers: %#v", gearNumbers))
}

func getGearRatios(engine Engine) []int {
	var gearRatios []int
	for y, line := range engine {
		fmt.Printf("= ")
		for x := 0; x < len(line); x++ {
			if rune(line[x]) == '*' {
				gearRatio, ok := getGearRatio(engine, x, y)
				if ok {
					fmt.Printf("%v ", gearRatio)
					gearRatios = append(gearRatios, gearRatio)
				} else {
					fmt.Printf("~ ")
				}
			}
		}
		fmt.Printf("=\n")
	}
	fmt.Printf("\n")
	return gearRatios
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

	switch args.Part {
	case 1:
		numbers := getPartNumbers(engine)
		sum := 0
		for _, number := range numbers {
			sum += number
		}
		fmt.Printf("%v\n", sum)

	case 2:
		gearRatios := getGearRatios(engine)
		sum := 0
		for _, gearRatio := range gearRatios {
			sum += gearRatio
		}
		fmt.Printf("%v\n", sum)

	default:
		return fmt.Errorf("unknown part %v", args.Part)
	}

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
