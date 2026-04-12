package utils

const positionStep = 1000.0

// CalcPosition returns the position value for an item being placed
// between lower and upper.
//
//	upper = nil  → append at end: lower + positionStep
//	upper = 0.0  → container is empty: positionStep
//	otherwise    → midpoint between lower and upper
func CalcPosition(lower float64, upper *float64) float64 {
	if upper == nil {
		return lower + positionStep
	}
	return (lower + *upper) / 2
}
