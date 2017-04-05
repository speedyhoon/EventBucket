package main

import (
	"os/exec"
)

func openBrowser(url string){
	if exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", url).Start() != nil {
		warn.Println("Unable to open a web browser for", url)
	}
}