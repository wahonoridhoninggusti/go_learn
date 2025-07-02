package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	result := CountWordFrequency("  Spaces,   tabs,\t\tand\nnew-lines are ignored!  ")
	fmt.Println(result)
}

func CountWordFrequency(text string) map[string]int {
	counts := make(map[string]int)
	re := regexp.MustCompile(`[\t\n\-]+`)
	cleaned := re.ReplaceAllString(text, " ")
	words := strings.SplitSeq(cleaned, " ")
	for t := range words {
		if t == "" {
			continue
		}
		t = strings.ToLower(t)
		re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
		output := re.ReplaceAllString(t, "")
		fmt.Println(output)
		counts[output]++
	}
	return counts
}
