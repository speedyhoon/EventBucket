package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

// Converts numeric string to uint.
func stoU(id string) (uint, error) {
	u, err := strconv.ParseUint(strings.TrimSpace(id), 10, strconv.IntSize)
	if err != nil {
		log.Printf("Unable to convert %v to uint", id)
	}
	return uint(u), err
}

// mkDir attempts to create the path supplied if it doesn't exist.
func mkDir(path string) error {
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		err = os.Mkdir(path, os.ModeDir)
		if err != nil {
			log.Printf("Unable to create directory %v %v", path, err)
		}
	}
	return err
}
