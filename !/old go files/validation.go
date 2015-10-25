package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func checkForm(options []Inputs, r *http.Request) map[string]string {
	//TODO add a struct parameter to this function (like Club, Event, etc) so that it can assign the values to the struct rather than return newValues
	//TODO when the form doesn't meet the requirements send user back to certain page and display an error message
	//TODO the base template should be able to handle error messages to it and display them accordingly.
	//TODO given the form elements convert the string to type X and return the instance of the struct
	//TODO add validation for a group of options like <select>
	r.ParseForm()
	form := r.Form
	newValues := make(map[string]string)
	for _, option := range options {
		if option.name != "" {
			array, ok := form[option.name]
			if ok {
				if (option.required && array[0] != "") || !option.required {
					if len(array) > 1 {
						newValues[option.name] = strings.Join(array, ",") //TODO it would be nice to trim whitepspace when joining the items together
					} else {
						newValues[option.name] = strings.TrimSpace(array[0])
						/*//Test if the data length is less than specified
						if option.maxLength > 0 {
							length := len(array[0])
							if length <= option.maxLength && length >= option.minLength {
								newValues[option.name] = strings.TrimSpace(array[0])
							}
						} else {
							newValues[option.name] = strings.TrimSpace(array[0])
						}*/
					}
				} else {
					warning.Printf("options[%v] is REQUIRED OR is not in array", option)
				}
			}
		}
	}
	return newValues
}

/*func v8(inputs []Inputs, r *http.Request) (map[string]string, bool) {
	//TODO possibly call this function directly when a post is made and then pass the function if valid to execute the update/insert
	err := r.ParseForm()
	if err != nil {
		//		http.Error(w, "Invalid request form data", 400)
		return map[string]string{}, false
	}
	form := r.Form
	passed := true
	newValues := make(map[string]string)
	for _, input := range inputs {
		if input.Name != "" {
			dataSlice, ok := form[input.Name]
			if ok && len(dataSlice) == 1 {
				data := strings.TrimSpace(dataSlice[0])
				dataLength := len(data)
				if (!input.Required && dataLength >= 1) || (input.Required && dataLength <= input.MaxLength && dataLength >= input.MinLength) {
					newValues[input.Name] = data
				} else {
					passed = false
					//					dump("Display info that this field didn't have enough characters", input.Name, data)
				}
			} else {
				passed = false
				//				dump("There was too many items in dataSlice", input.Name, dataSlice)
			}
			//			Error.Println("forgot to add a name to input!")
			//			export(input)
		}
	}
	return newValues, passed
}*/

func valid8(options []Inputs, r *http.Request) (M, bool) {
	//TODO add a struct parameter to this function (like Club, Event, etc) so that it can assign the values to the struct rather than return newValues
	//TODO when the form doesn't meet the requirements send user back to certain page and display an error message
	//TODO the base template should be able to handle error messages to it and display them accordingly.
	//TODO given the form elements convert the string to type X and return the instance of the struct
	//TODO add validation for a group of options like <select>
	r.ParseForm()
	form := r.Form
	newValues := make(M)
	var passedV8tion []bool
	var formArray []string
	var value string
	var ok bool
	var err error
	var valueInt int
	for _, option := range options {
		if option.name != "" {
			passedV8tion = append(passedV8tion, false)
			formArray, ok = form[option.name]
			if ok && (!option.required || (option.required && len(formArray) > 0)) {
				if len(formArray) > 1 {
					//					newValues[option.name] = strings.Join(array, ",")
				} else {
					value = strings.TrimSpace(formArray[0])
				}
				switch option.varType {
				case "int":
					valueInt, err = strconv.Atoi(value)
					if err == nil {
						newValues[option.name] = valueInt
						passedV8tion = append(passedV8tion, true)
					} else {
						passedV8tion = append(passedV8tion, false)
					}
				case "string":
					valueInt = len(value)
					if valueInt >= option.minLength && valueInt <= option.minLength {
						newValues[option.name] = value
						passedV8tion = append(passedV8tion, true)
					} else {
						passedV8tion = append(passedV8tion, false)
						//						dump("string failed")
					}
				}
			} else {
				warning.Printf("options[%v] is REQUIRED OR is not in array", option)
			}
		}
	}

	if !testAllTrue(passedV8tion) {
		//TODO output all these error messages to screen at once. A form might have several invalid fields at the same time
		info.Println("validation was not good")
		return M{}, false
	}

	var event Event
	var tempInt int
	for name, value2 := range newValues {
		switch name {
		case "eventid":
			event, err = getEvent(fmt.Sprintf("%v", value2))
			if err == nil {
				newValues["event"] = event
			} else {
				warning.Printf("event with id '%v' doesn't exist", value2)
				return M{}, false
			}
		case "rangeid":
			tempInt = value2.(int)
			if tempInt < 0 || tempInt > len(event.Ranges) {
				warning.Printf("event with range id '%v' doesn't exist", value2)
				return M{}, false
			}
		case "shooterid":
			//TODO this might be better as a pointer to check that index is not null
			tempInt = value2.(int)
			if tempInt < 0 || tempInt > len(event.Shooters) {
				warning.Printf("event with shooter id '%v' doesn't exist", value2)
				return M{}, false
			}
		}
	}
	return newValues, true
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
