package main

import (
	"fmt"
	"sort"
)

/*func main() {

	slice := []int{5, 6, 7, 2, 1, 0}
	fmt.Println("--- Unsorted --- \n\n", slice)
	quicksort(slice)
	fmt.Println("--- Sorted ---\n\n", slice, "\n")
}*/

func sortList(a map[string]int64) []string {

	keys := make([]string, 0, len(a))

	for key := range a {
		keys = append(keys, key)
	}

	fmt.Println(a)
	fmt.Println(keys)

	sort.SliceStable(keys, func(i, j int) bool {
		return a[keys[i]] < a[keys[j]]
	})

	return keys
}
