package main

import (
	"fmt"
	"math"
)

func CelciusToFahrenheit(c float64) float64 {
	return (c * 9 / 5) + 32
}

func FahrenheitToCelcius(f float64) float64 {
	return (f - 32) * 5 / 9
}

func Round(value float64, decimals int) float64 {
	multiplier := math.Pow(10, float64(decimals))
	return math.Round(value*multiplier) / multiplier
}

func main() {
	fmt.Println(FahrenheitToCelcius(30))
}
