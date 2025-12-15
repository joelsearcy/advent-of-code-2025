package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// define an enum of type of values: "stops", "passes"
type CountMethod string

const (
	Stops  CountMethod = "stops"
	Passes CountMethod = "passes"
)

func main() {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	startTime := time.Now()

	// Split the input data into lines
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")

	// Starting position of the padlock
	ringSize := 100
	position := 50
	countOfZeros := 0
	countMethod := Passes
	landedOnZero := 0
	visitedZero := 0

	for _, line := range lines {
		if len(line) < 2 {
			continue
		}

		// each line represents a number of clicks on a padlock, preceded by a direction (L or R)
		direction := line[0]
		clicks, err := strconv.Atoi(line[1:])
		if err != nil {
			panic(err)
		}

		unit := 1
		if direction == 'L' {
			unit = -1
		}

		adjustedClicks := clicks % ringSize
		passesZero := clicks / ringSize
		if (unit == -1 && position > 0 && position <= adjustedClicks) || position+unit*adjustedClicks >= ringSize {
			passesZero += 1
		}
		newPosition := (position + unit*adjustedClicks + ringSize) % ringSize
		// fmt.Printf("Direction: %c, Clicks: %d, From: %d To: %d PassesZero: %d\n", direction, clicks, position, newPosition, passesZero)

		if newPosition == 0 {
			landedOnZero += 1
		}
		visitedZero += passesZero

		position = newPosition

		if countMethod == Stops && position == 0 {
			countOfZeros++
		} else if countMethod == Passes {
			countOfZeros += passesZero
		}
	}
	fmt.Printf("Final position: %d\n", position)
	fmt.Printf("landedOnZero: %d\n", landedOnZero)
	fmt.Printf("visitedZero: %d\n", visitedZero)
	fmt.Printf("Number of times landed on or passing 0: %d\n", countOfZeros)
	fmt.Printf("Execution time: %s\n", time.Since(startTime))
}
