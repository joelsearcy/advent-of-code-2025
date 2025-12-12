package main

import (
	"math"
	"slices"
)

func FindMinButtonPresses(buttons []int, joltageTargets []int) int {
	numButtons := len(buttons)
	numTargets := len(joltageTargets)

	// Precompute: for each button, what's the max it can be pressed (limited by all slots it affects)
	buttonMaxPresses := make([]int, numButtons)
	for b := 0; b < numButtons; b++ {
		buttonMaxPresses[b] = math.MaxInt
		for s := 0; s < numTargets; s++ {
			if (buttons[b] & (1 << s)) != 0 {
				if joltageTargets[s] < buttonMaxPresses[b] {
					buttonMaxPresses[b] = joltageTargets[s]
				}
			}
		}
	}

	// Build slot info: which buttons affect each slot
	type slotInfo struct {
		index         int
		target        int
		affectingBtns []int
	}
	slots := make([]slotInfo, numTargets)
	for s := 0; s < numTargets; s++ {
		slots[s] = slotInfo{index: s, target: joltageTargets[s]}
		for b := 0; b < numButtons; b++ {
			if (buttons[b] & (1 << s)) != 0 {
				slots[s].affectingBtns = append(slots[s].affectingBtns, b)
			}
		}
	}

	slices.SortFunc(slots, func(a, b slotInfo) int {
		if len(a.affectingBtns) != len(b.affectingBtns) {
			return len(a.affectingBtns) - len(b.affectingBtns)
		}
		return b.target - a.target
	})

	// Check for impossible case: slot with no affecting buttons but positive target
	for _, slot := range slots {
		if len(slot.affectingBtns) == 0 && slot.target > 0 {
			return math.MaxInt
		}
	}

	// Button press counts - use -1 to indicate "not yet assigned"
	buttonPresses := make([]int, numButtons)
	for i := range buttonPresses {
		buttonPresses[i] = -1
	}

	minPresses := math.MaxInt

	// Helper: compute upper bound for a button based on remaining targets in unprocessed slots
	getButtonUpperBound := func(btn int, slotIdx int) int {
		maxAllowed := buttonMaxPresses[btn]
		for si := slotIdx; si < numTargets; si++ {
			s := slots[si].index
			if (buttons[btn] & (1 << s)) != 0 {
				// Current contribution from assigned buttons
				contrib := 0
				for _, b := range slots[si].affectingBtns {
					if buttonPresses[b] >= 0 {
						contrib += buttonPresses[b]
					}
				}
				remaining := joltageTargets[s] - contrib
				if remaining < maxAllowed {
					maxAllowed = remaining
				}
			}
		}
		return maxAllowed
	}

	// DFS through slots
	var dfs func(slotIdx int, currentPresses int)
	dfs = func(slotIdx int, currentPresses int) {
		if currentPresses >= minPresses {
			return
		}

		if slotIdx == numTargets {
			// Verify all slots satisfied
			for s := 0; s < numTargets; s++ {
				sum := 0
				for b := 0; b < numButtons; b++ {
					if (buttons[b] & (1 << s)) != 0 {
						if buttonPresses[b] < 0 {
							return // Unassigned button
						}
						sum += buttonPresses[b]
					}
				}
				if sum != joltageTargets[s] {
					return
				}
			}
			if currentPresses < minPresses {
				minPresses = currentPresses
			}
			return
		}

		slot := slots[slotIdx]

		// Calculate current contribution from assigned buttons
		// Must allocate fresh slice since recursive calls would corrupt shared buffer
		currentSum := 0
		var unassignedBtns []int
		for _, b := range slot.affectingBtns {
			if buttonPresses[b] >= 0 {
				currentSum += buttonPresses[b]
			} else {
				unassignedBtns = append(unassignedBtns, b)
			}
		}

		remaining := slot.target - currentSum
		if remaining < 0 {
			return // Already exceeded target
		}

		if len(unassignedBtns) == 0 {
			// All buttons assigned, check if constraint satisfied
			if remaining == 0 {
				dfs(slotIdx+1, currentPresses)
			}
			return
		}

		// Single unassigned button: it must take exactly 'remaining'
		if len(unassignedBtns) == 1 {
			b := unassignedBtns[0]
			maxAllowed := getButtonUpperBound(b, slotIdx+1)
			if remaining > maxAllowed {
				return // Can't assign this much
			}
			buttonPresses[b] = remaining
			dfs(slotIdx+1, currentPresses+remaining)
			buttonPresses[b] = -1
			return
		}

		// Multiple unassigned buttons: distribute 'remaining' among them
		// Must allocate fresh slice since tryDistribute calls dfs which could corrupt shared buffer
		upperBounds := make([]int, len(unassignedBtns))
		for i, b := range unassignedBtns {
			upperBounds[i] = getButtonUpperBound(b, slotIdx+1)
			if upperBounds[i] > remaining {
				upperBounds[i] = remaining
			}
		}

		var tryDistribute func(btnIdx int, remainingToDistribute int, pressesAdded int)
		tryDistribute = func(btnIdx int, remainingToDistribute int, pressesAdded int) {
			if currentPresses+pressesAdded >= minPresses {
				return
			}

			if btnIdx == len(unassignedBtns)-1 {
				// Last button takes whatever is remaining
				b := unassignedBtns[btnIdx]
				if remainingToDistribute > upperBounds[btnIdx] {
					return
				}
				buttonPresses[b] = remainingToDistribute
				dfs(slotIdx+1, currentPresses+pressesAdded+remainingToDistribute)
				buttonPresses[b] = -1
				return
			}

			b := unassignedBtns[btnIdx]
			maxForThis := upperBounds[btnIdx]
			if maxForThis > remainingToDistribute {
				maxForThis = remainingToDistribute
			}

			// Calculate minimum this button must take
			maxFromRest := 0
			for i := btnIdx + 1; i < len(unassignedBtns); i++ {
				maxFromRest += upperBounds[i]
			}
			minForThis := remainingToDistribute - maxFromRest
			if minForThis < 0 {
				minForThis = 0
			}

			for assign := minForThis; assign <= maxForThis; assign++ {
				buttonPresses[b] = assign
				tryDistribute(btnIdx+1, remainingToDistribute-assign, pressesAdded+assign)
				buttonPresses[b] = -1
			}
		}

		tryDistribute(0, remaining, 0)
	}

	dfs(0, 0)

	return minPresses
}

