package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func checkForm(options []Inputs, r *http.Request) map[string]string {
	//TODO add a struct parameter to this function (like Club, Event, etc) so that it can assign the values to the struct rather than return new_values
	//TODO when the form doesn't meet the requirements send user back to certain page and display an error message
	//TODO the base template should be able to handle error messages to it and display them accordingly.
	//TODO given the form elements convert the string to type X and return the instance of the struct
	//TODO add validation for a group of options like <select>
	r.ParseForm()
	form := r.Form
	new_values := make(map[string]string)
	for _, option := range options {
		if option.Name != "" {
			array, ok := form[option.Name]
			if ok {
				if (option.Required && array[0] != "") || !option.Required {
					if len(array) > 1 {
						new_values[option.Name] = strings.Join(array, ",") //TODO it would be nice to trim whitepspace when joining the items together
					} else {
						new_values[option.Name] = strings.TrimSpace(array[0])
					}
				} else {
					Warning.Printf("options[%v] is REQUIRED OR is not in array", option)
				}
			}
		}
	}
	return new_values
}

func valid8(options []Inputs, r *http.Request) (M, bool) {
	//TODO add a struct parameter to this function (like Club, Event, etc) so that it can assign the values to the struct rather than return new_values
	//TODO when the form doesn't meet the requirements send user back to certain page and display an error message
	//TODO the base template should be able to handle error messages to it and display them accordingly.
	//TODO given the form elements convert the string to type X and return the instance of the struct
	//TODO add validation for a group of options like <select>
	r.ParseForm()
	form := r.Form
	new_values := make(M)
	var passedV8tion []bool
	var formArray []string
	var value string
	var ok bool
	var err error
	var valueInt int
	for _, option := range options {
		if option.Name != "" {
			passedV8tion = append(passedV8tion, false)
			formArray, ok = form[option.Name]
			if ok && (!option.Required || (option.Required && len(formArray) > 0)) {
				if len(formArray) > 1 {
					//					new_values[option.Name] = strings.Join(array, ",")
				} else {
					value = strings.TrimSpace(formArray[0])
				}
				switch option.VarType {
				case "int":
					valueInt, err = strconv.Atoi(value)
					if err == nil {
						new_values[option.Name] = valueInt
						passedV8tion = append(passedV8tion, true)
					} else {
						passedV8tion = append(passedV8tion, false)
					}
				case "string":
					valueInt = len(value)
					if valueInt >= option.VarMinLen && valueInt <= option.VarMinLen {
						new_values[option.Name] = value
						passedV8tion = append(passedV8tion, true)
					} else {
						passedV8tion = append(passedV8tion, false)
						//						dump("string failed")
					}
				}
			} else {
				Warning.Printf("options[%v] is REQUIRED OR is not in array", option)
			}
		}
	}

	if !testAllTrue(passedV8tion) {
		//TODO output all these error messages to screen at once. A form might have several invalid fields at the same time
		Info.Println("validation was not good")
		return make(M), false
	}

	var event Event
	var tempInt int
	for name, value := range new_values {
		switch name {
		case "eventid":
			event, err = getEvent(fmt.Sprintf("%v", value))
			if err == nil {
				new_values["event"] = event
			} else {
				Warning.Printf("event with id '%v' doesn't exist", value)
				return make(M), false
			}
		case "rangeid":
			tempInt = value.(int)
			if tempInt < 0 || tempInt > len(event.Ranges) {
				Warning.Printf("event with range id '%v' doesn't exist", value)
				return make(M), false
			}
		case "shooterid":
			//TODO this might be better as a pointer to check that index is not null
			tempInt = value.(int)
			if tempInt < 0 || tempInt > len(event.Shooters) {
				Warning.Printf("event with shooter id '%v' doesn't exist", value)
				return make(M), false
			}
		}
	}
	return new_values, true
}

func testAllTrue(array []bool) bool {
	if len(array) > 0 {
		for _, test := range array {
			if test == false {
				return false
			}
		}
	} else {
		return false
	}
	return true
}
