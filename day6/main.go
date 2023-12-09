package main

import (
	"bufio"
	"errors"
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

type Race struct {
	Time     int
	Distance int
}

func parseNumbersPrefix(s string, prefix string) ([]int, error) {
	s, ok := strings.CutPrefix(s, prefix)
	if !ok {
		return nil, fmt.Errorf("Invalid string doesn't start with %#v: %#v", prefix, s)
	}
	return parseNumbers(s)
}

func scanNumbersPrefix(scanner *bufio.Scanner, prefix string) ([]int, error) {
	if !scanner.Scan() {
		return nil, errors.New("Line expected, no line found")
	}
	return parseNumbersPrefix(scanner.Text(), prefix)
}

func getInput(path string) ([]Race, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	times, err := scanNumbersPrefix(scanner, "Time:")
	if err != nil {
		return nil, err
	}
	distances, err := scanNumbersPrefix(scanner, "Distance:")
	if err != nil {
		return nil, err
	}

	if len(times) != len(distances) {
		return nil, fmt.Errorf("Times (%#v) and distances (%#v) lengths don't match", times, distances)
	}

	output := make([]Race, len(times))
	for i := 0; i < len(times); i++ {
		output[i] = Race{Time: times[i], Distance: distances[i]}
	}
	return output, nil
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

	races, err := getInput(args.InputPath)
	if err != nil {
		return err
	}

	switch args.Part {
	case 1:
		// Part 1
		marginOfError := 1
		fmt.Printf("Races: %+v\n", races)
		for raceI, race := range races {
			numWaysToWin := 0
			for buttonTime := 0; buttonTime <= race.Time; buttonTime++ {
				speed := buttonTime
				speedTime := race.Time - buttonTime
				distance := speed * speedTime
				fmt.Printf("Race %v: button %vms speed %vms @ %vmm/s distance %v", raceI, buttonTime, speedTime, speed, distance)
				if distance > race.Distance {
					fmt.Printf(" (win)")
					numWaysToWin++
				}
				fmt.Printf("\n")
			}
			fmt.Printf("Race %v: %v ways to win\n", raceI, numWaysToWin)
			marginOfError *= numWaysToWin
		}
		fmt.Printf("Margin of error: %v\n", marginOfError)

	case 2:
		// Part 2
		panic("NYI")
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
