package main

import (
	"fmt"
	//	"os"
	"bytes"
	"io/ioutil"
	//	"log"
)

/*var (
	Info		*log.Logger = log.New(os.Stdout, "INFO: ",    log.Ldate|log.Ltime|log.Lshortfile)
	Warning	*log.Logger = log.New(os.Stderr, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
)

func info(format string, a ...interface{}){
	if !PRODUCTION {
		Info.Printf(format, a...)
	}
}
func warning(format string, a ...interface{}){
	if !PRODUCTION {
		Warning.Printf(format, a...)
	}
}*/

func dev_mode_loadHTM(page_name string, existing_minified_file []byte) []byte {
	bytes, err := ioutil.ReadFile(fmt.Sprintf(PATH_HTML_SOURCE, page_name))
	checkErr(err)
	bytes = dev_mode_minifyHtml(page_name, bytes)
	existing_len := len(existing_minified_file)
	new_len := len(bytes)
	if existing_len != new_len {
		ioutil.WriteFile(fmt.Sprintf(PATH_HTML_MINIFIED, page_name), bytes, 0777)
		fmt.Printf("Page '%v' had %v bytes removed (%v percent), total: %v, from: %v", page_name, new_len-existing_len, (existing_len*100/new_len-100)*-1, existing_len, new_len)
		return bytes
	}
	return bytes
	//	return existing_minified_file
}

func dev_mode_check_form(check bool, message string){
	if !check{
		dump(message)
	}
}

func dev_mode_minifyHtml(page_name string, html []byte) []byte {
	minify := html

	if bytes.Contains(minify, []byte("ZgotmplZ")) {
		fmt.Println("Template generation error: ZgotmplZ")
		return []byte("")
	}

	remove_chars := map[string]string{
		"	": "", //Tab
		"\n": "", //new line
		"\r": "", //carriage return
	}
	//TODO remove spaces between block elements like: </div> <div> but keep between inline elements like </span> <span>
	//TODO use improved regex for better searching & replacement
	replace_chars := map[string]string{
		"  ":            " ", //double spaces
		"type=text":     "",
		"type=\"text\"": "",
		"type='text'":   "",
		" >":            ">",
		" <":            "<",
		"< ":            "<",
		">  <":          "> <",
		" />":           "/>",
		"/ >":           "/>",
		"<br/>":         "<br>",
		"</br>":         "<br>",
		"<br />":        "<br>",
	}
	for search, replace := range remove_chars {
		minify = bytes.Replace(minify, []byte(search), []byte(replace), -1)
	}

	backup := minify
	for !bytes.Equal(minify, backup) {
		for search, replace := range replace_chars {
			length := len(minify)
			minify = bytes.Replace(minify, []byte(search), []byte(replace), -1)
			if length != len(minify) {
				fmt.Printf("A dodgy character (%v) was found in the source! Please replace with (%v).", search, replace)
				//				warning("A dodgy character (%v) was found in the source! Please replace with (%v).", search, replace)
			}
		}
	}
	//TODO why is the string not being replaced here, even though it is 100% running?
	minify = bytes.Replace(minify, []byte("~~~"), []byte(" "), -1)
	return minify
}