// FindMinButtonPressesPartition uses partition-based DFS with aggressive button removal
// Based on the Rust implementation that commits all buttons for a slot at once
func FindMinButtonPressesPartition(buttons []int, joltageTargets []int) int {
	numButtons := len(buttons)
	availableMask := (1 << numButtons) - 1 // All buttons initially available
	
	joltage := make([]int, len(joltageTargets))
	copy(joltage, joltageTargets)
	
	return dfsPartition(joltage, availableMask, buttons)
}

func dfsPartition(joltage []int, availableMask int, buttons []int) int {
	// Base case: all targets satisfied
	allZero := true
	for _, j := range joltage {
		if j != 0 {
			allZero = false
			break
		}
	}
	
	if allZero {
		return 0
	}

	// Find slot with lowest number of available affecting buttons (and highest target as tiebreaker)
	minSlot := -1
	minButtonCount := math.MaxInt
	maxTarget := -1
	
	for i, target := range joltage {
		if target == 0 {
			continue
		}
		
		// Count available buttons affecting this slot
		buttonCount := 0
		for b := 0; b < len(buttons); b++ {
			if (availableMask & (1 << b)) != 0 && (buttons[b] & (1 << i)) != 0 {
				buttonCount++
			}
		}
		
		if buttonCount < minButtonCount || (buttonCount == minButtonCount && target > maxTarget) {
			minButtonCount = buttonCount
			maxTarget = target
			minSlot = i
		}
	}
	
	if minSlot == -1 || minButtonCount == 0 {
		return math.MaxInt // No buttons available for a slot that needs them
	}
	
	// Get all available buttons that affect minSlot
	var matchingButtons []int
	for b := 0; b < len(buttons); b++ {
		if (availableMask & (1 << b)) != 0 && (buttons[b] & (1 << minSlot)) != 0 {
			matchingButtons = append(matchingButtons, b)
		}
	}
	
	// Create new mask excluding these buttons (they'll be committed)
	newMask := availableMask
	for _, b := range matchingButtons {
		newMask &= ^(1 << b)
	}
	
	result := math.MaxInt
	
	// Try all integer partitions of joltage[minSlot] across matchingButtons
	partition := make([]int, len(matchingButtons))
	partition[len(partition)-1] = joltage[minSlot]
	
	for {
		// Apply partition and check if valid
		newJoltage := make([]int, len(joltage))
		copy(newJoltage, joltage)
		
		valid := true
		totalPresses := 0
		
		for bi, presses := range partition {
			if presses == 0 {
				continue
			}
			
			btn := matchingButtons[bi]
			totalPresses += presses
			
			// Apply button presses to all affected slots
			for s := 0; s < len(joltage); s++ {
				if (buttons[btn] & (1 << s)) != 0 {
					if newJoltage[s] >= presses {
						newJoltage[s] -= presses
					} else {
						valid = false
						break
					}
				}
			}
			
			if !valid {
				break
			}
		}
		
		if valid {
			// Recurse with updated joltage and removed buttons
			subResult := dfsPartition(newJoltage, newMask, buttons)
			if subResult != math.MaxInt {
				candidate := totalPresses + subResult
				if candidate < result {
					result = candidate
				}
			}
		}
		
		// Generate next partition
		if !nextPartition(partition) {
			break
		}
	}
	
	return result
}

