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
	num, err := strconv.Atoi(strNum)
	if err != nil {
		return num, err.Error()
	}
	if num >= field.min && num <= field.max && (field.step == 0 || field.step != 0 && num%field.step == 0) || !field.Required && num == 0 {
		return num, ""
	}
	//		return num, errors.New("field integer doesn't pass validation")
	return num, "field integer doesn't pass validation"
}
func isValidStr(str string, field field) (interface{}, string) {
	info.Println(field.name, "r=", field.Required)
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
	info.Println("__________________________________value='" + str + "'")
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
	//	var err error
	valid := true
	for i, field := range fields {
		//TODO remove developer message
		if field.v8 == nil {
			field.Error = "No v8 function setup!"
			warn.Println("No v8 function setup!")
			continue
		}

		fieldValue, ok = r.Form[field.name]
		trace.Println(field.name, "ok=", ok, "req=", field.Required)
		if !ok {
			trace.Println("!ok", field.name)

			if field.Required && field.defValue != nil {
				fieldValue = field.defValue()
				trace.Println("err0")
				//			} else if field.Required {
				//				valid = false
				//				fields[i].Error = "Missing value"
				//				trace.Println("err1")
				//				continue //Skip to the next loop iteration.
			} else {
				fieldValue = []string{""}
				trace.Println("err2")
				//				continue //Skip to the next loop iteration.
			}
		}

		fields[i].internalValue, fields[i].Error = field.v8(strings.TrimSpace(fieldValue[0]), field)
		info.Println("\n\n", field.name, "\nval=", fieldValue[0], "\nErr=", fields[i].Error)
		if fields[i].Error != "" {
			warn.Println("err2")
			valid = false
			//		}else{
			//			temp := field.kind.(type)
			//			fields[i].internalValue = fields[i].Vaue.(temp)
			//		} else {
			//			switch field.kind.(type) {
			//			case string:
			//
			//			}
		}
		fields[i].Value = fmt.Sprintf("%v", fields[i].internalValue)
		/*switch field.kind.(type) {
		case bool:
			fmt.Printf("boolean %t\n", field.kind)
		case int:
			//			fields[i].Value, err = isValidInt(strings.TrimSpace(fieldValue[0]), field)
			if err == nil {
				//				fields[i].isValid = true
			} else {
				valid = false
				fields[i].Error = "integer supplied was wrong"
				fmt.Println(field.name + "integer supplied was wrong")
			}
		case string:
			if len(fieldValue) > 0 {
				fields[i].Value, err = isValidStr(strings.TrimSpace(fieldValue[0]), field)
			} else {
				fields[i].Value, err = isValidStr("", field)
			}
			//			fields[i].isValid = fields[i].Error == ""
			if err != nil {
				//				fields[i].Error = fmt.Sprintf("%v", err)
				fields[i].Error = err.Error()
				valid = false
			}
			//			if fields[i].Error != "" {
			//				valid = false
			//				fields[i].Error = "string doesn't match"
			//				fmt.Println(field.name + "string doesn't match")
			//			}
		case []string:
			fmt.Println("within string slice")
			for key, thingy := range fieldValue {
				fmt.Println(key, thingy)
			}
			//			fields[i].Value = fieldValue
		default:
			warn.Printf("unexpected type %T", field.kind) // %T prints whatever type t is
		}*/
	}
	return fields, valid
}
