package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type AshRock bool

func (ar AshRock) isRock() bool {
	return ar == Rock
}

func (ar AshRock) isAsh() bool {
	return ar == Ash
}

func (ar AshRock) String() string {
	if ar.isRock() {
		return "#"
	} else {
		return "."
	}
}

const (
	Ash  AshRock = false
	Rock AshRock = true
)

type Pattern [][]AshRock

func (p Pattern) String() string {
	result := ""
	for i, row := range p {
		if i > 0 {
			result += "\n"
		}
		for _, ar := range row {
			result += ar.String()
		}
	}
	return result
}

func getInput(path string) ([]Pattern, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var patterns []Pattern
	var currentPattern Pattern
	for y := 0; scanner.Scan(); y++ {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			if currentPattern != nil {
				patterns = append(patterns, currentPattern)
				currentPattern = nil
			}
			continue
		}

		var row []AshRock
		for _, r := range line {
			var ar AshRock
			switch r {
			case '#':
				ar = Rock
			case '.':
				ar = Ash
			default:
				return nil, fmt.Errorf("invalid rune %c", r)
			}
			row = append(row, ar)
		}
		currentPattern = append(currentPattern, row)
	}
	if currentPattern != nil {
		patterns = append(patterns, currentPattern)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return patterns, nil
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
		return Args{}, fmt.Errorf("invalid part number %#v. Expected 1/2", os.Args[1])
	}
	return Args{Part: part, InputPath: os.Args[2]}, nil
}

func run() error {
	args, err := parseArgs()
	if err != nil {
		return err
	}
	fmt.Printf("Args: %+v\n", args)

	patterns, err := getInput(args.InputPath)
	if err != nil {
		return err
	}

	for i, pattern := range patterns {
		horizontalReflection := pattern.getHorizontalReflection()
		verticalReflection := pattern.getVerticalReflection()
		if i > 0 {
			fmt.Println("")
		}
		printHorizontalNumberIndicator(pattern)
		printPatternHorizontalReflectionIndicator(pattern, horizontalReflection)
		for i, row := range pattern {
			fmt.Printf("%v", (i+1)%10)
			printPatternVerticalReflectionIndicator(verticalReflection, i)
			for _, ar := range row {
				fmt.Printf("%v", ar)
			}
			printPatternVerticalReflectionIndicator(verticalReflection, i)
			fmt.Printf("%v\n", (i+1)%10)
		}
		printPatternHorizontalReflectionIndicator(pattern, horizontalReflection)
		printHorizontalNumberIndicator(pattern)
	}

	// TODO
	// if args.Part == 2 {
	// }

	return nil
}

func (p Pattern) isValidHorizontalReflection(row []AshRock, rowI int, reflectionI int) bool {
	startI := (rowI - len(row)/2) * 2
	if startI < 0 {
		startI = 0
	}
	for i := 0; i < reflectionI; i++ {

	}
	return true
}

func (p Pattern) getHorizontalReflection() int {
	for reflectionI := 1; reflectionI < len(p[0]); reflectionI++ {
		validReflection := true
		for rowI, row := range p {
			if !p.isValidHorizontalReflection(row, rowI, reflectionI) {
				validReflection = false
				break
			}
		}
		if validReflection {
			return reflectionI
		}
	}
	return 0
}

func (p Pattern) getVerticalReflection() int {
	return 3
}

func printHorizontalNumberIndicator(pattern Pattern) {
	fmt.Printf("  ")
	for i := range pattern[0] {
		fmt.Printf("%v", (i+1)%10)
	}
	fmt.Printf("  \n")
}

func printPatternVerticalReflectionIndicator(verticalReflection int, i int) {
	if verticalReflection > 0 && i == verticalReflection-1 {
		fmt.Printf("v")
	} else if verticalReflection > 0 && i == verticalReflection {
		fmt.Printf("^")
	} else {
		fmt.Printf(" ")
	}
}

func printPatternHorizontalReflectionIndicator(pattern Pattern, horizontalReflection int) {
	fmt.Printf("  ")
	for i := range pattern[0] {
		if horizontalReflection > 0 && i == horizontalReflection-1 {
			fmt.Printf(">")
		} else if horizontalReflection > 0 && i == horizontalReflection {
			fmt.Printf("<")
		} else {
			fmt.Printf(" ")
		}
	}
	fmt.Printf("  \n")
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
