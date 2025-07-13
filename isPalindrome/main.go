package main

import (
	"fmt"
	"regexp"
	"strings"
)

func isPalindrome(s string) bool {
	tolower := strings.ToLower(s)
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	clean := re.ReplaceAllString(tolower, "")
	runes := []rune(clean)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return clean == string(runes)
}

func main() {
	fungsi := isPalindrome("wahono, rig")

	fmt.Println(fungsi)
}
