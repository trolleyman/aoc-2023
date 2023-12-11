package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Vec2 struct {
	X int
	Y int
}

func add(a, b Vec2) Vec2 {
	return Vec2{a.X + b.X, a.Y + b.Y}
}

type Tile struct {
	Rune       rune
	Position   Vec2
	IsPipe     bool
	Connecting [4]bool
}

func (t Tile) isConnecting(d Direction) bool {
	return t.Connecting[d]
}

func (t Tile) getTile(pm *PipeMaze, d Direction) *Tile {
	return pm.getTile(add(t.Position, d.getOffset()))
}

func (t *Tile) String() string {
	return fmt.Sprintf("%c", t.Rune)
}

type Direction byte

const (
	DirectionNorth Direction = iota
	DirectionEast
	DirectionSouth
	DirectionWest
)

var Directions [4]Direction = [4]Direction{DirectionNorth, DirectionEast, DirectionSouth, DirectionWest}

func (d Direction) getOffset() Vec2 {
	switch d {
	case DirectionNorth:
		return Vec2{Y: -1}
	case DirectionEast:
		return Vec2{X: 1}
	case DirectionSouth:
		return Vec2{Y: 1}
	case DirectionWest:
		return Vec2{X: -1}
	}
	panic("invalid direction")
}

func (d Direction) invert() Direction {
	switch d {
	case DirectionNorth:
		return DirectionSouth
	case DirectionEast:
		return DirectionWest
	case DirectionSouth:
		return DirectionNorth
	case DirectionWest:
		return DirectionEast
	}
	panic("invalid direction")
}

type PipeMaze struct {
	Start Vec2
	Size  Vec2
	Tiles [][]*Tile
}

func (pm *PipeMaze) getTile(position Vec2) *Tile {
	if position.X < 0 || position.Y < 0 || position.X >= pm.Size.X || position.Y >= pm.Size.Y {
		return nil
	}
	return pm.Tiles[position.Y][position.X]
}

func isConnectingIgnoring(tile *Tile, direction Direction, ignoring func(*Tile) bool) bool {
	if tile == nil {
		return false
	}
	if ignoring(tile) {
		return false
	}
	return tile.isConnecting(direction)
}

func (pm *PipeMaze) isIntersectionFree(intersection Vec2, direction Direction, ignoring func(*Tile) bool) bool {
	newIntersection := add(intersection, direction.getOffset())
	if newIntersection.X < 0 || newIntersection.Y < 0 || newIntersection.X > pm.Size.X || newIntersection.Y > pm.Size.Y {
		return false
	}
	nw := pm.getTile(Vec2{intersection.X - 1, intersection.Y - 1})
	ne := pm.getTile(Vec2{intersection.X, intersection.Y - 1})
	sw := pm.getTile(Vec2{intersection.X - 1, intersection.Y})
	se := pm.getTile(Vec2{intersection.X, intersection.Y})
	switch direction {
	case DirectionNorth:
		return !isConnectingIgnoring(nw, DirectionEast, ignoring) && !isConnectingIgnoring(ne, DirectionWest, ignoring)
	case DirectionEast:
		return !isConnectingIgnoring(ne, DirectionSouth, ignoring) && !isConnectingIgnoring(se, DirectionNorth, ignoring)
	case DirectionSouth:
		return !isConnectingIgnoring(sw, DirectionEast, ignoring) && !isConnectingIgnoring(se, DirectionWest, ignoring)
	case DirectionWest:
		return !isConnectingIgnoring(nw, DirectionSouth, ignoring) && !isConnectingIgnoring(sw, DirectionNorth, ignoring)
	}
	panic("invalid direction")
}

func (pm *PipeMaze) String() string {
	result := ""
	for i, tileRow := range pm.Tiles {
		if i != 0 {
			result += "\n"
		}
		for _, t := range tileRow {
			result += t.String()
		}
	}
	return result
}

