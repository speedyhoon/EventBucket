package main

import (
	"fmt"
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
			trace.Println("validation -- ok", ok, "len()", len(fieldValue), "value", g)
		}

		//TODO change fieldValue to a string instead of slice of strings. Almost all fields submit a string instead of an array.
		fieldValue, ok = r.Form[field.name]

		//if fieldValue is empty
		if !ok || len(fieldValue) == 0 || (len(fieldValue) == 1 && strings.TrimSpace(fieldValue[0]) == "") {
			if field.Required {
				fields[i].Error = "Value is required"
			}
			//If field is empty and it has a default value function assign the default value.
			if field.defValue != nil {
				trace.Println("set default value")
				fieldValue = field.defValue()
			}
		} else {
			//Otherwise validate user input
			fields[i].internalValue, fields[i].Error = field.v8(field, fieldValue...)
		}

		if fields[i].Error != "" {
			//Set the first field with failed validation to have focus onscreen
			if valid {
				fields[i].AutoFocus = true
				valid = false
			}
		} else {
			fields[i].Value = fmt.Sprintf("%v", fields[i].internalValue)
			switch fields[i].internalValue.(type) {
			case bool:
				fields[i].Checked = fields[i].internalValue.(bool)
			}
		}
	}
	return fields, valid
}
