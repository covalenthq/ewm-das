package sampler

import "math"

// CalculateSamplesNeeded calculates the minimum number of samples needed to achieve the desired confidence level.
func CalculateSamplesNeeded(R, E int, C float64) int {
	// Calculate the number of samples S using the given formula S \geq R - R \times \left(1 - C\right)^{\frac{1}{E}}
	S := float64(R) - float64(R)*math.Pow(1-C, 1/float64(E))
	return int(math.Ceil(S)) // Round up to the nearest whole number since you can't take a fraction of a sample
}
