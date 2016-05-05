package main

import (
	"strconv"
	"strings"
)

//Converts base36 string to uint
func b36tou(id string) (uint, error) {
	u, err := strconv.ParseUint(strings.TrimSpace(id), 36, 64)
	if err != nil {
		warn.Printf("Unable to convert %v to uint", id)
	}
	return uint(u), err
}

//Converts numeric string to uint
//TODO rename stoU
func strToUint(id string) (uint, error) {
	u, err := strconv.ParseUint(strings.TrimSpace(id), 10, 64)
	if err != nil {
		warn.Printf("Unable to convert %v to uint", id)
	}
	return uint(u), err
}
