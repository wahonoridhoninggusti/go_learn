package main

import (
	"fmt"
	"sort"
	"strings"
)

func SlowSort(data []int) []int {
	results := make([]int, len(data))
	copy(results, data)

	for range data {
		for j := 0; j < len(data)-1; j++ {
			if results[j] > results[j+1] {
				results[j], results[j+1] = results[j+1], results[j]
			}
		}
	}
	return results
}

func OptimizedSort(data []int) []int {
	results := make([]int, len(data))
	copy(results, data)
	sort.Ints(results)
	return results
}

func InefficientStringBuilder(parts []string, repeatCount int) string {
	result := ""

	for i := 0; i < repeatCount; i++ {
		for _, part := range parts {
			result += part
		}
	}
	return result
}

func OptimizedStringBuilder(parts []string, repeatCount int) string {
	var sb strings.Builder

	totalLen := 0
	for _, part := range parts {
		totalLen = len(part)
	}
	sb.Grow(totalLen * repeatCount)
	for i := 0; i < repeatCount; i++ {
		for _, part := range parts {
			sb.WriteString(part)
		}
	}
	return sb.String()
}

func ExpensiveCalculation(n int) int {
	if n <= 0 {
		return 0
	}
	sum := 0
	for i := 0; i <= n; i++ {
		sum += fibonacci(i)
	}

	return sum
}

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func OptimizedCalculation(n int) int {
	if n <= 0 {
		return 0
	}
	sum := 0
	for i := 0; i <= n; i++ {
		sum += fibbonacci(i)
	}

	return sum
}

var memo = map[int]int{}

func fibbonacci(n int) int {
	if n <= 1 {
		return n
	}

	if val, ok := memo[n]; ok {
		return val
	}

	memo[n] = fibbonacci(n-1) + fibbonacci(n-2)
	return memo[n]
}

func HighAllocationSearch(text, substr string) map[int]string {
	result := make(map[int]string)

	lowerText := strings.ToLower(text)
	lowerSubstr := strings.ToLower(substr)

	for i := 0; i < len(lowerText); i++ {
		if i+len(lowerSubstr) <= len(lowerText) {
			potentialMatch := lowerText[i : i+len(lowerSubstr)]
			if potentialMatch == lowerSubstr {
				result[i] = text[i : i+len(substr)]
			}
		}
	}

	return result
}

func OptimizedSearch(text, substr string) map[int]string {
	result := make(map[int]string)
	if len(substr) == 0 || len(text) < len(substr) {
		return result
	}

	textLower := strings.ToLower(text)
	substrLower := strings.ToLower(substr)
	substrLen := len(substrLower)

	for i := 0; i <= len(textLower)-substrLen; i++ {
		if textLower[i:i+substrLen] == substrLower {
			result[i] = text[i : i+substrLen]
		}
	}

	return result

}

func main() {
	fmt.Println(OptimizedSort([]int{1, 4, 7, 8, 3, 0, 3, 5}))
	fmt.Println(InefficientStringBuilder([]string{"Hello", " ", "World"}, 10))
	fmt.Println(OptimizedCalculation(1))
}
