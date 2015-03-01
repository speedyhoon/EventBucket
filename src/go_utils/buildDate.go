package main

import (
	"bytes"
	"io/ioutil"
	"time"
	"fmt"
	"path/filepath"
	"os"
)
//loop through all files & sub directories and replace the search variables with the actual values

const (
	ROOT_DIR = "../eventbucketM/"
	COPY_TO_DIR = ROOT_DIR +"!/"

	BD_ARGS = `--dbpath", databasePath, "--port", "38888", "--nssize", "1", "--smallfiles", "--noscripting", `
)

var (
	CURRENT_DIR = ""
	DEV = "true"
	ReplaceChars = map[string]interface{}{
		"VersionNumber": 58,		//TODO get the Git tag from the last commit
		"DIR_ROOT": "",
		"BuildDate": time.Now().Format("January 2, 2006"),
		"schemaSHOOTER": "S",
		"schemaAutoInc": "U",
		"schemaRANGE": "R",
		"schemaSORT": "o",
		"schemaGRADES": "g",
	}
	DevMode = map[string]interface{}{
		"DbArgs": BD_ARGS + `"--noauth", "--slowms", "3", "--cpu", "--profile", "2", "--objcheck", "--notablescan", "--rest`,
	}
	ProdMode = map[string]interface{}{
		"DbArgs": BD_ARGS + `"--nohttpinterface`,
	}
)
//TODO eventually move these settings to a json or yaml file.

func main(){
	joinSettings()
	filepath.Walk(ROOT_DIR + "golang", walkPath)
//	filepath.Walk(ROOT_DIR + "sass", walkPath)
//	filepath.Walk(ROOT_DIR + "js", walkPath)
//	filepath.Walk(ROOT_DIR + "html", walkPath)
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
		if CURRENT_DIR == "golang"{
			ioutil.WriteFile(COPY_TO_DIR+f.Name(), source, 0777)
		}
	}else{
		CURRENT_DIR = f.Name()
	}
	return nil
}
