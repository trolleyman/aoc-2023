package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	t "github.com/barweiss/go-tuple"
)

type Card int

const (
	card2 Card = iota + 2
	card3
	card4
	card5
	card6
	card7
	card8
	card9
	cardT
	cardJ
	cardQ
	cardK
	cardA
)

func parseCard(r rune) (Card, error) {
	switch r {
	case '2':
		return card2, nil
	case '3':
		return card3, nil
	case '4':
		return card4, nil
	case '5':
		return card5, nil
	case '6':
		return card6, nil
	case '7':
		return card7, nil
	case '8':
		return card8, nil
	case '9':
		return card9, nil
	case 'T':
		return cardT, nil
	case 'J':
		return cardJ, nil
	case 'Q':
		return cardQ, nil
	case 'K':
		return cardK, nil
	case 'A':
		return cardA, nil
	}
	return 0, fmt.Errorf("invalid card %+v", r)
}

func (c Card) String() string {
	switch c {
	case card2:
		return "2"
	case card3:
		return "3"
	case card4:
		return "4"
	case card5:
		return "5"
	case card6:
		return "6"
	case card7:
		return "7"
	case card8:
		return "8"
	case card9:
		return "9"
	case cardT:
		return "T"
	case cardJ:
		return "J"
	case cardQ:
		return "Q"
	case cardK:
		return "K"
	case cardA:
		return "A"
	}
	return fmt.Sprintf("?%v?", int(c))
}

type HandType int

const (
	handTypeHighCard HandType = iota + 1
	handTypeOnePair
	handTypeTwoPair
	handTypeThreeOfAKind
	handTypeFullHouse
	handTypeFourOfAKind
	handTypeFiveOfAKind
)

type Hand struct {
	Cards    [5]Card
	HandType HandType
}

func All[T any](vs []T, f func(T) bool) bool {
	for _, v := range vs {
		if !f(v) {
			return false
		}
	}
	return true
}

func getHandType(cards [5]Card) HandType {
	countsMap := make(map[Card]int)
	for _, card := range cards {
		countsMap[card] += 1
	}
	var counts []t.T2[Card, int]
	for k, v := range countsMap {
		counts = append(counts, t.New2(k, v))
	}
	slices.SortFunc(counts, func(a, b t.T2[Card, int]) int { return b.V2 - a.V2 })
	switch len(counts) {
	case 1:
		return handTypeFiveOfAKind
	case 2:
		switch counts[0].V2 {
		case 4:
			return handTypeFourOfAKind
		case 3:
			return handTypeFullHouse
		}
		panic("invalid state")
	case 3:
		switch counts[0].V2 {
		case 3:
			return handTypeThreeOfAKind
		case 2:
			if counts[1].V2 != 2 {
				panic("invalid state")
			}
			return handTypeTwoPair
		}
		panic("invalid state")
	case 4:
		if counts[0].V2 != 2 {
			panic("invalid state")
		}
		return handTypeOnePair
	case 5:
		return handTypeHighCard
	}
	panic("invalid state")
}

func createHand(cards [5]Card) Hand {
	return Hand{Cards: cards, HandType: getHandType(cards)}
}

func (h Hand) String() string {
	return fmt.Sprintf("%v%v%v%v%v", h.Cards[0], h.Cards[1], h.Cards[2], h.Cards[3], h.Cards[4])
}

func parseHand(handString string) (Hand, error) {
	handString = strings.TrimSpace(handString)
	var hand []Card
	for _, r := range handString {
		card, err := parseCard(r)
		if err != nil {
			return Hand{}, err
		}
		hand = append(hand, card)
		if len(hand) > 5 {
			return Hand{}, errors.New("hand too large")
		}
	}
	if len(hand) != 5 {
		return Hand{}, errors.New("hand not 5 cards")
	}
	return createHand([5]Card{hand[0], hand[1], hand[2], hand[3], hand[4]}), nil
}

type HandBid struct {
	Hand Hand
	Bid  int
}

func getInput(path string) ([]HandBid, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var handBids []HandBid
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		lineSplit := strings.Fields(line)
		if len(lineSplit) != 2 {
			return nil, fmt.Errorf("invalid line (too many fields): %+v", line)
		}
		hand, err := parseHand(lineSplit[0])
		if err != nil {
			return nil, err
		}
		bid, err := strconv.Atoi(lineSplit[1])
		if err != nil {
			return nil, err
		}
		handBids = append(handBids, HandBid{hand, bid})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return handBids, nil
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

func sortHand(hand1 Hand, hand2 Hand) int {
	handTypeDiff := int(hand1.HandType - hand2.HandType)
	if handTypeDiff != 0 {
		return handTypeDiff
	}
	for i := range hand1.Cards {
		card1 := hand1.Cards[i]
		card2 := hand2.Cards[i]
		cardDiff := int(card1 - card2)
		if cardDiff != 0 {
			return cardDiff
		}
	}
	return 0
}

func sortHandBid(handBid1 HandBid, handBid2 HandBid) int {
	return sortHand(handBid1.Hand, handBid2.Hand)
}

func run() error {
	args, err := parseArgs()
	if err != nil {
		return err
	}
	fmt.Printf("Args: %+v\n", args)

	handBids, err := getInput(args.InputPath)
	if err != nil {
		return err
	}

	switch args.Part {
	case 1:
		// Part 1
		slices.SortStableFunc(handBids, sortHandBid)
		totalWinnings := 0
		for i, handBid := range handBids {
			rank := i + 1
			winnings := handBid.Bid * rank
			totalWinnings += winnings
			fmt.Printf("Rank %v: %+v wins %v\n", rank, handBid, winnings)
		}
		fmt.Printf("Total winnings: %v\n", totalWinnings)

	case 2:
		// Part 2
		// TODO
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
