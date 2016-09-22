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
		fieldValue, ok = urlValues[field.name]

		//if fieldValue is empty and field is required
		if !ok || len(fieldValue) == 0 || (len(fieldValue) == 1 && strings.TrimSpace(fieldValue[0]) == "") {
			if field.Required {
				fields[i].Error = "Please fill in this field."
			}
			//else if field is not required - do nothing.
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

//Converts numeric string to uint
func stoU(id string) (uint, error) {
	u, err := strconv.ParseUint(strings.TrimSpace(id), 10, systemArch)
	if err != nil {
		warn.Printf("Unable to convert %v to uint", id)
	}
	return uint(u), err
}
