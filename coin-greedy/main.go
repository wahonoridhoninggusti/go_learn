// package main

// import "fmt"

// func Reverse(arr []int) []int {
// 	result := make([]int, len(arr))
// 	for i, v := range arr {
// 		result[len(arr)-1-i] = v
// 	}
// 	return result
// }

// func CoinCombination(amount int, denominations []int) map[int]int {
// 	results := make(map[int]int)

// 	for _, coins := range denominations {
// 		if coins <= amount {
// 			mod := amount % coins
// 			res := (amount - mod) / coins
// 			results[coins] = res
// 			amount = mod
// 		}
// 	}

// 	if amount != 0 {
// 		return map[int]int{}
// 	}

// 	return results
// }

// func MinCoins(amount int, denominations []int) int {
// 	results := 0

// 	for _, coins := range denominations {
// 		if coins <= amount {
// 			mod := amount % coins
// 			res := (amount - mod) / coins
// 			results += res
// 			amount = mod
// 		}
// 	}

// 	if amount != 0 {
// 		return -1
// 	}

// 	return results
// }

// func main() {
// 	data := []int{1, 5, 10, 25, 50}
// 	dataDenomin := Reverse(data)
// 	fmt.Println(CoinCombination(87, dataDenomin))
// 	fmt.Println(CoinCombination(42, dataDenomin))
// 	fmt.Println(MinCoins(87, dataDenomin))
// }

package main

import "fmt"

func reversedCopy(slice []int) []int {
	result := make([]int, len(slice))
	for i, v := range slice {
		result[len(slice)-1-i] = v
	}
	return result
}

func CoinCombination(amount int, denominations []int) map[int]int {
	results := make(map[int]int)
	data := reversedCopy(denominations)
	for _, coins := range data {
		if coins <= amount {
			mod := amount % coins
			res := (amount - mod) / coins
			results[coins] = res
			amount = mod
		}
	}

	if amount != 0 {
		return map[int]int{}
	}

	return results
}

func MinCoins(amount int, denominations []int) int {
	results := 0
	data := reversedCopy(denominations)

	for _, coins := range data {
		if coins <= amount {
			mod := amount % coins
			res := (amount - mod) / coins
			results += res
			amount = mod
		}
	}

	if amount != 0 {
		return -1
	}

	return results
}

func main() {
	data := []int{1, 5, 10, 25, 50}

	fmt.Println(CoinCombination(87, data))
	fmt.Println(MinCoins(87, data))
}
