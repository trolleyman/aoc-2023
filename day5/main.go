package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type RangeMapItem struct {
	SourceStart      int
	DestinationStart int
	Length           int
}

type RangeMap []RangeMapItem

func (rangeMap RangeMap) getDestination(source int) int {
	for _, item := range rangeMap {
		if source >= item.SourceStart && source < item.SourceStart+item.Length {
			return source - item.SourceStart + item.DestinationStart
		}
	}
	return source
}

func (rangeMap RangeMap) getDestinations(sourceRanges []Range) []Range {
	destinationRanges := make([]Range, 0, len(sourceRanges))
	for _, sourceRange := range sourceRanges {
		for _, item := range rangeMap {
			if sourceRange.Length <= 0 {
				break
			}
			newRangeStart := max(sourceRange.Start, item.SourceStart)
			newRangeLength := min(sourceRange.Start+sourceRange.Length, item.SourceStart+item.Length) - newRangeStart
			if newRangeLength > 0 {
				sourceRange = Range{Start: newRangeStart, Length: max(sourceRange.Start+sourceRange.Length, item.SourceStart+item.Length) - newRangeStart}
				destinationRanges = append(destinationRanges, Range{Start: item.DestinationStart, Length: newRangeLength})
			}
		}
		if sourceRange.Length > 0 {
			destinationRanges = append(destinationRanges, sourceRange)
		}
	}
	return destinationRanges
}

type Almanac struct {
	Seeds                 []int
	SeedToSoil            RangeMap
	SoilToFertilizer      RangeMap
	FertilizerToWater     RangeMap
	WaterToLight          RangeMap
	LightToTemperature    RangeMap
	TemperatureToHumidity RangeMap
	HumidityToLocation    RangeMap
}

type Range struct {
	Start  int
	Length int
}

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

func parseMap(scanner *bufio.Scanner) (RangeMap, error) {
	var rangeMap RangeMap
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		numbers, err := parseNumbers(line)
		if err != nil {
			return nil, err
		}
		if len(numbers) != 3 {
			return nil, fmt.Errorf("invalid line: Expected 3 numbers, got: %#v", numbers)
		}
		rangeMap = append(rangeMap, RangeMapItem{SourceStart: numbers[1], DestinationStart: numbers[0], Length: numbers[2]})
	}
	return rangeMap, nil
}

func getInput(path string) (Almanac, error) {
	file, err := os.Open(path)
	if err != nil {
		return Almanac{}, err
	}
	defer file.Close()

	var seeds []int
	var seedToSoil RangeMap
	var soilToFertilizer RangeMap
	var fertilizerToWater RangeMap
	var waterToLight RangeMap
	var lightToTemperature RangeMap
	var temperatureToHumidity RangeMap
	var humidityToLocation RangeMap

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		splitLine := strings.Split(line, ":")
		if len(splitLine) != 2 {
			return Almanac{}, fmt.Errorf("invalid line: Expected singular ':', but found: %#v", line)
		}

		key := strings.TrimSpace(splitLine[0])
		switch key {
		case "seeds":
			seeds, err = parseNumbers(splitLine[1])
			if err != nil {
				return Almanac{}, err
			}
		case "seed-to-soil map":
			seedToSoil, err = parseMap(scanner)
			if err != nil {
				return Almanac{}, err
			}
		case "soil-to-fertilizer map":
			soilToFertilizer, err = parseMap(scanner)
			if err != nil {
				return Almanac{}, err
			}
		case "fertilizer-to-water map":
			fertilizerToWater, err = parseMap(scanner)
			if err != nil {
				return Almanac{}, err
			}
		case "water-to-light map":
			waterToLight, err = parseMap(scanner)
			if err != nil {
				return Almanac{}, err
			}
		case "light-to-temperature map":
			lightToTemperature, err = parseMap(scanner)
			if err != nil {
				return Almanac{}, err
			}
		case "temperature-to-humidity map":
			temperatureToHumidity, err = parseMap(scanner)
			if err != nil {
				return Almanac{}, err
			}
		case "humidity-to-location map":
			humidityToLocation, err = parseMap(scanner)
			if err != nil {
				return Almanac{}, err
			}
		default:
			return Almanac{}, fmt.Errorf("invalid line: Unknown key %#v", key)
		}
	}
	if err := scanner.Err(); err != nil {
		return Almanac{}, err
	}

	return Almanac{
		Seeds:                 seeds,
		SeedToSoil:            seedToSoil,
		SoilToFertilizer:      soilToFertilizer,
		FertilizerToWater:     fertilizerToWater,
		WaterToLight:          waterToLight,
		LightToTemperature:    lightToTemperature,
		TemperatureToHumidity: temperatureToHumidity,
		HumidityToLocation:    humidityToLocation,
	}, nil
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
	// fmt.Printf("Args: %+v\n", args)

	almanac, err := getInput(args.InputPath)
	if err != nil {
		return err
	}

	switch args.Part {
	case 1:
		// Part 1
		fmt.Println("Seed -> Soil -> Fertilizer -> Water -> Light -> Temperature -> Humidity -> Location")
		seedToLocation := []RangeMap{almanac.SeedToSoil, almanac.SoilToFertilizer, almanac.FertilizerToWater, almanac.WaterToLight, almanac.LightToTemperature, almanac.TemperatureToHumidity, almanac.HumidityToLocation}
		var locations []int
		for _, seed := range almanac.Seeds {
			fmt.Printf("%v", seed)
			value := seed
			for _, rangeMap := range seedToLocation {
				value = rangeMap.getDestination(value)
				fmt.Printf(" -> %v", value)
			}
			fmt.Println()
			locations = append(locations, value)
		}
		fmt.Printf("\nLocations: %#v", locations)
		minLocation := slices.Min(locations)
		fmt.Printf("\nMin location: %v\n", minLocation)

	case 2:
		// Part 2
		// TODO
		fmt.Println("Seed -> Soil -> Fertilizer -> Water -> Light -> Temperature -> Humidity -> Location")
		seedToLocation := []RangeMap{almanac.SeedToSoil, almanac.SoilToFertilizer, almanac.FertilizerToWater, almanac.WaterToLight, almanac.LightToTemperature, almanac.TemperatureToHumidity, almanac.HumidityToLocation}
		var locationRanges []Range
		for i := 0; i < len(almanac.Seeds); i += 2 {
			seedRange := Range{Start: almanac.Seeds[i], Length: almanac.Seeds[i+1]}
			fmt.Printf("%v", seedRange)
			values := []Range{seedRange}
			for _, rangeMap := range seedToLocation {
				values = rangeMap.getDestinations(values)
				fmt.Printf(" -> %v", values)
			}
			fmt.Println()
			locationRanges = append(locationRanges, values...)
		}
		fmt.Printf("\nLocation ranges: %v", locationRanges)
		minLocations := make([]int, 0, len(locationRanges))
		for _, locationRange := range locationRanges {
			minLocations = append(minLocations, locationRange.Start)
		}
		minLocation := slices.Min(minLocations)
		fmt.Printf("\nMin location: %v\n", minLocation)
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