// nextPartition generates the next integer partition
// Given [a, b, c, d] that sums to N, generates all ways to partition N
func nextPartition(partition []int) bool {
	// Find rightmost non-zero element (not the last one)
	i := len(partition) - 1
	for i >= 0 && partition[i] == 0 {
		i--
	}
	
	if i <= 0 {
		return false
	}
	
	// Move one unit from position i to i-1, and accumulate rest at end
	val := partition[i]
	partition[i-1]++
	partition[i] = 0
	partition[len(partition)-1] = val - 1
	
	return true
}

type Constraint struct {
	coeffs []float64 // coefficients for each button variable
	rhs    float64   // right-hand side (target or bound)
	eqType int       // 0: ==, 1: >=, -1: <=
}

// FindMinButtonPressesMILP solves using Mixed Integer Linear Programming
// with simplex algorithm + branch-and-bound, falling back to CSP if needed
func FindMinButtonPressesMILP(buttons []int, joltageTargets []int) int {
	numButtons := len(buttons)
	numSlots := len(joltageTargets)

	// Build constraint matrix A and target vector
	// Constraints: for each slot s, sum of (button[i] affects slot s) * x[i] == target[s]
	// Also: x[i] >= 0 for all buttons
	// Objective: minimize sum of x[i]

	var constraints []Constraint

	// Equality constraints: one per slot
	for s := 0; s < numSlots; s++ {
		c := Constraint{
			coeffs: make([]float64, numButtons),
			rhs:    float64(joltageTargets[s]),
			eqType: 0, // equality
		}
		for b := 0; b < numButtons; b++ {
			if (buttons[b] & (1 << s)) != 0 {
				c.coeffs[b] = 1.0
			}
		}
		constraints = append(constraints, c)
	}

	// Non-negativity constraints: x[i] >= 0 (implicit in simplex)
	// Objective: minimize sum of all x[i] (coefficients all 1)
	objective := make([]float64, numButtons)
	for i := range objective {
		objective[i] = 1.0
	}

	// Solve using branch-and-bound
	result := branchAndBound(numButtons, objective, constraints)
	if result == nil {
		// MILP failed - fallback to CSP solver
		return FindMinButtonPresses(buttons, joltageTargets)
	}

	// Validate solution satisfies all constraints
	for s := 0; s < numSlots; s++ {
		sum := 0
		for b := 0; b < numButtons; b++ {
			if (buttons[b] & (1 << s)) != 0 {
				sum += int(math.Round(result[b]))
			}
		}
		if sum != joltageTargets[s] {
			// Solution doesn't satisfy constraints - fallback to CSP solver
			return FindMinButtonPresses(buttons, joltageTargets)
		}
	}

	// Sum up integer solution
	total := 0
	for _, v := range result {
		total += int(math.Round(v))
	}
	return total
}

const epsilon = 1e-9

// bbNode represents a node in the branch-and-bound tree
type bbNode struct {
	constraints []Constraint
	depth       int
	lpValue     float64 // LP relaxation value for priority ordering
}

