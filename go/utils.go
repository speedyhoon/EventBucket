package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

//Converts numeric string to uint
func stoU(id string) (uint, error) {
	u, err := strconv.ParseUint(strings.TrimSpace(id), 10, sysArch)
	if err != nil {
		log.Printf("Unable to convert %v to uint", id)
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

//mkDir attempts to create the path supplied if it doesn't exist.
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

// startBrowser tries to open the URL in a browser, and returns whether it succeed.
func openBrowser(url string) bool {
	var args []string
	switch runtime.GOOS {
	case "darwin":
		//macOS, iOS
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		//android, dragonfly, freebsd, linux, nacl, netbsd, openbsd, plan9, solaris
		args = []string{"xdg-open"}
	}
	return exec.Command(args[0], append(args[1:], url)...).Start() == nil
}

func trimFloat(num float32) string {
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.6f", num), "0"), ".")
}
