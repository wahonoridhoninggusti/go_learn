package main

import "fmt"

func NaiveSearch(text, s string) []int {
	var positions []int
	n := len(text)
	m := len(s)

	for i := 0; i <= n-m; i++ {
		temp := 0
		for j := range m {
			if text[i+j] != s[j] {
				break
			}
			temp++
		}

		if temp == m {
			positions = append(positions, i)
		}
	}
	return positions
}

func ComputeLPS(s string) []int {
	lps := make([]int, len(s))

	for i := 1; i < len(s); i++ {
		length := lps[i-1]
		for length > 0 && s[i] != s[length] {
			length = lps[length-1]
		}
		if s[i] == s[length] {
			length++
		}
		lps[i] = length
	}
	return lps
}

func KMPSearch(text, s string) []int {
	komput := ComputeLPS(s)
	var res []int
	i, j := 0, 0

	for i < len(text) {
		if text[i] == s[j] {
			i++
			j++
		}
		if j == len(s) {
			res = append(res, i-j)
			j = komput[j-1]
		} else if i < len(text) && text[i] != s[j] {
			if j != 0 {
				j = komput[j-1]
			} else {
				i++
			}
		}

	}
	return res
}

const base = 256
const prime = 101

func rabinKarp(text, s string) []int {
	n := len(text)
	m := len(s)
	if m > n {
		return nil
	}

	var res []int

	pHash := 0
	tHash := 0
	h := 1

	for i := 0; i < m-1; i++ {
		h = (h * base) % prime
	}

	for i := 0; i < m; i++ {
		pHash = (base*pHash + int(s[i])) % prime
		tHash = (base*tHash + int(text[i])) % prime

	}

	for i := 0; i <= n-m; i++ {
		if pHash == tHash {
			match := true
			for j := 0; j < m; j++ {
				if text[i+j] != s[j] {
					match = false
					break
				}
			}
			if match {
				res = append(res, i)
			}
		}

		if i < n-m {
			tHash = (base*(tHash-int(text[i])*h) + int(text[i+m])) % prime
			if tHash < 0 {
				tHash += prime
			}
		}
	}
	fmt.Println(res)
	return res
}
func main() {
	fmt.Println(KMPSearch("abababc", "abaabaa"))
	NaiveSearch("abababc", "abababd")

	rabinKarp("abababababaabss", "aba")
}
