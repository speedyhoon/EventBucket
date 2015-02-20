package main

import (
	"bytes"
	"io/ioutil"
	"time"
	"fmt"
)

func main(){
	fileName := "settings.go"
	//TODO Move these settings from hardcoded to a yaml file or something
	versionNumber := 58
	source, err := ioutil.ReadFile(fileName)
	if err == nil{
		date := time.Now().Format("January 2, 2006")
		newSettings := []byte(fmt.Sprintf("VERSION = 0.%v\n\tBUILDDATE = \"Compiled on %v  by Cam Webb\"", versionNumber, date))
		source = bytes.Replace(source, []byte("/*IMPORT EXTERNAL SETTINGS HERE*/"), []byte(newSettings), 1)
		fmt.Printf("version: %v, date: %v", versionNumber, date)
		ioutil.WriteFile(fileName, source, 0777)
	}
}
