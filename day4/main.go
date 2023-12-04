package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
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

type Card struct {
	Id      int
	Winning []int
	Numbers []int
}

const cardPrefix = "Card "

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

func parseCard(line string) (Card, error) {
	if !strings.HasPrefix(line, cardPrefix) {
		return Card{}, errors.New(fmt.Sprintf("Line does not start with prefix %#v: %#v", cardPrefix, line))
	}

	splitLine := strings.Split(line, ":")
	if len(splitLine) != 2 {
		return Card{}, errors.New(fmt.Sprintf("Line does not have a singular ':' %#v", line))
	}

	cardIdString := strings.TrimSpace(splitLine[0][len(cardPrefix):])
	cardId, err := strconv.Atoi(cardIdString)
	if err != nil {
		return Card{}, errors.New(fmt.Sprintf("Invalid card ID %#v", cardIdString))
	}

	winningNonSplit := strings.Split(splitLine[1], "|")
	if len(winningNonSplit) != 2 {
		return Card{}, errors.New(fmt.Sprintf("Line does not have a singular '|' %#v", line))
	}

	winning, err := parseNumbers(winningNonSplit[0])
	if err != nil {
		return Card{}, err
	}

	numbers, err := parseNumbers(winningNonSplit[1])
	if err != nil {
		return Card{}, err
	}

	return Card{cardId, winning, numbers}, nil
}

func getInput(path string) ([]Card, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cards []Card
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		card, err := parseCard(line)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return cards, nil
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
		return Args{}, errors.New(fmt.Sprintf("Invalid arguments. Expected %v <part> <inputPath>", os.Args[0]))
	}
	var part int
	switch os.Args[1] {
	case "1":
		part = 1
	case "2":
		part = 2
	default:
		return Args{}, errors.New(fmt.Sprintf("Invalid part number %#v. Expected 1/2", os.Args[1]))
	}
	return Args{Part: part, InputPath: os.Args[2]}, nil
}

func (c Card) getPoints() int {
	points := 0
	for _, number := range c.Numbers {
		if slices.Contains(c.Winning, number) {
			if points == 0 {
				points = 1
			} else {
				points *= 2
			}
		}
	}
	return points
}

func run() error {
	args, err := parseArgs()
	if err != nil {
		return err
	}
	// fmt.Printf("Args: %+v\n", args)

	cards, err := getInput(args.InputPath)
	if err != nil {
		return err
	}

	switch args.Part {
	case 1:
		// Part 1
		pointSum := 0
		for _, card := range cards {
			pointSum += card.getPoints()
		}
		fmt.Printf("%v\n", pointSum)

	case 2:
		// Part 2

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