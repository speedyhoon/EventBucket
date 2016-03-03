package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func isValidUint64(inp []string, field field) (interface{}, string) {
	strNum := strings.TrimSpace(inp[0])
	var num uint64
	num = 0
	if field.step == 0 {
		warn.Println("Are you sure about step == 0?")
		return num, "Step supplied = 0"
	}
	var err error
	if strNum != "" {
		num, err = strconv.ParseUint(strNum, 10, 64)
		if err != nil {
			return num, err.Error()
		}
	} else if field.Required {
		return num, "Please fill in this field"
	}
	if field.max == 0 {
		field.max = math.MaxInt32
	}
	if num >= uint64(field.min) && num <= uint64(field.max) && num%uint64(field.step) == 0 || !field.Required && num == 0 {
		return num, ""
	}
	return num, "field integer doesn't pass validation"
}

func isValidFloat64(inp []string, field field) (interface{}, string) {
	strNum := strings.TrimSpace(inp[0])
	if field.step == 0 {
		warn.Println("Are you sure about step == 0?")
		return 0, "Step supplied = 0"
	}
	var num float64
	var err error
	if strNum != "" {
		num, err = strconv.ParseFloat(strNum, 64)
		if err != nil {
			return num, err.Error()
		}
	} else if field.Required {
		return num, "Please fill in this field"
	}
	if num >= float64(field.min) && num <= float64(field.max) && (field.step == 0 || field.step != 0 /*&& num%step == 0*/) || !field.Required && num == 0 {
		return num, ""
	}
	return num, "field integer doesn't pass validation"
}

func isValidFloat32(inp []string, field field) (interface{}, string) {
	strNum := strings.TrimSpace(inp[0])
	if field.step == 0 {
		warn.Println("Are you sure about step == 0?")
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
		return num, "Please fill in this field"
	}
	if num >= float32(field.min) && num <= float32(field.max) && (field.step == 0 || field.step != 0 /*&& num%step == 0*/) || !field.Required && num == 0 {
		return num, ""
	}
	return num, "field integer doesn't pass validation"
}

func isValidStr(inp []string, field field) (interface{}, string) {
	//TODO check if the string passes regex
	//TODO check if it matches one of the options provided
	str := strings.TrimSpace(inp[0])
	length := uint(len(str))
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
	return str, "Please fill in this field"
}
func isValidID(inp []string, field field) (interface{}, string) {
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

func isValidBool(inp []string, field field) (interface{}, string) {
	str := strings.TrimSpace(inp[0])
	checked := len(str) >= 1
	if field.Required && !checked {
		return false, "Please check this field"
	}
	return checked, ""
}

func isValidRangeIDs(ranges []string, field field) (interface{}, string) {
	var rangeIDs []uint64
	var num uint64
	var err error
	for _, r := range ranges {
		num, err = strconv.ParseUint(r, 10, 64)
		if err != nil {
			return []uint64{}, "Contains invalid range ids."
		}
		rangeIDs = append(rangeIDs, num)
	}
	return rangeIDs, ""

	//	trace.Println("TODO")
	//	return "fds", "fsd"
	//get eventid
	//get list of ranges
	//get their ids
	//check each
	// 			[]range has index str
}
