package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

const (
	fillItIn = "Please fill in this field"
)

func listUint64(field field, inp ...string) (interface{}, string) {
	//TODO add a minimum qty of items. most lists should be at least one or two items long.
	check := make(map[uint64]bool, 1)
	var ids []uint64
	if field.Required && len(inp) == 0 {
		return check, "There are no items selected"
	}
	var oo string
	var temp interface{}
	var id uint64
	for _, in := range inp {
		trace.Println("unvalidated rangeID", in, "Isn't empty", in != "")
		temp, oo = isValidUint64(field, in)
		if oo != "" {
			return ids, "This value is invalid."
		}
		id = temp.(uint64)
		//Using a map here to prevent duplicates
		_, ok := check[id]
		if ok {
			return ids, "duplicate id found in list"
		}
		check[id] = true
		ids = append(ids, id)
	}
	return ids, ""
}

func isValidUint(field field, inp ...string)(interface{}, string){
	if debug{
		if field.step == 0 {
			warn.Println("Are you sure about step == 0? isValidUint", field.name)
			return 0, "Step supplied = 0"
		}
		if field.max == 0{
			warn.Println("Are you sure about max == 0? isValidUint", field.name)
		}
	}

	//TODO switch between 64 and 32 on different architectures.
	n64, err := strconv.ParseUint(strings.TrimSpace(inp[0]), 10, 64)
	num := uint(n64)
	if err != nil {
		//Return error if input string failed to convert.
		return num, err.Error()
	}
	if !field.Required && num == 0{
		return 0, ""
	}
	if num < uint(field.min) || num > uint(field.max){
		return num, fmt.Sprintf("Must be between %d and %d", field.min, field.max)
	}
	if num % uint(field.step) != 0{
		return num, "Please enter a valid value. The two nearest values are %d and %d"
	}
	return num, ""
}

func isValidUint64(field field, inp ...string) (interface{}, string) {
	strNum := strings.TrimSpace(inp[0])
	var num uint64
	num = 0
	if debug && field.step == 0 {
		warn.Println("Are you sure about step == 0? isValidUint64", field.name)
		return num, "Step supplied = 0"
	}
	var err error
	if strNum != "" {
		num, err = strconv.ParseUint(strNum, 10, 64)
		if err != nil {
			return num, err.Error()
		}
	} else if field.Required {
		return num, fillItIn
	}
	if field.max == 0 {
		if num >= uint64(field.min) && num <= math.MaxUint64 && num%uint64(field.step) == 0 || !field.Required && num == 0 {
			return num, ""
		}
	}
	if num >= uint64(field.min) && num <= uint64(field.max) && num%uint64(field.step) == 0 || !field.Required && num == 0 {
		return num, ""
	}
	return num, "field integer doesn't pass validation"
}

func isValidFloat32(field field, inp ...string) (interface{}, string) {
	strNum := strings.TrimSpace(inp[0])
	if debug && field.step == 0 {
		warn.Println("Are you sure about step == 0? isValidFloat32", field.name)
		return 0, "Step supplied = 0"
	}
	var num float32
	var x float64
	var err error
	if strNum != "" {
		x, err = strconv.ParseFloat(strNum, 64)
		if err != nil {
			return num, err.Error()
		}
		num = float32(x)
	} else if field.Required {
		return num, fillItIn
	}
	if num >= float32(field.min) && num <= float32(field.max) && (field.step == 0 || field.step != 0 /*&& num%step == 0*/) || !field.Required && num == 0 {
		return num, ""
	}
	return num, "field integer doesn't pass validation"
}

func isValidStr(field field, inp ...string) (interface{}, string) {
	//TODO check if the string passes regex
	//TODO check if it matches one of the options provided
	str := strings.TrimSpace(inp[0])
	length := uint(len(str))
	if field.Required && length == 0 {
		return str, fillItIn
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
	return str, fillItIn
}
func isValidID(field field, inp ...string) (interface{}, string) {
	str := strings.TrimSpace(inp[0])
	if field.regex == nil {
		trace.Println("missing regex for field:", field.name)
		return str, "Missing regex to check against"
	}
	if field.regex.MatchString(str) {
		return str, ""
	}
	return str, "ID supplied is incorrect"
}

func isValidBool(field field, inp ...string) (interface{}, string) {
	str := strings.TrimSpace(inp[0])
	checked := len(str) >= 1
	if field.Required && !checked {
		return false, "Please check this field"
	}
	return checked, ""
}
