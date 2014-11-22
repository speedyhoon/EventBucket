package main

import (
	"net/http"
)

func check_form(options map[string]Inputs, r *http.Request)(map[string]string){
	/*TODO add a struct parameter to this function (like Club, Event, etc)
	so that it can assign the values to the struct rather than return new_values*/
	//TODO when the form doesn't meet the requirements send user back to certain page and display an error message
	//TODO the base template should be able to handle error messages to it and display them accordingly.
	//TODO given the form elements convert the string to type X and return the instance of the struct
	r.ParseForm()
	form := r.Form
	new_values := make(map[string]string)
	for option := range options {
		if options[option].Html != "submit" {
			array, ok := form[option]
			if ok && ((options[option].Required && array[0] != "") || !options[option].Required) {
				new_values[option] = array[0]
			}else {
				log("options[%v] is REQUIRED OR is not in array", option)
			}
		}
	}
	return new_values
}
