package main

import (
	"fmt"
	"net/http"
)

//Is it worth while to auto add failed forms to session so it doesn't have to be done in each http handler?
func isValid(r *http.Request, fields []field) ([]field, bool) {
	r.ParseForm()
	if len(r.Form) == 0 {
		return fields, false
	}
	//Process the post request as normal if len(r.Form) > len(fields)
	var fieldValue []string
	var ok bool
	valid := true
	for i, field := range fields {
		//TODO remove developer message
		if field.v8 == nil {
			field.Error = "No v8 function setup!"
			warn.Println("No v8 function setup! for", field.name)
			continue
		}

		//TODO change fieldValue to a string instead of slice of strings. Almost all fields submit a string instead of an array.
		fieldValue, ok = r.Form[field.name]

		//trace.Println(field.name, "ok=", ok, "req=", field.Required, "value=", fieldValue)
		/*	if !ok {
						trace.Println("!ok", field.name)

				//		if field.Required && field.defValue != nil {
					//		fieldValue = field.defValue()
						//	trace.Println("err0")
							//			} else if field.Required {
							//				valid = false
							//				fields[i].Error = "Missing value"
							//				trace.Println("err1")
							//				continue //Skip to the next loop iteration.
					//	} else {
							fieldValue = []string{""}
			//				trace.Println("err2")
							//				continue //Skip to the next loop iteration.
						//}
					}*/

		//trace.Println("check", field.defValue != nil, !ok, fieldValue == nil, len(fieldValue))
		/*if len(fieldValue) > 0 {
			trace.Println(field.name, "contents: ", fieldValue[0], fieldValue[0] == "")
		} else {
			trace.Println(field.name, "no contents: ", fieldValue)
		}*/
		if field.defValue != nil && (!ok || fieldValue == nil || len(fieldValue) == 0 || (len(fieldValue) == 1 && fieldValue[0] == "")) {
			//			trace.Println("set default value")
			fieldValue = field.defValue()
		} else if !ok {
			//			trace.Println("set empty string")
			fieldValue = []string{""}
		}

		fields[i].internalValue, fields[i].Error = field.v8(field, fieldValue...)
		if fields[i].Error != "" {
			//Set the first field with failed validation to have focus onscreen
			if valid {
				fields[i].AutoFocus = true
			}
			valid = false
		}
		fields[i].Value = fmt.Sprintf("%v", fields[i].internalValue)
		switch fields[i].internalValue.(type) {
		case bool:
			fields[i].Checked = fields[i].internalValue.(bool)
		}
	}
	return fields, valid
}
