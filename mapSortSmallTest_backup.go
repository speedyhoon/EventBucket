package main

import (
	"fmt"
	"sort"
)

type dataSlice []*data

type data struct {
	count int64
	size  int64
}

// Len is part of sort.Interface.
func (d dataSlice) Len() int {
	return len(d)
}

// Swap is part of sort.Interface.
func (d dataSlice) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

// Less is part of sort.Interface. We use count as the value to sort by
func (d dataSlice) Less(i, j int) bool {
	return d[i].count < d[j].count
}

func main() {
	M := map[string]*data {
		"x": {0, 0},
		"y": {2, 9},
		"z": {1, 7},
	}
	
	s := make(dataSlice, 0, len(M))

	for _, d := range M {
		s = append(s, d)
	}		
	
	// We just add 3 to one of our structs
	d := M["x"]
	d.count += 3
	
	sort.Sort(s)
	
	for _, d := range s {
		fmt.Printf("%+v\n", *d)
	}
}