// branchAndBound implements iterative branch-and-bound MILP solver with best-first search
func branchAndBound(numVars int, objective []float64, constraints []Constraint) []float64 {
	bestValue := math.Inf(1)
	var bestSolution []float64
	nodesExplored := 0
	maxNodes := 50000  // Increased significantly for better accuracy
	maxDepth := 50

	// Priority queue (min-heap) for best-first search - explore best LP values first
	pq := []bbNode{{constraints: constraints, depth: 0, lpValue: 0}}

	for len(pq) > 0 && nodesExplored < maxNodes {
		// Find and remove node with minimum lpValue (best-first)
		minIdx := 0
		for i := 1; i < len(pq); i++ {
			if pq[i].lpValue < pq[minIdx].lpValue {
				minIdx = i
			}
		}
		node := pq[minIdx]
		pq[minIdx] = pq[len(pq)-1]
		pq = pq[:len(pq)-1]

		nodesExplored++
		if node.depth > maxDepth {
			continue
		}

		// Solve LP relaxation
		value, solution := solveSimplex(numVars, objective, node.constraints)

		// Prune if infeasible or worse than current best
		if solution == nil || value >= bestValue-epsilon {
			continue
		}

		// Check if solution is integer and find most fractional variable
		fracVarIdx := -1
		maxFractionality := 0.0
		allInteger := true
		integerTolerance := 1e-6 // Slightly looser for floating point errors

		for i, v := range solution {
			rounded := math.Round(v)
			frac := math.Abs(v - rounded)
			if frac > integerTolerance {
				allInteger = false
				// Choose variable closest to 0.5 (most fractional)
				fractionality := math.Min(frac, 1.0-frac)
				if fractionality > maxFractionality {
					maxFractionality = fractionality
					fracVarIdx = i
				}
			}
		}

		if allInteger {
			// All integer - round and update best solution
			if value < bestValue {
				bestValue = value
				bestSolution = make([]float64, len(solution))
				for i, v := range solution {
					bestSolution[i] = math.Round(v)
				}
			}
			continue
		}

		if fracVarIdx == -1 {
			continue // Shouldn't happen, but be safe
		}

		// Branch on most fractional variable
		fracValue := solution[fracVarIdx]
		lowerBound := math.Floor(fracValue)
		upperBound := math.Ceil(fracValue)

		// Create two child nodes with estimated LP values
		// Lower bound branch: x[fracVarIdx] <= floor(fracValue)
		c1 := Constraint{
			coeffs: make([]float64, numVars),
			rhs:    lowerBound,
			eqType: -1, // <=
		}
		c1.coeffs[fracVarIdx] = 1.0
		lowerConstraints := make([]Constraint, len(node.constraints), len(node.constraints)+1)
		copy(lowerConstraints, node.constraints)
		lowerConstraints = append(lowerConstraints, c1)

		// Upper bound branch: x[fracVarIdx] >= ceil(fracValue)
		c2 := Constraint{
			coeffs: make([]float64, numVars),
			rhs:    upperBound,
			eqType: 1, // >=
		}
		c2.coeffs[fracVarIdx] = 1.0
		upperConstraints := make([]Constraint, len(node.constraints), len(node.constraints)+1)
		copy(upperConstraints, node.constraints)
		upperConstraints = append(upperConstraints, c2)

		// Estimate child LP values (use current value as estimate)
		// In best-first, these will be re-solved when popped from queue
		pq = append(pq, bbNode{
			constraints: lowerConstraints,
			depth:       node.depth + 1,
			lpValue:     value, // Will be refined when processed
		})
		pq = append(pq, bbNode{
			constraints: upperConstraints,
			depth:       node.depth + 1,
			lpValue:     value, // Will be refined when processed
		})
	}

	return bestSolution
}

