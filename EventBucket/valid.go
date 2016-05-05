package main

import (
	"net/http"
	"net/url"
	"strings"
)

func validPost(r *http.Request, fields []field) ([]field, bool) {
	err := r.ParseForm()
	if err != nil {
		warn.Println(err)
		return fields, false
	}
	return isValid(r.Form, fields)
}

func validGet(r *http.Request, fields []field) ([]field, bool) {
	u, err := url.Parse(r.RequestURI)
	if err != nil {
		warn.Println(err)
		return fields, false
	}
	return isValid(u.Query(), fields)
}

//Is it worth while to auto add failed forms to session so it doesn't have to be done in each http handler?
func isValid(urlValues url.Values, fields []field) ([]field, bool) {
	if len(urlValues) == 0 {
		return fields, false
	}
	//Process the post request as normal if len(urlValues) >= len(fields).
	var fieldValue []string
	var ok bool
	valid := true
	for i, field := range fields {
		//TODO remove developer message
		if debug {
			if field.v8 == nil {
				field.Error = "No v8 function setup!"
				warn.Println("No v8 function setup! for", field.name)
				continue
			}
			if field.manyRequiredQty > 0 || len(field.manyRequired) > 0 {
				t.Println("manyRequiredQty:", field.manyRequiredQty, "field.manyRequired:", field.manyRequired)
			}
		}

		fieldValue, ok = urlValues[field.name]
		t.Println(fieldValue, ok)

		//if fieldValue is empty and...
		if !ok || len(fieldValue) == 0 || (len(fieldValue) == 1 && strings.TrimSpace(fieldValue[0]) == "") {
			if field.Required {
				fields[i].Error = "Please fill in this field"
			}
			//else if field is not required - do nothing.

			//If field is empty and it has a default value function assign the default value.
			if field.defValue != nil {
				t.Println("set default value")
				fieldValue = field.defValue()
			}
			//else if field doesn't have a default value - do nothing.
		} else {
			//Otherwise validate user input
			field.v8(&fields[i], fieldValue...)
		}

		//TODO write comments!
		//TODO this could be simplified - manyRequiredQty isn't used here.
		/*if fields[i].Error == "" && field.manyRequiredQty >= 1 {
			if debug && len(field.manyRequired) < 2 {
				//manyRequired: []int{0, 1, 2}, manyRequiredQty: 1,
				warn.Println("field %v manyRequired slice doesn't have enough indexes. Need at least 2 items %v", field.name, field.manyRequired)
			}
			var pass, errors bool
			for _, index := range field.manyRequired {
				if fields[index].Error != "" {
					errors = true
					pass = false
					break
				}
				if fields[index].Value != "" && fields[index].Error == "" {
					pass = true
				}
			}

			//If it hasn't passed and there were no errors
			if !pass && !errors {
				//Assign the error message to the first field in the list
				fields[field.manyRequired[0]].Error = "Please fill out one of these fields"
				valid = false
			}
		}*/

		if fields[i].Error != "" {
			warn.Println(fields[i].name, fields[i].Error)
			//Set the first field with failed validation to have focus onscreen
			if valid {
				fields[i].AutoFocus = true
				valid = false
			}
		}
	}
	return fields, valid
}
