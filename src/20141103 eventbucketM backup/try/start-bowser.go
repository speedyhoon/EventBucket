package main

import (
//	"encoding/json"
	"fmt"
//	"os"
//	"reflect"
//	"text/template"
	"runtime"
	"os/exec"
)

func main() {
	var err error

	fmt.Printf("%v", runtime)

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", "http://localhost:4001/").Start()
	case "windows", "darwin":
//		err = exec.Command("open", "http://localhost:4001/").Start()
//		exec.Command(`C:\Windows\System32\rundll32.exe`, "url.dll,FileProtocolHandler", "http://localhost:4001/").Start()
		exec.Command(`rundll32.exe`, "url.dll,FileProtocolHandler", "http://localhost/").Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil{
		fmt.Printf("%v", err)
	}
}
