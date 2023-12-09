package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func parseNumbers(numbersString string, ignoreSpaces bool) ([]int, error) {
	var numbers []int
	numbersString = strings.TrimSpace(numbersString)
	if ignoreSpaces {
		newNumberString := ""
		for _, numberString := range strings.Fields(numbersString) {
			newNumberString += numberString
		}
		number, err := strconv.Atoi(newNumberString)
		if err != nil {
			return nil, err
		}
		return []int{number}, nil
	} else {
		for _, numberString := range strings.Fields(numbersString) {
			numberString = strings.TrimSpace(numberString)
			number, err := strconv.Atoi(numberString)
			if err != nil {
				return nil, err
			}
			numbers = append(numbers, number)
		}
	}
	return numbers, nil
}

type Race struct {
	Time     int
	Distance int
}

func parseNumbersPrefix(s string, prefix string, ignoreSpaces bool) ([]int, error) {
	s, ok := strings.CutPrefix(s, prefix)
	if !ok {
		return nil, fmt.Errorf("invalid string doesn't start with %#v: %#v", prefix, s)
	}
	return parseNumbers(s, ignoreSpaces)
}

func scanNumbersPrefix(scanner *bufio.Scanner, prefix string, ignoreSpaces bool) ([]int, error) {
	if !scanner.Scan() {
		return nil, errors.New("line expected, no line found")
	}
	return parseNumbersPrefix(scanner.Text(), prefix, ignoreSpaces)
}

func getInput(path string, ignoreSpaces bool) ([]Race, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	times, err := scanNumbersPrefix(scanner, "Time:", ignoreSpaces)
	if err != nil {
		return nil, err
	}
	distances, err := scanNumbersPrefix(scanner, "Distance:", ignoreSpaces)
	if err != nil {
		return nil, err
	}

	if len(times) != len(distances) {
		return nil, fmt.Errorf("times (%#v) and distances (%#v) lengths don't match", times, distances)
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

func solveQuadratic(a float64, b float64, c float64) (float64, float64, bool) {
	sqrtIn := b*b - 4*a*c
	if sqrtIn < 0 {
		return 0, 0, false
	}
	sqrt := math.Sqrt(sqrtIn)
	a2 := a * 2
	// fmt.Printf("a:%v b:%v c:%v sqrt:%v a2:%v\n", a, b, c, sqrt, a2)
	minSolution := (-b - sqrt) / a2
	maxSolution := (-b + sqrt) / a2
	if minSolution > maxSolution {
		minSolution, maxSolution = maxSolution, minSolution
	}
	return minSolution, maxSolution, true
}

func getWinButtonTimes(race Race) (min int, max int, found bool) {
	minSolution, maxSolution, found := solveQuadratic(-1, float64(race.Time), -float64(race.Distance))
	// fmt.Printf("minSolution: %v, maxSolution: %v\n", minSolution, maxSolution)
	if !found {
		return 0, 0, false
	}
	return int(math.Ceil(minSolution)), int(math.Floor(maxSolution)), true
}

func run() error {
	args, err := parseArgs()
	if err != nil {
		return err
	}
	fmt.Printf("Args: %+v\n", args)

	races, err := getInput(args.InputPath, args.Part == 2)
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
				if canWinRace(raceI, race, buttonTime) {
					numWaysToWin++
				}
			}
			fmt.Printf("Race %v: %v ways to win\n", raceI, numWaysToWin)
			marginOfError *= numWaysToWin
		}
		fmt.Printf("Margin of error: %v\n", marginOfError)

	case 2:
		// Part 2
		marginOfError := 1
		fmt.Printf("Races: %+v\n", races)
		for raceI, race := range races {
			numWaysToWin := 0
			min, max, found := getWinButtonTimes(race)
			if found {
				numWaysToWin = max - min + 1
			}
			fmt.Printf("Race %v: %v ways to win (min=%v, max=%v)\n", raceI, numWaysToWin, min, max)
			marginOfError *= numWaysToWin
		}
		fmt.Printf("Margin of error: %v\n", marginOfError)
	}

	return nil
}

func canWinRace(raceI int, race Race, buttonTime int) bool {
	speed := buttonTime
	speedTime := race.Time - buttonTime
	distance := speed * speedTime
	fmt.Printf("Race %v: button %vms speed %vms @ %vmm/s distance %v", raceI, buttonTime, speedTime, speed, distance)
	canWin := distance > race.Distance
	if canWin {
		fmt.Printf(" (win)")
	}
	fmt.Printf("\n")
	return canWin
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
