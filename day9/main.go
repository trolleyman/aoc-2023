package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	t "github.com/barweiss/go-tuple"
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
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return histories, nil
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

	multipleRouters := args.Part == 2
	inputMap, err := getInput(args.InputPath, multipleRouters)
	if err != nil {
		return err
	}

	fmt.Printf("Directions (%v): ", len(inputMap.Directions))
	for _, direction := range inputMap.Directions {
		fmt.Printf("%v", direction)
	}
	fmt.Println("")
	totalSteps := getTotalStepsAnalytical(inputMap)
	fmt.Printf("Total steps: %v\n", totalSteps)

	return nil
}

func gcd(a, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

func lcm(a, b int) int {
	return a / gcd(a, b) * b
}

func getTotalStepsAnalytical(inputMap Map) int {
	totalSteps := 1
	for i, node := range inputMap.Network.Starts {
		offset, length, endSteps := getPeriodicity(inputMap, node)
		fmt.Printf("Ghost %v: start=%v, offset=%v, length=%v /lendir=%v endSteps=%v\n", i, node.Name, offset, length, length/len(inputMap.Directions), endSteps)
		totalSteps = lcm(totalSteps, length)
	}
	return totalSteps
}

func getTotalSteps(inputMap Map) int {
	printSteps := true
	step := 0
	for nodes := inputMap.Network.Starts; !inputMap.Network.areAllEndNodes(nodes); step++ {
		direction := inputMap.Directions[step%len(inputMap.Directions)]
		if printSteps {
			fmt.Printf("Step %v (%v): ", step, direction)
		}
		for i, node := range nodes {
			if printSteps {
				if i > 0 {
					fmt.Print("; ")
				}
				fmt.Printf("%v -> ", node.Name)
			}
			nodes[i] = node.getNextNode(direction)
			if printSteps {
				fmt.Printf("%v", nodes[i].Name)
			}
		}
		if printSteps {
			fmt.Println("")
		}
	}
	return step
}

func getPeriodicity(inputMap Map, startNode *Node) (offset int, length int, endSteps []int) {
	nodeDirectionIndexSteps := make(map[t.T2[*Node, int]]int)
	for step, node := 0, startNode; ; step++ {
		directionIndex := step % len(inputMap.Directions)
		nodeDirectionIndex := t.New2(node, directionIndex)
		if nodeDirectionIndexStep, found := nodeDirectionIndexSteps[nodeDirectionIndex]; found {
			// Periodicity found
			offset = nodeDirectionIndexStep
			length = step - nodeDirectionIndexStep
			return
		}
		if inputMap.Network.isEndNode(node) {
			endSteps = append(endSteps, step)
		}
		nodeDirectionIndexSteps[nodeDirectionIndex] = step
		node = node.getNextNode(inputMap.Directions[directionIndex])
	}
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
