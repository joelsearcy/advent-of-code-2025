package main

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Problem struct {
	lineNum               int
	buttons               []int
	desiredIndicatorState string
	joltageTarget         []int
}

type Result struct {
	lineNum        int
	indicatorPress int
	joltagePress   int
}

func main() {
	// Toggle between different solvers
	solverType := "partition" // default: partition-based
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--milp":
			solverType = "milp"
			fmt.Println("Using MILP solver (simplex + branch-and-bound)")
		case "--csp":
			solverType = "csp"
			fmt.Println("Using CSP solver (constraint propagation + DFS)")
		case "--partition":
			solverType = "partition"
			fmt.Println("Using Partition solver (aggressive button removal)")
		default:
			fmt.Println("Using Partition solver (aggressive button removal) - default")
		}
	} else {
		fmt.Println("Using Partition solver (aggressive button removal) - default")
	}

	data, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	// split the input data into lines
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")

	startTime := time.Now()

	maxJoltage := 0
	maxButtons := 0
	maxSlots := 0

	// lines are in the format "[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}"
	// we need to parse the square brackets at the front for the desired state of the indicator lights (zero indexed from left to right)
	// we can ignore the cruly braces at the end for now (part 1 only)
	// in between are the button definitions, each in parentheses, indicating which indicator lights each button toggles
	// initial state of all indicator lights is off (represented by '.')

	// Parse all problems first
	problems := make([]Problem, len(lines))
	for n, line := range lines {
		parts := strings.Split(line, " ")
		partsSize := len(parts)
		desiredIndicatorState := parts[0][1 : len(parts[0])-1]                   // remove square brackets
		buttonSchemas := parts[1 : partsSize-1]                                  // ignore curly braces section at the end
		desiredJoltageState := parts[partsSize-1][1 : len(parts[partsSize-1])-1] // remove curly braces

		numSlots := len(desiredIndicatorState)
		if numSlots > maxSlots {
			maxSlots = numSlots
		}
		numButtons := len(buttonSchemas)
		if numButtons > maxButtons {
			maxButtons = numButtons
		}
		buttons := make([]int, numButtons)
		for i, schema := range buttonSchemas {
			buttons[i] = 0
			buttonDef := schema[1 : len(schema)-1] // remove parentheses
			indices := strings.Split(buttonDef, ",")
			for _, idxStr := range indices {
				idx, _ := strconv.Atoi(idxStr)
				buttons[i] |= (1 << idx)
			}
		}

		joltageTarget := make([]int, numSlots)
		for i, numStr := range strings.Split(desiredJoltageState, ",") {
			num, _ := strconv.Atoi(numStr)
			joltageTarget[i] = num
			if num > maxJoltage {
				maxJoltage = num
			}
		}

		problems[n] = Problem{
			lineNum:               n,
			buttons:               buttons,
			desiredIndicatorState: desiredIndicatorState,
			joltageTarget:         joltageTarget,
		}
	}

	// Level 1 parallelization: process problems concurrently
	numWorkers := runtime.NumCPU()
	fmt.Printf("Processing %d problems with %d workers\n", len(problems), numWorkers)

	results := make([]Result, len(problems))
	var wg sync.WaitGroup
	problemChan := make(chan int, len(problems))

	var completed int64

	// Track in-progress problems for visibility
	inProgress := make([]int64, numWorkers)

	// Start workers
	for w := 0; w < numWorkers; w++ {
		workerID := w
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range problemChan {
				p := problems[idx]
				atomic.StoreInt64(&inProgress[workerID], int64(p.lineNum+1))
				problemStart := time.Now()

				indicatorPresses := FindMinButtonPressesForIndicatorTarget(p.buttons, p.desiredIndicatorState)

				var joltagePresses int
				switch solverType {
				case "milp":
					joltagePresses = FindMinButtonPressesMILP(p.buttons, p.joltageTarget)
				case "csp":
					joltagePresses = FindMinButtonPresses(p.buttons, p.joltageTarget)
				case "partition":
					joltagePresses = FindMinButtonPressesPartition(p.buttons, p.joltageTarget)
				}

				results[idx] = Result{
					lineNum:        p.lineNum,
					indicatorPress: indicatorPresses,
					joltagePress:   joltagePresses,
				}

				atomic.StoreInt64(&inProgress[workerID], 0)
				done := atomic.AddInt64(&completed, 1)

				// Show which problems are still in progress
				stillWorking := []int64{}
				for i := 0; i < numWorkers; i++ {
					if prob := atomic.LoadInt64(&inProgress[i]); prob > 0 {
						stillWorking = append(stillWorking, prob)
					}
				}
				duration := time.Since(problemStart)
				if duration > 2*time.Second {
					fmt.Printf("Completed %d/%d (problem %d) in %s | Still working: %v\n",
						done, len(problems), p.lineNum+1, duration, stillWorking)
				}
			}
		}()
	}

	// Send work
	for i := range problems {
		problemChan <- i
	}
	close(problemChan)

	wg.Wait()

	// Sum results
	totalPresses := int64(0)
	totalJoltagePresses := int64(0)
	impossibleCount := 0
	for _, r := range results {
		totalPresses += int64(r.indicatorPress)
		// Check for math.MaxInt (impossible case) to avoid overflow
		if r.joltagePress == math.MaxInt {
			impossibleCount++
			fmt.Printf("Warning: Problem %d returned impossible (math.MaxInt)\n", r.lineNum+1)
		} else {
			totalJoltagePresses += int64(r.joltagePress)
		}
	}

	if impossibleCount > 0 {
		fmt.Printf("WARNING: %d problems were impossible to solve!\n", impossibleCount)
	}

	fmt.Printf("Max joltage target found: %d\n", maxJoltage)
	fmt.Printf("Max buttons found: %d\n", maxButtons)
	fmt.Printf("Max slots found: %d\n", len(lines[0][1:strings.Index(lines[0], "]")]))
	fmt.Printf("Total minimum indicator button presses: %d\n", totalPresses)
	fmt.Printf("Total minimum joltage button presses: %d\n", totalJoltagePresses)
	fmt.Printf("Execution time: %s\n", time.Since(startTime))
}
