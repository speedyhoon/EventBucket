package main

import (
	"fmt"

	"net/url"
	"net/http"
)

//Basic validation to get only the items listed in options, all others are ignored
//func validate(form url.Values, options map[string][]int)interface{}{
func validate(form url.Values, options map[string]InputTypers) interface{} {
	tryThis := map[string]string{}
	for option, min_max := range options {
		if options[option].Html != "submit" && options[option].Html != "form" {
			array, ok := form[option]
			if ok {
				length := len(array[0])
				if length >= min_max.Min && length <= min_max.Max {
					tryThis[option] = array[0]
				}else if min_max.Required {
					fmt.Print("\nlength is not within range")
					return false
				}else {
					fmt.Print("\nELSE length is not within range " + option)
				}
			}else if min_max.Required {
				fmt.Print("\nform[option] not in array")
				return false
			}else {
				fmt.Print("\nELSE form[option] not in array " + option)
			}
		}
	}
	return tryThis
}

func valid8(form url.Values, options map[string]Inputs) interface{} {
	tryThis := map[string]string{}
	for option := range options {
		if options[option].Html != "submit" {
			array, ok := form[option]
			if ok {
				tryThis[option] = array[0]
			}else {
				fmt.Print("\nELSE form[option] not in array " + option)
			}
		}
	}
	return tryThis
}

//research http://net.tutsplus.com/tutorials/client-side-security-best-practices/
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

type Inputs struct {
	Html, Label, Help, Value         string
	Select                           []string
	SelectValues map[string]string
	Checked, MultiSelect             bool
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
	return "\"" + input + "\""
}





func getClubSelectBox(club_list []map[string]interface{}) map[string]string {
	drop_down := make(map[string]string)
	for _, row := range club_list {
		drop_down[fmt.Sprintf("%v", row["_id"])[13:37]] = fmt.Sprintf("%v", row["name"])
	}
	return drop_down
}


func redirectTo(runThisFirst func(http.ResponseWriter, *http.Request), path string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		runThisFirst(w, r)
		http.Redirect(w, r, path, http.StatusSeeOther)//303 mandating the change of request type to GET
	}
}
func redirectPermanent(path string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, path, http.StatusMovedPermanently)//Search engine Optimisation
	}
}
