package main

import (
	"strconv"
	"strings"
)

//Converts numeric string to uint
func stoU(id string) (uint, error) {
	u, err := strconv.ParseUint(strings.TrimSpace(id), 10, sysArch)
	if err != nil {
		warn.Printf("Unable to convert %v to uint", id)
	}
	return uint(u), err
}

func plural(length int, single, multiple string) string {
	if length != 1 {
		if multiple != "" {
			return multiple
		}
		return "s"
	}
	if single != "" {
		return single
	}
	return ""
}