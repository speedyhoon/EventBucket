package main

import (
	"strconv"
	"strings"
)

//Converts base36 string to uint
func b36tou(id string) (uint, error) {
	t.Println("64 bit version")
	u, err := strconv.ParseUint(strings.TrimSpace(id), 36, 64)
	return uint(u), err
}

//Converts base36 string to uint
func strToUint(id string) (uint, error) {
	t.Println("64 bit version")
	u, err := strconv.ParseUint(strings.TrimSpace(id), 10, 64)
	return uint(u), err
}
