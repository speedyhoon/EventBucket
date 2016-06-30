package main

import (
	"strconv"
	"strings"
)

//Converts numeric string to uint
func stoU(id string) (uint, error) {
	u, err := strconv.ParseUint(strings.TrimSpace(id), 10, 64)
	if err != nil {
		warn.Printf("Unable to convert %v to uint", id)
	}
	return uint(u), err
}
