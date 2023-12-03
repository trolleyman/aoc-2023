package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
)

type Color int

const (
	red Color = iota + 1
	blue
	green
)

var colors = [3]Color{red, green, blue}

func parseColor(s string) (Color, error) {
	switch s {
	case "red":
		return red, nil
	case "blue":
		return blue, nil
	case "green":
		return green, nil
	default:
		return Color(0), errors.New(fmt.Sprintf("Invalid color %#v", s))
	}
}

func (c Color) String() string {
	switch c {
	case red:
		return "red"
	case green:
		return "green"
	case blue:
		return "blue"
	default:
		panic(fmt.Sprintf("Invalid color %d", int(c)))
	}
}

type Game struct {
	Id   int
	Sets []Set
}

type Set map[Color]int

const gamePrefix = "Game "

func parseSet(setString string) (Set, error) {
	set := make(map[Color]int)
	setString = strings.TrimSpace(setString)
	splitSetString := strings.Split(setString, ",")
	for _, setItemString := range splitSetString {
		setItemString = strings.TrimSpace(setItemString)
		setItemStringSplit := strings.Split(setItemString, " ")
		if len(setItemStringSplit) != 2 {
			return nil, errors.New(fmt.Sprintf("Invalid set item %#v", setItemString))
		}
		colorCount, err := strconv.Atoi(setItemStringSplit[0])
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Invalid set item count %#v", setItemStringSplit[0]))
		}
		color, err := parseColor(setItemStringSplit[1])
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Invalid set item color %#v", setItemStringSplit[1]))
		}
		set[color] = colorCount
	}
	return set, nil
}

func parseGame(line string) (Game, error) {
	if !strings.HasPrefix(line, gamePrefix) {
		return Game{}, errors.New(fmt.Sprintf("Line does not start with prefix %#v: %#v", gamePrefix, line))
	}

	splitLine := strings.Split(line, ":")
	if len(splitLine) != 2 {
		return Game{}, errors.New(fmt.Sprintf("Line does not have a singular ':' %#v", line))
	}

	gameIdString := splitLine[0][len(gamePrefix):]
	gameId, err := strconv.Atoi(gameIdString)
	if err != nil {
		return Game{}, errors.New(fmt.Sprintf("Invalid game ID %#v", gameIdString))
	}

	setsStrings := strings.Split(splitLine[1], ";")
	var sets []Set = nil
	for _, setString := range setsStrings {
		set, err := parseSet(setString)
		if err != nil {
			return Game{}, err
		}
		sets = append(sets, set)
	}
	return Game{gameId, sets}, nil
}

func getInput(path string) ([]Game, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var games []Game
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		game, err := parseGame(line)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return games, nil
}

type Args struct {
	InputPath string
	InputSet  *Set
}

func parseArgs() (Args, error) {
	var inputSet *Set = nil
	switch len(os.Args) {
	case 2:
		break
	case 3:
		parsedInputSet, err := parseSet(os.Args[2])
		if err != nil {
			return Args{}, errors.New(fmt.Sprintf("Invalid set %#v: %v", os.Args[2], err))
		}
		inputSet = &parsedInputSet
	default:
		return Args{}, errors.New(fmt.Sprintf("Invalid arguments. Expected %v <inputPath> [inputSet]", os.Args[0]))
	}
	return Args{InputPath: os.Args[1], InputSet: inputSet}, nil
}

func isGamePossible(game Game, inputSet Set) bool {
	for _, set := range game.Sets {
		for color, count := range set {
			inputCount := inputSet[color]
			if count > inputCount {
				return false
			}
		}
	}
	return true
}

func max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func getMinimumInputSet(game Game) Set {
	minimumSet := make(Set)
	for _, set := range game.Sets {
		for _, color := range colors {
			minimumSet[color] = max(minimumSet[color], set[color])
		}
	}
	return minimumSet
}

func getSetPower(minimumSet Set) int {
	power := 1
	for _, count := range minimumSet {
		power *= count
	}
	return power
}

func run() error {
	args, err := parseArgs()
	if err != nil {
		return err
	}
	// fmt.Printf("Args: %+v\n", args)

	games, err := getInput(args.InputPath)
	if err != nil {
		return err
	}

	if args.InputSet != nil {
		// Part 1
		gameIdSum := 0
		for _, game := range games {
			possible := isGamePossible(game, *args.InputSet)
			// fmt.Printf("%v: %+v\n", possible, game)
			if possible {
				gameIdSum += game.Id
			}
		}
		fmt.Printf("%v\n", gameIdSum)

	} else {
		// Part 2
		powerSum := 0
		for _, game := range games {
			minimumSet := getMinimumInputSet(game)
			minimumSetPower := getSetPower(minimumSet)
			fmt.Printf("minimumSet=%+v power=%v: %+v\n", minimumSet, minimumSetPower, game)
			powerSum += minimumSetPower
		}
		fmt.Printf("%v\n", powerSum)
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
