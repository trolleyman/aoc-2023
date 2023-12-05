package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

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
		return Card{}, fmt.Errorf("line does not start with prefix %#v: %#v", cardPrefix, line)
	}

	splitLine := strings.Split(line, ":")
	if len(splitLine) != 2 {
		return Card{}, fmt.Errorf("line does not have a singular ':' %#v", line)
	}

	cardIdString := strings.TrimSpace(splitLine[0][len(cardPrefix):])
	cardId, err := strconv.Atoi(cardIdString)
	if err != nil {
		return Card{}, fmt.Errorf("invalid card ID %#v", cardIdString)
	}

	winningNonSplit := strings.Split(splitLine[1], "|")
	if len(winningNonSplit) != 2 {
		return Card{}, fmt.Errorf("line does not have a singular '|' %#v", line)
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

func (c Card) getMatchingWinningNumbersCount() int {
	count := 0
	for _, number := range c.Numbers {
		if slices.Contains(c.Winning, number) {
			count += 1
		}
	}
	return count
}

func intPow(n, m int) int {
	if m == 0 {
		return 1
	}
	result := n
	for i := 2; i <= m; i++ {
		result *= n
	}
	return result
}

func (c Card) getPoints() int {
	count := c.getMatchingWinningNumbersCount()
	if count == 0 {
		return 0
	} else {
		return intPow(2, count-1)
	}
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
		originalCards := cards[:]
		for i := 0; i < len(cards); i++ {
			card := cards[i]
			// fmt.Printf("%+v\n", card)
			winningCount := card.getMatchingWinningNumbersCount()
			for j := card.Id; j < len(originalCards) && j < card.Id+winningCount; j++ {
				extraCard := originalCards[j]
				// fmt.Printf("+ %+v\n", extraCard)
				cards = append(cards, extraCard)
			}
		}
		fmt.Printf("%v\n", len(cards))
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
