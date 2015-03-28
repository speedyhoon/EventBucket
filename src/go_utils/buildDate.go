package main

import (
	"bytes"
	"io/ioutil"
	"time"
	"fmt"
	"path/filepath"
	"strings"
	"os"
)
//loop through all files & sub directories and replace the search variables with the actual values

const (
	rootDir = "../eventbucketM/"
	copyToDir = rootDir +"!/"
	//Folder structure
	dirCSS  = "/c/"
	dirJPEG = "/e/"
	dirJS   = "/j/"	//TODO It might be better to replace the constants directly with the string value. e.g __dirCSS__ and replace with "/c/"
	dirPNG  = "/p/"
	dirSVG  = "/v/"
	dirWEBP  = "/w/"
	favicon = dirPNG+"a" //TODO Create a custom icon with sizes 16x16, 32, 48, 64 and 128
	dirWOF  = "/f/"
	dirWOF2 = "/2/"

	dbArguments = `ebd", "--dbpath", databasePath, "--port", "38888", "--nssize", "1", "--smallfiles", "--noscripting", `
)

var (
	currentDir = ""
	dev = "true"
	replaceChars = map[string]interface{}{
		"VersionNumber": 58,		//TODO get the Git tag from the last commit
		"dirRoot": "",
		"dirCSS": dirCSS,
		"dirJPEG": dirJPEG,
		"dirJS": dirJS,
		"dirPNG": dirPNG,
		"dirSVG": dirSVG,
		"dirWEBP": dirWEBP,
		"dirWOF": dirWOF,
		"dirWOF2": dirWOF2,
		"Favicon": favicon,
		"BuildDate": time.Now().Format("January 2, 2006"),
		"schemaSHOOTER": "S",
		"schemaAutoInc": "U",
		"schemaRANGE": "R",
		"schemaSORT": "o",
		"schemaGRADES": "g",
	}
	devMode = map[string]interface{}{
		"DbArgs": dbArguments + `"--noauth", "--slowms", "3", "--cpu", "--profile", "2", "--objcheck", "--notablescan", "--rest`,
		//"NewRelic": "true",		//TODO there has to be a better way to do this? Maybe use Gulp.js instead?
	}

	prodMode = map[string]interface{}{
		"DbArgs": dbArguments + `"--nohttpinterface`,
	}
)
//TODO eventually move these settings to a json or yaml file.

func main(){
	joinSettings()
	loadHtmlSnippets()
	//TODO get the version number frim Git "C:\Program Files (x86)\Git\bin\git.exe" describe --tags
	filepath.Walk(rootDir + "golang", walkPath)
//	filepath.Walk(rootDir + "sass", walkPath)
	filepath.Walk(rootDir + "js", walkPath)
	filepath.Walk(rootDir + "html", walkPath)
//	filepath.Walk(rootDir + "htm", walkPath)
}

func joinSettings(){
	additionalSettings := prodMode
	if dev != ""{
		additionalSettings = devMode
	}
	for name, setting := range additionalSettings{
		replaceChars[name] = setting
	}
}

func walkPath(path string, f os.FileInfo, err error) error {
	if err != nil{
		return err
	}
	if !f.IsDir() {
		source, _ := ioutil.ReadFile(path)
		if err == nil{
			source = replaceContents(replaceChars, source)
			switch currentDir{
			case "golang":
				if f.Name() != "unused.go" {
					err = ioutil.WriteFile(copyToDir+f.Name(), source, 0777)
				}
			case "js":
				err = ioutil.WriteFile(copyToDir+"j/"+f.Name(), source, 0777)
			case "html":
				err = ioutil.WriteFile(copyToDir+"/h/"+strings.Replace(f.Name(), ".html", ".htm", -1), minifyHtml(f.Name(), source), 0777)
				err = ioutil.WriteFile(copyToDir+"/html/"+f.Name(), source, 0777)
			}
			if err != nil{
				fmt.Printf("ERROR: %v", err)
			}
		}
	}else{
		currentDir = f.Name()
	}
	return nil
}

func replaceContents(replaceSearch map[string]interface{}, source []byte)[]byte{
	for search, replace := range replaceSearch {
		source = bytes.Replace(source, []byte("^^"+search+"^^"), []byte(fmt.Sprintf("%v",replace)), -1)
	}
	return source
}

func loadHtmlSnippets(){
	//TODO remove this hacky code when multiple templates can be used easily with Ace!
	var fileContents []byte
	var err error

	fileContents, err = ioutil.ReadFile(rootDir + "html/NetworkAdaptor.html")
	replaceChars["NetworkAdaptor"] = string(replaceContents(replaceChars, fileContents)[:])
	if err != nil{
		fmt.Println("Unable to load NetworkAdaptor html contents")
	}

	if exists(replaceChars, "NewRelic") != "" {
		fileContents, err = ioutil.ReadFile(rootDir + "html/newRelic.html")
		if err != nil{
			fmt.Println("Unable to load NewRelic html contents")
		}
		replaceChars["NewRelic"] = string(replaceContents(replaceChars, fileContents)[:])
	}else {
		replaceChars["NewRelic"] = ""
	}
}

func exists(dict map[string]interface{}, key string) string {
	if val, ok := dict[key]; ok {
		return fmt.Sprintf("%v", val)
	}
	return ""
}

func minifyHtml(pageName string, minify []byte) []byte {
	if bytes.Contains(minify, []byte("ZgotmplZ")) {
		fmt.Println("Template generation error: ZgotmplZ")
		return []byte("")
	}
	removeChars := []string{
		"	", //tab
		"\n", //new line
		"\r", //carriage return
	}
	//TODO remove spaces between block elements like: </div> <div> but keep between inline elements like </span> <span>
	//TODO use improved regex for better searching & replacement
	replaceChars := map[string]string{
		"  ":            " ", //double spaces
		"type=text":     "",
		"type=\"text\"": "",
		"type='text'":   "",
		" >":            ">",
		"< ":            "<",
		">  <":          "> <",
		" />":           "/>",
		"/ >":           "/>",
		"<br/>":         "<br>",
		"</br>":         "<br>",
		"<br />":        "<br>",
	}
	for _, search := range removeChars {
		minify = bytes.Replace(minify, []byte(search), []byte(""), -1)
	}
	for search, replace := range replaceChars {
		length := len(minify)
		minify = bytes.Replace(minify, []byte(search), []byte(replace), -1)
		if length != len(minify) {
			fmt.Printf("A dodgy character (%v) was found in '%v'! Please replace with (%v).", search, pageName, replace)
		}
	}
	return minify
}