// solveSimplex solves linear program using simplex algorithm
// Returns (objective value, solution vector) or (Inf, nil) if infeasible
func solveSimplex(numVars int, objective []float64, constraints []Constraint) (float64, []float64) {
	// Early feasibility check: if any constraint has all zero coefficients but non-zero RHS, infeasible
	for _, c := range constraints {
		allZero := true
		for _, coeff := range c.coeffs {
			if math.Abs(coeff) > epsilon {
				allZero = false
				break
			}
		}
		if allZero && math.Abs(c.rhs) > epsilon {
			return math.Inf(1), nil // Infeasible
		}
	}
	// Convert to standard form: minimize c^T x subject to Ax = b, x >= 0
	// We need to add slack/surplus variables for inequalities

	numOriginalVars := numVars
	slackVars := 0
	artificialVars := 0

	// Count slack and artificial variables needed
	for _, c := range constraints {
		if c.eqType != 0 {
			slackVars++
		}
		if c.eqType >= 0 {
			artificialVars++ // Need artificial variable for >= or ==
		}
	}

	totalVars := numOriginalVars + slackVars + artificialVars
	numConstraints := len(constraints)

	// Build tableau
	// Rows: constraints + objective row
	// Cols: variables + slack + artificial + RHS
	tableau := make([][]float64, numConstraints+1)
	for i := range tableau {
		tableau[i] = make([]float64, totalVars+1)
	}

	// Fill constraint rows
	slackIdx := numOriginalVars
	artificialIdx := numOriginalVars + slackVars

	for i, c := range constraints {
		// Original variables
		copy(tableau[i][:numOriginalVars], c.coeffs)

		// Slack/surplus variables
		if c.eqType == -1 {
			// <= constraint: add slack variable
			tableau[i][slackIdx] = 1.0
			slackIdx++
		} else if c.eqType == 1 {
			// >= constraint: subtract surplus variable, add artificial
			tableau[i][slackIdx] = -1.0
			slackIdx++
			tableau[i][artificialIdx] = 1.0
			artificialIdx++
		} else {
			// == constraint: add artificial variable
			tableau[i][artificialIdx] = 1.0
			artificialIdx++
		}

		// RHS
		tableau[i][totalVars] = c.rhs
	}

	// Objective row (last row)
	for i, coeff := range objective {
		tableau[numConstraints][i] = -coeff // Negative for minimization
	}

	// Basic feasible solution: artificials and slacks are basic
	// Need to eliminate artificials from objective if present
	if artificialVars > 0 {
		// Phase 1: minimize sum of artificial variables
		// Not implementing full two-phase simplex for brevity
		// For this problem, we assume constraints are feasible
	}

	// Simplex algorithm with iteration limit
	maxIterations := 500
	iteration := 0
	
	for iteration < maxIterations {
		iteration++
		
		// Find entering variable (most negative coefficient in objective row)
		enteringCol := -1
		minCoeff := 0.0
		for j := 0; j < totalVars; j++ {
			if tableau[numConstraints][j] < minCoeff-epsilon {
				minCoeff = tableau[numConstraints][j]
				enteringCol = j
			}
		}

		if enteringCol == -1 {
			// Optimal solution found
			break
		}

		// Find leaving variable (minimum ratio test with tie-breaking)
		leavingRow := -1
		minRatio := math.Inf(1)
		for i := 0; i < numConstraints; i++ {
			if tableau[i][enteringCol] > epsilon {
				ratio := tableau[i][totalVars] / tableau[i][enteringCol]
				if ratio < minRatio-epsilon || (math.Abs(ratio-minRatio) < epsilon && leavingRow > i) {
					minRatio = ratio
					leavingRow = i
				}
			}
		}

		if leavingRow == -1 {
			// Unbounded
			return math.Inf(1), nil
		}

		// Pivot
		pivot := tableau[leavingRow][enteringCol]
		for j := 0; j <= totalVars; j++ {
			tableau[leavingRow][j] /= pivot
		}

		for i := 0; i <= numConstraints; i++ {
			if i != leavingRow {
				factor := tableau[i][enteringCol]
				for j := 0; j <= totalVars; j++ {
					tableau[i][j] -= factor * tableau[leavingRow][j]
				}
			}
		}
	}
	
	// Check if hit iteration limit
	if iteration >= maxIterations {
		return math.Inf(1), nil // Treat as infeasible
	}

	// Extract solution (only original variables)
	solution := make([]float64, numOriginalVars)
	for j := 0; j < numOriginalVars; j++ {
		// Check if this variable is basic
		isBasic := false
		basicRow := -1
		for i := 0; i < numConstraints; i++ {
			if math.Abs(tableau[i][j]-1.0) < epsilon {
				// Check if all other entries in this column are 0
				allZero := true
				for k := 0; k < numConstraints; k++ {
					if k != i && math.Abs(tableau[k][j]) > epsilon {
						allZero = false
						break
					}
				}
				if allZero {
					isBasic = true
					basicRow = i
					break
				}
			}
		}

		if isBasic {
			solution[j] = math.Max(0, tableau[basicRow][totalVars]) // Ensure non-negative
		} else {
			solution[j] = 0.0
		}
	}

	objectiveValue := -tableau[numConstraints][totalVars] // Negative because we negated it

	return objectiveValue, solution
}
