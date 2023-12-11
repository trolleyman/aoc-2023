package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Space bool

func (s Space) String() string {
	if s {
		return "#"
	} else {
		return "."
	}
}

type Universe [][]Space

func (u Universe) String() string {
	result := ""
	for i, row := range u {
		if i != 0 {
			result += "\n"
		}
		for _, s := range row {
			result += s.String()
		}
	}
	return result
}

func (u Universe) expand() Universe {
	emptyRows, emptyColumns := u.getEmptyRowsColumns()

	newU := make([][]Space, 0, len(u)+len(emptyRows))
	for y, row := range u {
		newRowLen := len(row) + len(emptyColumns)
		if emptyRows[y] {
			newU = append(newU, make([]Space, newRowLen), make([]Space, newRowLen))
		} else {
			newRow := make([]Space, 0, newRowLen)
			for x, s := range row {
				if emptyColumns[x] {
					newRow = append(newRow, false, false)
				} else {
					newRow = append(newRow, s)
				}
			}
			newU = append(newU, newRow)
		}
	}
	return newU
}

func (u Universe) getEmptyRowsColumns() (emptyRows map[int]bool, emptyColumns map[int]bool) {
	emptyRows = make(map[int]bool)
	for y, row := range u {
		isRowEmpty := true
		for _, s := range row {
			if s {
				isRowEmpty = false
				break
			}
		}
		if isRowEmpty {
			emptyRows[y] = true
		}
	}

	emptyColumns = make(map[int]bool)
	for x := 0; x < len(u[0]); x++ {
		isColumnEmpty := true
		for _, row := range u {
			if row[x] {
				isColumnEmpty = false
				break
			}
		}
		if isColumnEmpty {
			emptyColumns[x] = true
		}
	}
	return
}

func getInput(path string) (Universe, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var universe Universe
	for y := 0; scanner.Scan(); y++ {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var row []Space
		for _, r := range line {
			var space Space
			switch r {
			case '#':
				space = true
			case '.':
				space = false
			default:
				return nil, fmt.Errorf("invalid rune %c", r)
			}
			row = append(row, space)
		}
		universe = append(universe, row)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return universe, nil
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

	universe, err := getInput(args.InputPath)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", universe.String())

	expandedUniverse := universe.expand()
	fmt.Printf("\n%v\n\n", expandedUniverse.String())

	shortestPathSum := expandedUniverse.getShortestPathSum()
	fmt.Printf("\nShortest path sum: %v\n", shortestPathSum)

	if args.Part == 2 {
		expandedTimes := 1000000
		shortestPathSum = universe.getShortestPathSumExpanded(expandedTimes)

		fmt.Printf("\nShortest path sum (expanded=%v): %v\n", expandedTimes, shortestPathSum)
	}

	return nil
}

type Vec2 struct {
	X int
	Y int
}

func (u Universe) getGalaxies() (galaxies []Vec2) {
	for y, row := range u {
		for x, s := range row {
			if s {
				galaxies = append(galaxies, Vec2{x, y})
			}
		}
	}
	return
}

func getShortestPathSumExpanded(galaxies []Vec2, i int, emptyRows map[int]bool, emptyColumns map[int]bool, expandedTimes int) (sum int) {
	galaxy := galaxies[i]
	for j := i + 1; j < len(galaxies); j++ {
		if i == j {
			continue
		}
		otherGalaxy := galaxies[j]
		minX, maxX := galaxy.X, otherGalaxy.X
		if minX > maxX {
			minX, maxX = maxX, minX
		}
		minY, maxY := galaxy.Y, otherGalaxy.Y
		if minY > maxY {
			minY, maxY = maxY, minY
		}
		distance := maxX - minX + maxY - minY
		for x := minX + 1; x < maxX; x++ {
			if emptyColumns[x] {
				distance += expandedTimes - 1
				// fmt.Printf("Empty column %v\n", x)
			}
		}
		for y := minY + 1; y < maxY; y++ {
			if emptyRows[y] {
				distance += expandedTimes - 1
				// fmt.Printf("Empty row %v\n", y)
			}
		}
		// fmt.Printf("Galaxy %v -> %v = %v\n", galaxy, otherGalaxy, distance)
		sum += distance
	}
	return
}

func (u Universe) getShortestPathSumExpanded(expandedTimes int) (sum int) {
	emptyRows, emptyColumns := u.getEmptyRowsColumns()
	galaxies := u.getGalaxies()
	for i := range galaxies {
		sum += getShortestPathSumExpanded(galaxies, i, emptyRows, emptyColumns, expandedTimes)
	}
	return
}

func getShortestPathSum(galaxies []Vec2, i int) (sum int) {
	galaxy := galaxies[i]
	for j := i + 1; j < len(galaxies); j++ {
		if i == j {
			continue
		}
		otherGalaxy := galaxies[j]
		xdiff := galaxy.X - otherGalaxy.X
		if xdiff < 0 {
			xdiff = -xdiff
		}
		ydiff := galaxy.Y - otherGalaxy.Y
		if ydiff < 0 {
			ydiff = -ydiff
		}
		distance := xdiff + ydiff
		// fmt.Printf("Galaxy %v -> %v = %v\n", galaxy, otherGalaxy, distance)
		sum += distance
	}
	return
}

func (u Universe) getShortestPathSum() (sum int) {
	galaxies := u.getGalaxies()
	for i := range galaxies {
		sum += getShortestPathSum(galaxies, i)
	}
	return
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
