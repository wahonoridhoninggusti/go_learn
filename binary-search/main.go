package main

import (
	"fmt"
)

func BinarySearch(arr []int, target int) int {
	low := 0
	high := len(arr) - 1

	for low <= high {
		mid := (low + high) / 2
		if target == arr[mid] {
			return mid
		} else if arr[mid] < target {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return -1
}

func BinarySearchRecursive(arr []int, target int, left int, right int) int {
	if left > right {
		return -1
	}

	mid := (left + right) / 2
	if arr[mid] < target {
		return BinarySearchRecursive(arr, target, mid+1, right)
	} else if arr[mid] > target {
		return BinarySearchRecursive(arr, target, left, mid-1)
	} else {
		return mid
	}
}

func FindInsertPosition(arr []int, target int) int {
	low, high := 0, len(arr)-1
	for low <= high {
		mid := (low + high) / 2
		if arr[mid] < target {
			low = mid + 1
		} else if arr[mid] > target {
			high = mid - 1
		} else {
			return low
		}
	}

	return low
}

func main() {
	fmt.Println(BinarySearch([]int{1, 3, 5, 7, 9}, 7))
	fmt.Println(FindInsertPosition([]int{1, 3, 5, 7, 9}, 6))
	fmt.Println(BinarySearchRecursive([]int{1, 3, 5, 7, 9}, 7, 0, 4))
}
