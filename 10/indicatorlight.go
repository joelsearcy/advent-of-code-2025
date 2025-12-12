package main

import (
	"math"
)

func FindMinButtonPressesForIndicatorTarget(buttons []int, desiredIndicatorState string) int {
	numButtons := len(buttons)

	// convert desired state to binary representation for easier comparison
	// '.' = 0 (off), '#' = 1 (on)
	desiredIndicatorStateBinary := 0
	for i, ch := range desiredIndicatorState {
		if ch == '#' {
			desiredIndicatorStateBinary |= (1 << i)
		}
	}
	// fmt.Printf("Desired state: %s (binary: %b)\n", desiredIndicatorState, desiredIndicatorStateBinary)

	minPresses := math.MaxInt
	// try all combinations of button presses (2^numButtons possibilities)
	// tracking the number of presses for each combination to find the matching desired state with the fewest presses
	for combo := 0; combo < (1 << numButtons); combo++ {
		currentState := 0
		pressCount := 0
		for b := 0; b < numButtons; b++ {
			if (combo & (1 << b)) != 0 {
				currentState ^= buttons[b]
				pressCount++
			}
		}
		if currentState == desiredIndicatorStateBinary {
			if pressCount < minPresses {
				minPresses = pressCount
			}
			continue
		}
	}

	return minPresses
}
