package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func isValidInt(strNum string, field field) (int, error) {
	num, err := strconv.Atoi(strNum)
	if err != nil {
		return num, err
	}
	if num >= field.min && num <= field.max && (field.step == 0 || field.step != 0 && num%field.step == 0) || !field.Required && num == 0 {
		return num, nil
	}
	return num, errors.New("field integer doesn't pass validation")
}
func isValidStr(str string, field field) (string, error) {
	info.Println(field.minLen, field.maxLen, field.Required, `"`+str+`"`)
	length := len(str)
	if length >= field.minLen && length <= field.maxLen || !field.Required && str == "" {
		info.Println(str, "ok")
		return str, nil
	}
	warn.Println(str, "fail")
	return str, errors.New("field string doesn't pass validation")
}

func isValid(r *http.Request, fields []field) ([]field, bool) {
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
			if field.Required && field.defValue != nil {
				fieldValue = field.defValue()
			} else {
				valid = false
				field.Error = "Missing value"
				continue //Skip to the next loop iteration.
			}
		}

		switch field.kind.(type) {
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
		}
	}
	return fields, valid
}
