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
	ROOT_DIR = "../eventbucketM/"
	COPY_TO_DIR = ROOT_DIR +"!/"
	//Folder structure
	DIR_CSS  = "/c/"
	DIR_JPEG = "/e/"
	DIR_JS   = "/j/"	//TODO It might be better to replace the constants directly with the string value. e.g __DIR_CSS__ and replace with "/c/"
	DIR_PNG  = "/p/"
	DIR_SVG  = "/v/"
	DIR_WEBP  = "/w/"
	FAVICON = DIR_PNG+"a" //TODO Create a custom icon with sizes 16x16, 32, 48, 64 and 128

	BD_ARGS = `--dbpath", databasePath, "--port", "38888", "--nssize", "1", "--smallfiles", "--noscripting", `
)

var (
	CURRENT_DIR = ""
	DEV = "true"
	ReplaceChars = map[string]interface{}{
		"VersionNumber": 58,		//TODO get the Git tag from the last commit
		"DIR_ROOT": "",
		"DirCss": DIR_CSS,
		"DirJpeg": DIR_JPEG,
		"DirJs": DIR_JS,
		"DirPng": DIR_PNG,
		"DirSvg": DIR_SVG,
		"DirWebp": DIR_WEBP,
		"Favicon": FAVICON,
		"BuildDate": time.Now().Format("January 2, 2006"),
		"schemaSHOOTER": "S",
		"schemaAutoInc": "U",
		"schemaRANGE": "R",
		"schemaSORT": "o",
		"schemaGRADES": "g",
		"NetworkAdaptor": "true",		//TODO there has to be a better way to do this?
	}
	DevMode = map[string]interface{}{
		"DbArgs": BD_ARGS + `"--noauth", "--slowms", "3", "--cpu", "--profile", "2", "--objcheck", "--notablescan", "--rest`,
		"NewRelic": "true",		//TODO there has to be a better way to do this?
	}

	ProdMode = map[string]interface{}{
		"DbArgs": BD_ARGS + `"--nohttpinterface`,
	}
)
//TODO eventually move these settings to a json or yaml file.

func main(){
	var err error
	if exists(DevMode, "NewRelic") {
		DevMode["NewRelic"], err = ioutil.ReadFile("html/newRelic.html")
	}else{
		ProdMode["NewRelic"] = ""
	}
	if exists(ReplaceChars, "NetworkAdaptor") {
		ReplaceChars["NetworkAdaptor"], _ := ioutil.ReadFile("html/NetworkAdaptor.html")
	}
	joinSettings()
	filepath.Walk(ROOT_DIR + "golang", walkPath)
//	filepath.Walk(ROOT_DIR + "sass", walkPath)
//	filepath.Walk(ROOT_DIR + "js", walkPath)
	filepath.Walk(ROOT_DIR + "html", walkPath)
}

func joinSettings(){
	additionalSettings := ProdMode
	if DEV != ""{
		additionalSettings = DevMode
	}
	for name, setting := range additionalSettings{
		ReplaceChars[name] = setting
	}
}

func walkPath(path string, f os.FileInfo, err error) error {
	if err != nil{
		return err
	}
	if !f.IsDir() {
		source, _ := ioutil.ReadFile(path)
		if err == nil{
			for search, replace := range ReplaceChars {
				source = bytes.Replace(source, []byte("^^"+search+"^^"), []byte(fmt.Sprintf("%v",replace)), -1)
			}
		}
		switch CURRENT_DIR{
		case "golang":
			ioutil.WriteFile(COPY_TO_DIR+f.Name(), source, 0777)
			break
		case "html":
			ioutil.WriteFile(COPY_TO_DIR+"/h/"+strings.Replace(f.Name(), ".html", ".htm", -1), source, 0777)
			ioutil.WriteFile(COPY_TO_DIR+"/html/"+f.Name(), source, 0777)
			break
		}
	}else{
		CURRENT_DIR = f.Name()
	}
	return nil
}

func exists(dict M, key string) string {
	if val, ok := dict[key]; ok {
		return fmt.Sprintf("%v", val)
	}
	return ""
}
