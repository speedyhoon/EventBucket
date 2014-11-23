package main

import (
	"net/http"
	"strings"
	"fmt"
	"strconv"
)

func check_form(options []Inputs, r *http.Request)(map[string]string){
	//TODO add a struct parameter to this function (like Club, Event, etc) so that it can assign the values to the struct rather than return new_values
	//TODO when the form doesn't meet the requirements send user back to certain page and display an error message
	//TODO the base template should be able to handle error messages to it and display them accordingly.
	//TODO given the form elements convert the string to type X and return the instance of the struct
	//TODO add validation for a group of options like <select>
	r.ParseForm()
	form := r.Form
	new_values := make(map[string]string)
//	for option := range options {
	for _, option := range options {
//		if options[option].Html != "submit" {
		if option.Html != "submit" {
//			array, ok := form[option]
			array, ok := form[option.Name]
//			if ok && ((options[option].Required && array[0] != "") || !options[option].Required) {
			if ok && ((option.Required && array[0] != "") || !option.Required) {
				if len(array) > 1{
//					new_values[option] = strings.Join(array, ",")
					new_values[option.Name] = strings.Join(array, ",")
				}else{
//					new_values[option] = strings.TrimSpace(array[0])
					new_values[option.Name] = strings.TrimSpace(array[0])
				}
			}else {
				fmt.Printf("options[%v] is REQUIRED OR is not in array", option)
//				warning("options[%v] is REQUIRED OR is not in array", option)
			}
		}
	}
	return new_values
}

func valid8(options []Inputs, r *http.Request)(M,bool){
	//TODO add a struct parameter to this function (like Club, Event, etc) so that it can assign the values to the struct rather than return new_values
	//TODO when the form doesn't meet the requirements send user back to certain page and display an error message
	//TODO the base template should be able to handle error messages to it and display them accordingly.
	//TODO given the form elements convert the string to type X and return the instance of the struct
	//TODO add validation for a group of options like <select>
	r.ParseForm()
	form := r.Form
	new_values := make(M)
	var passedV8tion []bool
//	passedValidation := 0	//not processed = 0, passed = 1, failed = 2
	var formArray []string
	var value string
	var ok bool
	var err error
	var valueInt int
	for index, option := range options {
		if option.Html != "submit" {
			passedV8tion[index] = false
			formArray, ok = form[option.Name]
//			if ok && ((option.Required && array[0] != "") || !option.Required) {
			if ok && (!option.Required || (option.Required && len(formArray) > 0)) {
				if len(formArray) > 1{
//					new_values[option.Name] = strings.Join(array, ",")
				}else{
//					new_values[option.Name] = strings.TrimSpace(array[0])
					value = strings.TrimSpace(formArray[0])
				}

				switch option.VarType{
				case "int":
					valueInt, err = strconv.Atoi(value)
					if err == nil{
						new_values[option.Name] = valueInt
						passedV8tion[index] = true
//						passedValidation = true
//					}else{
//						passedValidation = 2
					}
				case "string":
					valueInt = len(value)
					if valueInt >= option.VarMinLen && valueInt <= option.VarMinLen {
						new_values[option.Name] = value
						passedV8tion[index] = true
//						passedValidation = 1
//					}else{
//						passedValidation = 2
					}
				}
			}else{
				//TODO output this error message to screen
				fmt.Printf("options[%v] is REQUIRED OR is not in array", option)
//				passedValidation = 2
			}
		}
	}

	if !testAllTrue(passedV8tion) {
		return make(M), false
	}else{
		var event Event
		for name, value := range new_values {
			switch name{
			case "event_id":
				if event, ok = getEvent(fmt.Sprintf("%v", value)); ok {
					new_values["event"] = event
				}else {
					fmt.Printf("event with id '%v' doesn't exist", value)
					return make(M), false
				}
			case "range_id":
				if integer, _ := value.(int); integer >= len(event.Ranges) {
					fmt.Printf("event with range id '%v' doesn't exist", value)
					return make(M), false
				}
				/*case "shooter_id":
				//TODO how to check EventShooter is not empty?
				if _, ok := event.Shooters[value.(int)]; !ok{
					fmt.Printf("event with shooter id '%v' doesn't exist", value)
					return make(M), false
				}*/
			}
		}
	}
	return new_values, true
}

func testAllTrue(array []bool)bool{
	if len(array) > 0{
		for _, test := range array{
			if test == false{
				return false
			}
		}
	}else{
		return false
	}
	return true
}
