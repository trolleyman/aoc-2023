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
	Card2 Card = iota + 2
	Card3
	Card4
	Card5
	Card6
	Card7
	Card8
	Card9
	CardT
	CardJ
	CardQ
	CardK
	CardA
)

func parseCard(r rune) (Card, error) {
	switch r {
	case '2':
		return Card2, nil
	case '3':
		return Card3, nil
	case '4':
		return Card4, nil
	case '5':
		return Card5, nil
	case '6':
		return Card6, nil
	case '7':
		return Card7, nil
	case '8':
		return Card8, nil
	case '9':
		return Card9, nil
	case 'T':
		return CardT, nil
	case 'J':
		return CardJ, nil
	case 'Q':
		return CardQ, nil
	case 'K':
		return CardK, nil
	case 'A':
		return CardA, nil
	}
	return 0, fmt.Errorf("invalid card %+v", r)
}

func (c Card) String() string {
	switch c {
	case Card2:
		return "2"
	case Card3:
		return "3"
	case Card4:
		return "4"
	case Card5:
		return "5"
	case Card6:
		return "6"
	case Card7:
		return "7"
	case Card8:
		return "8"
	case Card9:
		return "9"
	case CardT:
		return "T"
	case CardJ:
		return "J"
	case CardQ:
		return "Q"
	case CardK:
		return "K"
	case CardA:
		return "A"
	}
	return fmt.Sprintf("?%v?", int(c))
}

type HandType int

const (
	HandTypeHighCard HandType = iota + 1
	HandTypeOnePair
	HandTypeTwoPair
	HandTypeThreeOfAKind
	HandTypeFullHouse
	HandTypeFourOfAKind
	HandTypeFiveOfAKind
)

func (ht HandType) String() string {
	switch ht {
	case HandTypeHighCard:
		return "HighCard"
	case HandTypeOnePair:
		return "OnePair"
	case HandTypeTwoPair:
		return "TwoPair"
	case HandTypeThreeOfAKind:
		return "ThreeOfAKind"
	case HandTypeFullHouse:
		return "FullHouse"
	case HandTypeFourOfAKind:
		return "FourOfAKind"
	case HandTypeFiveOfAKind:
		return "FiveOfAKind"
	}
	return "Unknown"
}

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

func getCounts[T comparable](values []T) map[T]int {
	counts := make(map[T]int)
	for _, value := range values {
		counts[value] += 1
	}
	return counts
}

func getHandType(cards [5]Card, jackIsJoker bool) HandType {
	countsMap := getCounts(cards[:])
	wildCount := 0
	if jackIsJoker {
		wildCount = countsMap[CardJ]
		delete(countsMap, CardJ)
	}

	var counts []t.T2[Card, int]
	for k, v := range countsMap {
		counts = append(counts, t.New2(k, v))
	}
	slices.SortFunc(counts, func(a, b t.T2[Card, int]) int { return b.V2 - a.V2 })
	if wildCount > 0 && len(counts) > 0 {
		counts[0] = t.New2(counts[0].V1, counts[0].V2+wildCount)
	}
	switch len(countsMap) {
	case 0:
		return HandTypeFiveOfAKind
	case 1:
		return HandTypeFiveOfAKind
	case 2:
		if counts[0].V2 == 4 {
			return HandTypeFourOfAKind
		}
		return HandTypeFullHouse
	case 3:
		if counts[0].V2 == 3 {
			return HandTypeThreeOfAKind
		}
		return HandTypeTwoPair
	case 4:
		return HandTypeOnePair
	case 5:
		return HandTypeHighCard
	}
	panic(fmt.Sprintf("unexpected: cards=%v len(countsMap)=%v wildCount=%v", cards, len(countsMap), wildCount))
}

func createHand(cards [5]Card, jackIsJoker bool) Hand {
	return Hand{Cards: cards, HandType: getHandType(cards, jackIsJoker)}
}

func (h Hand) String() string {
	return fmt.Sprintf("%v%v%v%v%v (%v)", h.Cards[0], h.Cards[1], h.Cards[2], h.Cards[3], h.Cards[4], h.HandType)
}

func parseHand(handString string, jackIsJoker bool) (Hand, error) {
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
	return createHand([5]Card{hand[0], hand[1], hand[2], hand[3], hand[4]}, jackIsJoker), nil
}

type HandBid struct {
	Hand Hand
	Bid  int
}

func getInput(path string, jackIsJoker bool) ([]HandBid, error) {
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
		hand, err := parseHand(lineSplit[0], jackIsJoker)
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

func getSortCardFunc(jackIsJoker bool) func(Card, Card) int {
	return func(card1, card2 Card) int {
		cardInt1 := int(card1)
		cardInt2 := int(card2)
		if jackIsJoker {
			if card1 == CardJ {
				cardInt1 = int(Card2 - 1)
			}
			if card2 == CardJ {
				cardInt2 = int(Card2 - 1)
			}
		}
		return cardInt1 - cardInt2
	}
}

func getSortHandFunc(jackIsJoker bool) func(Hand, Hand) int {
	cardSortFunc := getSortCardFunc(jackIsJoker)
	return func(hand1, hand2 Hand) int {
		handTypeDiff := int(hand1.HandType - hand2.HandType)
		if handTypeDiff != 0 {
			return handTypeDiff
		}
		for i := range hand1.Cards {
			card1 := hand1.Cards[i]
			card2 := hand2.Cards[i]
			cardDiff := cardSortFunc(card1, card2)
			if cardDiff != 0 {
				return cardDiff
			}
		}
		return 0
	}
}

func getSortHandBidFunc(jackIsJoker bool) func(HandBid, HandBid) int {
	f := getSortHandFunc(jackIsJoker)
	return func(handBid1, handBid2 HandBid) int {
		return f(handBid1.Hand, handBid2.Hand)
	}
}

func run() error {
	args, err := parseArgs()
	if err != nil {
		return err
	}
	fmt.Printf("Args: %+v\n", args)

	jackIsJoker := args.Part == 2
	handBids, err := getInput(args.InputPath, jackIsJoker)
	if err != nil {
		return err
	}

	fmt.Printf("Total winnings: %v\n", getTotalWinnings(handBids, jackIsJoker))

	return nil
}

func getTotalWinnings(handBids []HandBid, jackIsJoker bool) int {
	slices.SortStableFunc(handBids, getSortHandBidFunc(jackIsJoker))
	totalWinnings := 0
	for i, handBid := range handBids {
		rank := i + 1
		winnings := handBid.Bid * rank
		totalWinnings += winnings
		fmt.Printf("Rank %v: %+v wins %v\n", rank, handBid, winnings)
	}
	return totalWinnings
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
