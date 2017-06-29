package main

import (
	"net/http"
	"net/url"
	"strings"
)

func validBoth(r *http.Request, formID uint8) (form, bool) {
	var err error
	var u *url.URL

	if r.Method == get {
		u, err = url.Parse(r.RequestURI)
	} else {
		err = r.ParseForm()
	}

	if err != nil {
		warn.Println(err)
		return getForm(formID), false
	}

	var values url.Values
	if r.Method == get {
		values = u.Query()
	} else {
		values = r.Form
	}

	return isValid(values, formID)
}

//Is it worth while to auto add failed forms to session so it doesn't have to be done in each http handler?
func isValid(urlValues url.Values, formID uint8) (form, bool) {
	f := getForm(formID)
	if len(urlValues) == 0 {
		return f, false
	}
	//Process the post request as normal if len(urlValues) >= len(fields).
	var fieldValue []string
	var ok bool
	valid := true
	for i, field := range f.Fields {
		/*// Output warning if validation function is not set for this field in the submitted form.
		if debug && field.v8 == nil {
			field.Error = "No v8 function setup!"
			warn.Println("No v8 function setup! for", field.name)
			continue
		}*/
		fieldValue, ok = urlValues[field.name]

		//if fieldValue is empty and field is required
		if !ok || len(fieldValue) == 0 || len(fieldValue) == 1 && strings.TrimSpace(fieldValue[0]) == "" {
			if field.Required {
				f.Fields[i].Error = "Please fill in this field."
			}
			//else if field is not required & its contents is empty - don't validate
		} else {
			//Otherwise validate user input
			field.v8(&f.Fields[i], fieldValue...)
		}

		if f.Fields[i].Error != "" {
			//Set the first field with failed validation to have focus onscreen
			if valid {
				f.Fields[i].AutoFocus = true
				valid = false
			}
		}
	}
	return f, valid
}
