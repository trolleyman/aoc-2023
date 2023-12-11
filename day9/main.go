package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parseNumbers(numbersString string) ([]int, error) {
	var numbers []int
	numbersString = strings.TrimSpace(numbersString)
	for _, numberString := range strings.Fields(numbersString) {
		numberString = strings.TrimSpace(numberString)
		number, err := strconv.Atoi(numberString)
		if err != nil {
			return nil, err
		}
		numbers = append(numbers, number)
	}
	return numbers, nil
}

func getInput(path string, multipleRouters bool) ([][]int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var histories [][]int
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		history, err := parseNumbers(line)
		if err != nil {
			return nil, err
		}
		histories = append(histories, history)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return histories, nil
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

	multipleRouters := args.Part == 2
	histories, err := getInput(args.InputPath, multipleRouters)
	if err != nil {
		return err
	}

	extrapolatedValueSum := 0
	for hi, history := range histories {
		diffs := [][]int{history}
		for i := 1; ; i++ {
			prevDiff := diffs[i-1]
			diff := make([]int, len(prevDiff)-1)
			isZero := true
			for j := 0; j < len(diff); j++ {
				diff[j] = prevDiff[j+1] - prevDiff[j]
				if diff[j] != 0 {
					isZero = false
				}
			}
			diffs = append(diffs, diff)
			if isZero {
				break
			}
			if len(diffs) == 1 {
				panic("empty diff")
			}
		}

		for i := len(diffs) - 2; i >= 0; i-- {
			currentDiff := diffs[i]
			prevDiff := diffs[i+1]
			currentDiff = append(currentDiff, currentDiff[len(currentDiff)-1]+prevDiff[len(prevDiff)-1])
			diffs[i] = currentDiff
		}

		for i := 0; i < len(diffs); i++ {
			fmt.Printf("H%v #%v: %v\n", hi, i, diffs[i])
		}
		extrapolatedValueSum += diffs[0][len(diffs[0])-1]
	}
	fmt.Printf("Extrapolated values sum: %v\n", extrapolatedValueSum)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
