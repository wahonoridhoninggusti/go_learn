package main

import "fmt"

func main() {
	graph := map[int][]int{
		0: {1, 2},
		1: {2, 3},
		2: {3},
		3: {4},
		4: {},
	}
	queries := []int{0, 1, 2}
	// FindLastNode(graph)
	// ConcurentBFSGraph(graph, queries, 2)
	result := ConcurrentBFSQueries(graph, queries, 2)
	fmt.Println(result)
}

// func FindLastNode(graph map[int][]int) {
// 	max := -1
// 	for node := range graph {
// 		if node > max {
// 			max = node
// 		}
// 	}
// 	// fmt.Println(max)
// 	return max
// }

// func ConcurentBFSGraph(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	result := make(map[int][]int)
	// var semuaJalur [][]int
	// var queue [][]int
	max := -1

	for node := range graph {
		if node > max {
			max = node
		}
	}

	for _, s := range queries {
		visited := make(map[int]bool)
		queue := []int{s}
		var cross []int

		for len(queue) > 0 {
			path := queue[0]
			queue = queue[1:]
			if visited[path] {
				continue
			}

			visited[path] = true
			cross = append(cross, path)

			for _, neighbor := range graph[path] {
				if !visited[neighbor] {
					queue = append(queue, neighbor)
				}
			}

		}
		result[s] = cross
	}

	// for len(queue) > 0 {
	// 	path := queue[0]
	// 	queue = queue[1:]
	// 	last := path[len(path)-1]
	// 	if last == max {
	// 		semuaJalur = append(semuaJalur, path)
	// 		continue
	// 	}

	// 	for _, neighbor := range graph[last] {
	// 		if IsContain(path, neighbor) {
	// 			continue
	// 		}

	// 		newPath := append([]int{}, path...)
	// 		newPath = append(newPath, neighbor)
	// 		queue = append(queue, newPath)
	// 	}
	// }
	return result
}

// func IsContain(path []int, node int) bool {
// 	for _, n := range path {
// 		if n == node {
// 			return true
// 		}
// 	}
// 	return false
// }
