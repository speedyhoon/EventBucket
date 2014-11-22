package main

import (
	"fmt"
	"strings"
	"net/url"
	"net/http"
	"strconv"
)

//Basic validation to get only the items listed in options, all others are ignored
//func validate(form url.Values, options map[string][]int)interface{}{
func validate(form url.Values, options map[string]InputTypers) interface{} {
	tryThis := map[string]string{}
//	for option, min_max := range options {
	for option := range options {
		if options[option].Html != "submit" && options[option].Html != "form" {
			array, ok := form[option]
			if ok {
//				length := len(array[0])
//				if length >= min_max.Min && length <= min_max.Max {
					tryThis[option] = array[0]
//				}else if min_max.Required {
//					fmt.Print("\nlength is not within range")
//					return false
//				}else {
//					fmt.Print("\nELSE length is not within range " + option)
//				}
//			}else if min_max.Required {
//				fmt.Print("\nform[option] not in array")
//				return false
			}else {
				fmt.Print("\nELSE form[option] not in array " + option)
			}
		}
	}
	return tryThis
}

func valid8(form url.Values, options map[string]Inputs) map[string]interface{} {
	tryThis := make(map[string]interface{})
	for option := range options {
		if options[option].Html != "submit" {
			array, ok := form[option]
			if ok {
				tryThis[option] = array[0]
			}else {
				fmt.Print(fmt.Sprintf("\nELSE options[%v] not in array ", option))
			}
		}
	}
	return tryThis
}

func validInsert(form url.Values, options map[string]Inputs) map[string]interface{} {
	tryThis := make(map[string]interface{})
	for option := range options {
		if options[option].Html != "submit" {
			array, ok := form[option]
			if ok {
				if (options[option].Required && array[0] != "") || !options[option].Required{
					if schema(option) != ""{
						tryThis[schema(option)] = array[0]
					}else {
						if option == "rangeType"{
//							tryThis["a"] = make([]int)
						}else{
							tryThis[option] = array[0]
						}
//						tryThis[option] = array[0]
					}
				}else{
					fmt.Print(fmt.Sprintf("\nELSE options[%v] is REQUIRED", option))
				}
			}else {
				fmt.Print(fmt.Sprintf("\nELSE options[%v] not in array ", option))
			}
		}
	}
	return tryThis
}

func exists(dict map[string]interface{}, key string)string{
	if val, ok := dict[key]; ok {
		return fmt.Sprintf("%v", val)
	}
	return ""
}

//research http://net.tutsplus.com/tutorials/client-side-security-best-practices/
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

type Inputs struct {
	Html, Label, Help, Value			string
	Placeholder								string
	Select									[]string
	SelectValues 							map[string]string
	Checked, MultiSelect, Required	bool
	Min, Max									int
}

type InputTypers struct {
	Required                       bool
	Min, Max, RangeMin, RangeMax   int
	Name, Type                     string
	Html                           string
	Select                         []string
	SelectValues map[string]string
	Label, Help, PlaceHolder       string
	AutoCorrect, AutoCapitalize    bool
	Value                          string
	Disabled, Checked, MultiSelect bool

	method, action string
	table          bool
}

func addQuotes(input string) string {
	if strings.Contains(input, " ") {
		return "\"" + input + "\""
	}
	return input
}

func vardump(input map[string]interface{}) {
	for index, row := range input {
		log("\nindex: %v, ROW: %v", index, row)
		//		fmt.Print(fmt.Sprintf("\nrow: %v", row))
		//		for col, item := range row{
		//			fmt.Print(fmt.Sprintf("\ncol: %v", col))
		//			fmt.Print(fmt.Sprintf("\nitem: %v", item))
		//		}
		//		fmt.Print("\n\n")
	}
}

func getClubSelectBox(club_list []map[string]interface{}) map[string]string {
	drop_down := make(map[string]string)
	id := ""
	name := ""
	for _, row := range club_list {
		name = exists(row, schemaNAME)
		id = exists(row, "_id")
		if name != "" && id != "" {
			drop_down[id] = name
		}
	}
	return drop_down
}

func redirectVia(runThisFirst func(http.ResponseWriter, *http.Request), path string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		runThisFirst(w, r)
		http.Redirect(w, r, path, http.StatusSeeOther)//303 mandating the change of request type to GET
	}
}
func redirectTo(path string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, path, http.StatusSeeOther)
	}
}
func redirecter(path string, w http.ResponseWriter, r *http.Request){
//	return func(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, path, http.StatusSeeOther)
//	}
}
func redirectPermanent(path string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, path, http.StatusMovedPermanently)//Search engine Optimisation
	}
}
func log(format string, a ...interface{}){
	fmt.Printf("\n"+format, a...)
}
func dump(input interface{}){
	fmt.Printf("\n%v", input)
}
func echo(input interface{})string{
	return fmt.Sprintf("%v", input)
}
func str_to_int(input string)int{
	output, err := strconv.Atoi(input)
	checkErr(err)
	return output
}
