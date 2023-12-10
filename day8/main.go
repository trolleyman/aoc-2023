package main

import (
	"bufio"
	"errors"
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
	Nodes  map[string]*Node
	Starts []*Node
	Ends   []*Node
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

func getInput(path string, multipleStarts bool) (Map, error) {
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

	var startNodes []*Node
	var endNodes []*Node
	if multipleStarts {
		for nodeName, node := range nodes {
			if strings.HasSuffix(nodeName, "A") {
				startNodes = append(startNodes, node)
			}
			if strings.HasSuffix(nodeName, "Z") {
				endNodes = append(endNodes, node)
			}
		}
		if len(startNodes) == 0 {
			return Map{}, errors.New("no start nodes")
		}
		if len(endNodes) == 0 {
			return Map{}, errors.New("no end nodes")
		}
	} else {
		startNode, err := getNode(nodes, "AAA")
		if err != nil {
			return Map{}, err
		}
		endNode, err := getNode(nodes, "ZZZ")
		if err != nil {
			return Map{}, err
		}
		startNodes = []*Node{startNode}
		endNodes = []*Node{endNode}
	}
	network := Network{
		Nodes:  nodes,
		Starts: startNodes,
		Ends:   endNodes,
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

func (n Network) isEndNode(node *Node) bool {
	for _, endNode := range n.Ends {
		if node == endNode {
			return true
		}
	}
	return false
}

func (n Network) areAllEndNodes(nodes []*Node) bool {
	for _, node := range nodes {
		if !n.isEndNode(node) {
			return false
		}
	}
	return true
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

	multipleStarts := args.Part == 2
	inputMap, err := getInput(args.InputPath, multipleStarts)
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
	printSteps := true
	step := 0
	for nodes := inputMap.Network.Starts; !inputMap.Network.areAllEndNodes(nodes); step++ {
		turnRight := inputMap.Directions[step%len(inputMap.Directions)]
		var directionChar rune
		if turnRight {
			directionChar = 'R'
		} else {
			directionChar = 'L'
		}
		if printSteps {
			fmt.Printf("Step %v: %c: ", step, directionChar)
		}
		for i, node := range nodes {
			if printSteps {
				if i > 0 {
					fmt.Print("; ")
				}
				fmt.Printf("%v -> ", node.Name)
			}
			if turnRight {
				nodes[i] = node.Right
			} else {
				nodes[i] = node.Left
			}
			if printSteps {
				fmt.Printf("%v", nodes[i].Name)
			}
		}
		if printSteps {
			fmt.Println("")
		}
	}
	fmt.Printf("Total steps: %v\n", step)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
