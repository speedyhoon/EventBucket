package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	fieldMaxLen = 64
)

func isValidInt(strNum string, field field) (interface{}, string) {
	num := 0
	var err error
	if strNum != "" {
		num, err = strconv.Atoi(strNum)
	} else if field.Required {
		return num, "Please fill in this field"
	}
	if err != nil {
		return num, err.Error()
	}
	if num >= field.min && num <= field.max && (field.step == 0 || field.step != 0 && num%field.step == 0) || !field.Required && num == 0 {
		return num, ""
	}
	return num, "field integer doesn't pass validation"
}
func isValidStr(str string, field field) (interface{}, string) {
	//	info.Println(field.name, "r=", field.Required)
	//	info.Println(field.minLen, field.maxLen, field.Required, `"`+str+`"`)
	length := len(str)
	if field.Required && length == 0 {
		return str, "Please fill in this field"
	}

	if !field.Required && length == 0 {
		return "", ""
	}

	if field.maxLen == 0 {
		field.maxLen = fieldMaxLen
	}

	if length >= field.minLen && length <= field.maxLen {
		return str, ""
	}

	if length < field.minLen || length > field.maxLen {
		plural := "s"
		if length == 1 {
			plural = ""
		}
		return str, fmt.Sprintf("Please change this text be between %v & %v characters long (you are currently using %v character%v).", field.minLen, field.maxLen, length, plural)
	}
	//if field.Required {
	return str, "Please fill in this field"
	//}
}

func isValidBool(str string, field field) (interface{}, string) {
	//	info.Println("__________________________________value='" + str + "'")
	checked := len(str) >= 1
	if field.Required && !checked {
		return false, "Please check this field"
	}
	return checked, ""
}

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
			warn.Println("No v8 function setup!")
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
		if len(fieldValue) > 0 {
			//trace.Println("contents: ", fieldValue[0], fieldValue[0] == "")
		}
		if field.defValue != nil && (!ok || fieldValue == nil || len(fieldValue) == 0 || (len(fieldValue) == 1 && fieldValue[0] == "")) {
			//			trace.Println("set default value")
			fieldValue = field.defValue()
		} else if !ok {
			//			trace.Println("set empty string")
			fieldValue = []string{""}
		}

		fields[i].internalValue, fields[i].Error = field.v8(strings.TrimSpace(fieldValue[0]), field)
		//		info.Println("\n\n", field.name, "\nval=", fieldValue[0], "\nErr=", fields[i].Error)
		if fields[i].Error != "" {
			warn.Println("err2", fields[i].Error, field.name)
			valid = false
			//		}else{
			//			temp := field.kind.(type)
			//			fields[i].internalValue = fields[i].Vaue.(temp)
			//		} else {
			//			switch field.kind.(type) {
			//			case string:
		}
		fields[i].Value = fmt.Sprintf("%v", fields[i].internalValue)
		switch fields[i].internalValue.(type) {
		//		case string, int:
		//			fields[i].Value = fmt.Sprintf("%v", fields[i].internalValue)
		case bool:
			if !fields[i].internalValue.(bool) {
				fields[i].Value = ""
			}
		}
		//		trace.Println("mmm", fields[i].internalValue, fields[i].internalValue, "v=", fields[i].Value)
	}
	return fields, valid
}