func getInput(path string, multipleRouters bool) (PipeMaze, error) {
	file, err := os.Open(path)
	if err != nil {
		return PipeMaze{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var start *Tile
	var tiles [][]*Tile
	for y := 0; scanner.Scan(); y++ {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var tileRow []*Tile
		for x, r := range line {
			tile := &Tile{Rune: r, Position: Vec2{x, y}, IsPipe: true}
			switch r {
			case '|': // is a vertical pipe connecting north and south.
				tile.Connecting[DirectionNorth] = true
				tile.Connecting[DirectionSouth] = true
			case '-': // is a horizontal pipe connecting east and west.
				tile.Connecting[DirectionEast] = true
				tile.Connecting[DirectionWest] = true
			case 'L': // is a 90-degree bend connecting north and east.
				tile.Connecting[DirectionNorth] = true
				tile.Connecting[DirectionEast] = true
			case 'J': // is a 90-degree bend connecting north and west.
				tile.Connecting[DirectionNorth] = true
				tile.Connecting[DirectionWest] = true
			case '7': // is a 90-degree bend connecting south and west.
				tile.Connecting[DirectionSouth] = true
				tile.Connecting[DirectionWest] = true
			case 'F': // is a 90-degree bend connecting south and east.
				tile.Connecting[DirectionEast] = true
				tile.Connecting[DirectionSouth] = true
			case '.': // is ground; there is no pipe in this tile.
				tile.IsPipe = false
			case 'S': // is the starting position of the animal; there is a pipe on this tile, but your sketch doesn't show what shape the pipe has.
				start = tile
			default:
				return PipeMaze{}, fmt.Errorf("invalid rune %c", r)
			}
			tileRow = append(tileRow, tile)
		}
		tiles = append(tiles, tileRow)
	}
	if err := scanner.Err(); err != nil {
		return PipeMaze{}, err
	}

	pipeMaze := PipeMaze{Start: start.Position, Size: Vec2{len(tiles[0]), len(tiles)}, Tiles: tiles}

	// Fixup start connecting
	for _, direction := range Directions {
		t := start.getTile(&pipeMaze, direction)
		if t != nil {
			// fmt.Printf("t @ %v (conn=%v) (d=%v di=%v)\n", t.Position, t.Connecting, direction, direction.invert())
			if t.isConnecting(direction.invert()) {
				start.Connecting[direction] = true
			}
		}
	}
	return pipeMaze, nil
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
	pipeMaze, err := getInput(args.InputPath, multipleRouters)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n\n", pipeMaze.String())
	var intMaze [][]int
	for y := 0; y < pipeMaze.Size.Y; y++ {
		var intRow []int
		for x := 0; x < pipeMaze.Size.X; x++ {
			intRow = append(intRow, -1)
		}
		intMaze = append(intMaze, intRow)
	}

	step := 0
	for expandSet := []*Tile{pipeMaze.getTile(pipeMaze.Start)}; len(expandSet) > 0; step++ {
		var newExpandSet []*Tile
		for _, t := range expandSet {
			// fmt.Printf("Tile %v @ %v:", t, t.Position)
			if intMaze[t.Position.Y][t.Position.X] >= 0 {
				// fmt.Printf(" <excluded>\n")
				continue
			}
			intMaze[t.Position.Y][t.Position.X] = step
			// fmt.Printf(" <set to %v>", step)
			for _, d := range Directions {
				// fmt.Printf(" <dir %v>", d)
				if t.isConnecting(d) {
					// fmt.Printf(" <connecting %v>", d)
					newTile := pipeMaze.getTile(add(t.Position, d.getOffset()))
					if newTile != nil && intMaze[newTile.Position.Y][newTile.Position.X] < 0 {
						// fmt.Printf(" <added %v>", newTile.Position)
						newExpandSet = append(newExpandSet, newTile)
					}
				}
			}
			// fmt.Printf("\n")
		}
		expandSet = newExpandSet
	}
	step -= 1

	for y := 0; y < pipeMaze.Size.Y; y++ {
		for x := 0; x < pipeMaze.Size.X; x++ {
			val := intMaze[y][x]
			if val < 0 {
				fmt.Printf("%v", pipeMaze.getTile(Vec2{x, y}))
			} else if val == 0 {
				fmt.Print("S")
			} else if val == step {
				fmt.Print("X")
			} else {
				fmt.Printf("%v", val%10)
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")

	fmt.Printf("Max step: %v\n", step)

	if args.Part == 1 {
		return nil
	}

	// Now get all the tiles enclosed by the loop
	var intersectionExpandSet []Vec2
	var isIntersectionOutside [][]bool
	for y := 0; y <= pipeMaze.Size.Y; y++ {
		var intersectionRow []bool
		for x := 0; x <= pipeMaze.Size.X; x++ {
			isOutside := x == 0 || y == 0 || x == pipeMaze.Size.X || y == pipeMaze.Size.Y
			intersectionRow = append(intersectionRow, isOutside)
			if isOutside {
				intersectionExpandSet = append(intersectionExpandSet, Vec2{x, y})
			}
		}
		isIntersectionOutside = append(isIntersectionOutside, intersectionRow)
	}

	for len(intersectionExpandSet) > 0 {
		var newIntersectionExpandSet []Vec2
		for _, intersection := range intersectionExpandSet {
			// fmt.Printf("Intersection %v:", intersection)
			for _, direction := range Directions {
				// fmt.Printf(" <dir %v>", direction)
				if pipeMaze.isIntersectionFree(intersection, direction, func(t *Tile) bool { return intMaze[t.Position.Y][t.Position.X] < 0 }) {
					newIntersection := add(intersection, direction.getOffset())
					// fmt.Printf(" <intersection free %v>", newIntersection)
					if !isIntersectionOutside[newIntersection.Y][newIntersection.X] {
						isIntersectionOutside[newIntersection.Y][newIntersection.X] = true
						// fmt.Print(" <outside>")
						newIntersectionExpandSet = append(newIntersectionExpandSet, newIntersection)
					}
				}
			}
			// fmt.Printf("\n")
		}
		intersectionExpandSet = newIntersectionExpandSet
	}

	fmt.Printf("\n")
	for _, intersectionRow := range isIntersectionOutside {
		for _, isOutside := range intersectionRow {
			if isOutside {
				fmt.Print("-")
			} else {
				fmt.Print("@")
			}
		}
		fmt.Printf("\n")
	}

	fmt.Printf("\n")
	insideCount := 0
	for y, tileRow := range pipeMaze.Tiles {
		for x, tile := range tileRow {
			if tile.IsPipe && intMaze[y][x] >= 0 {
				fmt.Printf("%v", tile)
			} else {
				isOutside := isIntersectionOutside[y][x] && isIntersectionOutside[y+1][x] && isIntersectionOutside[y][x+1] && isIntersectionOutside[y+1][x+1]
				if isOutside {
					if true {
						fmt.Printf("%v", tile)
					} else {
						fmt.Printf("O")
					}
				} else {
					fmt.Printf("I")
					insideCount++
				}
			}
		}
		fmt.Printf("\n")
	}

	fmt.Printf("\nTotal inside: %v\n", insideCount)
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
