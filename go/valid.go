package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func isValidInt(strNum string, field field) (int, error) {
	num, err := strconv.Atoi(strNum)
	if err != nil {
		return num, err
	}
	if num >= int(field.min) && num <= int(field.max) && (field.step == 0 || field.step != 0 && num%int(field.step) == 0) || !field.required && num == 0 {
		return num, nil
	}
	return num, errors.New("field integer doesn't pass validation")
}
func isValidStr(str string, field field) (string, error) {
	length := len(str)
	if length >= field.minLen && length <= field.maxLen || !field.required && str == "" {
		return str, nil
	}
	return str, errors.New("field string doesn't pass validation")
}

func isValid(fields []field, r *http.Request) ([]field, bool) {
	r.ParseForm()
	if len(r.Form) < 1 {
		return fields, false
	}
	//Process the post request as normal if len(r.Form) > len(fields)
	var fieldValue []string
	var ok bool
	var err error
	valid := true
	for i, field := range fields {
		fieldValue, ok = r.Form[field.name]
		if !ok {
			valid = false
			continue //Skip to the next loop iteration.
		}

		switch field.kind.(type) {
		case bool:
			fmt.Printf("boolean %t\n", field.kind)
		case int:
			//			fields[i].value, err = isValidInt(strings.TrimSpace(fieldValue[0]), field)
			if err == nil {
				//				fields[i].isValid = true
			} else {
				valid = false
				fields[i].error = "integer supplied was wrong"
				fmt.Println(field.name + "integer supplied was wrong")
			}
		case string:
			//			fields[i].value, fields[i].error = isValidStr(strings.TrimSpace(fieldValue[0]), field)
			//			fields[i].isValid = fields[i].error == ""
			fields[i].error = fmt.Sprintf("%v", fields[i].error)
			if fields[i].error != "" {
				valid = false
				fields[i].error = "string doesn't match"
				fmt.Println(field.name + "string doesn't match")
			}
		case []string:
			fmt.Println("within string slice")
			for key, thingy := range fieldValue {
				fmt.Println(key, thingy)
			}
			//			fields[i].value = fieldValue
		default:
			fmt.Printf("unexpected type %T", field.kind) // %T prints whatever type t is
		}
	}
	return fields, valid
}
