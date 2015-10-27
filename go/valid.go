package main

import (
)

func isValidInt (strNum string, input Input)(int, error){
	num, err := strconv.Atoi(strNum)
	if err != nil {
		return num, err
	}
	if num >= int(input.min) && num <= int(input.max) && (input.step == 0 || input.step != 0 && num % int(input.step) == 0)  || !input.required && num == 0 {
		return num, nil
	}
	return num, errors.New("input integer doesn't pass validation")
}
func isValidStr (str string, input Input)(string, error){
	length := len(str)
	if length >= input.minLen && length <= input.maxLen || !input.required && str == ""{
		return str, nil
	}
	return str, errors.New("input string doesn't pass validation")
}

func isValid(fields []Input, r *http.Request) ([]Input, bool) {
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
			continue  //Skip to the next loop iteration.
		}

		switch field.kind.(type) {
		case bool:
			fmt.Printf("boolean %t\n", field.kind)
		case int:
			fields[i].value, err = isValidInt(strings.TrimSpace(fieldValue[0]), field)
			if err == nil {
				fields[i].isValid = true
			}else {
				valid = false
				fields[i].errMsg = "integer supplied was wrong"
				fmt.Println(field.name + "integer supplied was wrong")
			}
		case string:
			fields[i].value, fields[i].err = isValidStr(strings.TrimSpace(fieldValue[0]), field)
			fields[i].isValid = fields[i].err == nil
			fields[i].errMsg = fmt.Sprintf("%v", fields[i].err)
			if !fields[i].isValid {
				valid = false
				fields[i].errMsg = "string doesn't match"
				fmt.Println(field.name + "string doesn't match")
			}
		case []string:
			fmt.Println("within string slice")
			for key, thingy := range fieldValue {
				fmt.Println(key, thingy)
			}
			fields[i].value = fieldValue
		default:
			fmt.Printf("unexpected type %T", field.kind)		// %T prints whatever type t is
		}
	}
	return fields, valid
}