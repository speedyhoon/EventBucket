package main

import (
	"net/http"
	"strings"
	"fmt"
)

func check_form(options map[string]Inputs, r *http.Request)(map[string]string){
	//TODO add a struct parameter to this function (like Club, Event, etc) so that it can assign the values to the struct rather than return new_values
	//TODO when the form doesn't meet the requirements send user back to certain page and display an error message
	//TODO the base template should be able to handle error messages to it and display them accordingly.
	//TODO given the form elements convert the string to type X and return the instance of the struct
	//TODO add validation for a group of options like <select>
	r.ParseForm()
	form := r.Form
	new_values := make(map[string]string)
	for option := range options {
		if options[option].Html != "submit" {
			array, ok := form[option]
			if ok && ((options[option].Required && array[0] != "") || !options[option].Required) {
				if len(array) > 1{
					new_values[option] = strings.Join(array, ",")
				}else{
					new_values[option] = array[0]
				}
			}else {
				fmt.Printf("options[%v] is REQUIRED OR is not in array", option)
//				warning("options[%v] is REQUIRED OR is not in array", option)
			}
		}
	}
	return new_values
}
