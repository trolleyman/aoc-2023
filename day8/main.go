package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	t "github.com/barweiss/go-tuple"
)

type Node struct {
	Name  string
	Left  *Node
	Right *Node
}

type Network struct {
	Nodes map[string]*Node
	Start *Node
	End   *Node
}

type Map struct {
	// false is left, true is right
	Directions []bool
	Network    Network
}

func parseDirections(directionsString string) ([]bool, error) {
	result := make([]bool, 0, len(directionsString))
	for _, r := range directionsString {
		switch r {
		case 'L':
			result = append(result, false)
		case 'R':
			result = append(result, true)
		default:
			return nil, fmt.Errorf("invalid directions string (encountered %+v): %+v", r, directionsString)
		}
	}
	return result, nil
}

func getInput(path string) (Map, error) {
	file, err := os.Open(path)
	if err != nil {
		return Map{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	nodes := make(map[string]*Node)
	nodeDirections := make(map[string]t.T2[string, string])

	if !scanner.Scan() {
		return Map{}, fmt.Errorf("no directions input")
	}
	directionsString := strings.TrimSpace(scanner.Text())
	directions, err := parseDirections(directionsString)
	if err != nil {
		return Map{}, err
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		lineSplit := strings.Split(line, "=")
		if len(lineSplit) != 2 {
			return Map{}, fmt.Errorf("invalid line (expected 1 '='): %+v", line)
		}

		nodeName := strings.TrimSpace(lineSplit[0])
		nodes[nodeName] = &Node{Name: nodeName}
		leftRightString := strings.TrimSpace(lineSplit[1])
		leftRightString, foundPrefix := strings.CutPrefix(leftRightString, "(")
		if !foundPrefix {
			return Map{}, fmt.Errorf("invalid line (expected prefix '('): %+v", line)
		}
		leftRightString, foundSuffix := strings.CutSuffix(leftRightString, ")")
		if !foundSuffix {
			return Map{}, fmt.Errorf("invalid line (expected suffix ')'): %+v", line)
		}

		leftRightStringSplit := strings.Split(leftRightString, ",")
		if len(leftRightStringSplit) != 2 {
			return Map{}, fmt.Errorf("invalid line (expected 1 ','): %+v", line)
		}

		leftString := strings.TrimSpace(leftRightStringSplit[0])
		rightString := strings.TrimSpace(leftRightStringSplit[1])
		nodeDirections[nodeName] = t.New2(leftString, rightString)
	}
	if err := scanner.Err(); err != nil {
		return Map{}, err
	}

	for nodeName, nodeDirectionsTuple := range nodeDirections {
		node := nodes[nodeName]
		node.Left, err = getNode(nodes, nodeDirectionsTuple.V1)
		if err != nil {
			return Map{}, err
		}
		node.Right, err = getNode(nodes, nodeDirectionsTuple.V2)
		if err != nil {
			return Map{}, err
		}
	}

	nodeStart, err := getNode(nodes, "AAA")
	if err != nil {
		return Map{}, err
	}
	nodeEnd, err := getNode(nodes, "ZZZ")
	if err != nil {
		return Map{}, err
	}
	network := Network{
		Nodes: nodes,
		Start: nodeStart,
		End:   nodeEnd,
	}

	return Map{Network: network, Directions: directions}, nil
}

func getNode(nodes map[string]*Node, nodeName string) (*Node, error) {
	node := nodes[nodeName]
	if node == nil {
		return nil, fmt.Errorf("unknown node %+v", nodeName)
	}
	return node, nil
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

	inputMap, err := getInput(args.InputPath)
	if err != nil {
		return err
	}

	fmt.Printf("Directions: ")
	for _, isRight := range inputMap.Directions {
		if isRight {
			fmt.Print("R")
		} else {
			fmt.Print("L")
		}
	}
	fmt.Println("")
	i := 0
	for node := inputMap.Network.Start; node != inputMap.Network.End; i++ {
		turnRight := inputMap.Directions[i%len(inputMap.Directions)]
		var newNode *Node
		var directionChar rune
		if turnRight {
			directionChar = 'R'
			newNode = node.Right
		} else {
			directionChar = 'L'
			newNode = node.Left
		}
		fmt.Printf("Step %v: %v (%c) -> %v\n", i, node.Name, directionChar, newNode.Name)
		node = newNode
	}
	fmt.Printf("Total steps: %v\n", i)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
