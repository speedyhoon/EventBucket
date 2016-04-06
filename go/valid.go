package main

import (
	"net/http"
	"strings"
)

//Is it worth while to auto add failed forms to session so it doesn't have to be done in each http handler?
func isValid(r *http.Request, fields []field) ([]field, bool) {
	r.ParseForm()
	if len(r.Form) == 0 {
		return fields, false
	}
	//Process the post request as normal if len(r.Form) > len(fields).
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
			var g string
			if len(fieldValue) > 0 {
				g = fieldValue[0]
			}
			t.Println("validation -- ok", ok, "len()", len(fieldValue), "value", g)
		}

		//TODO change fieldValue to a string instead of slice of strings. Almost all fields submit a string instead of an array.
		fieldValue, ok = r.Form[field.name]

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

		if fields[i].Error != "" {
			//Set the first field with failed validation to have focus onscreen
			if valid {
				fields[i].AutoFocus = true
				valid = false
			}
		}
	}
	return fields, valid
}
