package helpers

import "math"

// RoundFloat allows for rounding of a float to a given number of decimal places.
func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
