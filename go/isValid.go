package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func listUint(f *field, inp ...string) {
	//TODO add a minimum qty of items. most lists should be at least one or two items long.
	check := make(map[uint]bool, len(inp))
	var list []uint

	//Use a temporary field as a pointer so isValidUint can assign values & errors.
	var g field

	for _, in := range inp {
		trace.Println("unvalidated rangeID", in, "Isn't empty", in != "")

		g.Error = ""
		isValidUint(&g, in)
		if g.Error != "" {
			f.Error = "This value is invalid."
			return
		}

		_, ok := check[g.valueUint]
		if ok {
			f.Error = "Duplicate value found in list"
			return
		}
		check[g.valueUint] = true
		list = append(list, g.valueUint)
	}

	f.valueUintSlice = list
	return
}

func isValidUint(f *field, inp ...string) {
	if debug {
		if f.step == 0 {
			warn.Println("Are you sure about step == 0? isValidUint", f.name)
		}
		if f.max == 0 {
			warn.Println("Are you sure about max == 0? isValidUint", f.name)
		}
	}

	//TODO switch between 64 and 32 on different architectures.
	n64, err := strconv.ParseUint(strings.TrimSpace(inp[0]), 10, 64)
	if err != nil {
		//Return error if input string failed to convert.
		f.Error = err.Error()
		return
	}
	num := uint(n64)

	if !f.Required && num == 0 {
		//f.valueUint is zero by default so assigning zero isn't required
		return
	}
	if num < uint(f.min) || num > uint(f.max) {
		f.Error = fmt.Sprintf("Must be between %d and %d", f.min, f.max)
		return
	}
	if num%uint(f.step) != 0 {
		//TODO calculate next and previous valid values
		f.Error = "Please enter a valid value. The two nearest values are %d and %d"
		return
	}
	f.valueUint = num
	return
}

func isValidFloat32(f *field, inp ...string) {
	if f.step == 0 {
		warn.Println("Are you sure about step == 0? isValidFloat32", f.name)
	}
	if f.max == 0 {
		warn.Println("Are you sure about max == 0? isValidFloat32", f.name)
	}

	f64, err := strconv.ParseFloat(strings.TrimSpace(inp[0]), 32)
	if err != nil {
		//Return error if input string failed to convert.
		f.Error = err.Error()
		return
	}
	num := float32(f64)

	if !f.Required && num == 0 {
		//f.ValueFloat32 is zero by default so assigning zero isn't required
		return
	}
	if num < float32(f.min) || num > float32(f.max) {
		f.Error = fmt.Sprintf("Must be between %d and %d", f.min, f.max)
		return
	}

	if math.Mod(float64(num), float64(f.step)) != 0 {
		//TODO calculate next and previous valid values
		f.Error = "Please enter a valid value. The two nearest values are %d and %d"
		return
	}
	f.valueFloat32 = num
	return
}

func isValidStr(f *field, inp ...string) {
	//Developer checks
	if f.maxLen == 0 {
		trace.Println("f.maxLen should be set: isValidStr", f.name)
	}
	if f.minLen == 0 {
		trace.Println("f.minLen should be set: isValidStr", f.name)
	}

	f.Value = strings.TrimSpace(inp[0])
	length := len(f.Value)

	//Check value matches regex
	if f.regex != nil && !f.regex.MatchString(f.Value) {
		f.Error = "Failed pattern"
		return
	}

	if length > f.maxLen {
		f.Error = fmt.Sprintf("Please shorten this text to %d characters or less (you are currently using %d character%v).", f.maxLen, length, plural(length))
		return
	}
	if length < f.minLen || length > f.maxLen {
		f.Error = fmt.Sprintf("Please lengthen this text to %d characters or more (you are currently using %d character%v).", f.minLen, length, plural(length))
		return
	}

	//Check value matches one of the options (optional).
	if len(f.Options) > 0 {
		matched := false
		for _, option := range f.Options {
			matched = option.Value == f.Value
			if matched {
				break
			}
		}
		if !matched {
			f.Error = "Value doesn't match any of the options"
			return
		}
	}
	return
}

//TODO look into using isValidStr instead of isValidID
func isValidID(f *field, inp ...string) {
	//TODO remove developer check
	if f.regex == nil {
		trace.Println("missing regex for field:", f.name)
		f.Error = "Missing regex to check against"
		return
	}

	f.Value = strings.TrimSpace(inp[0])
	if !f.regex.MatchString(f.Value) {
		f.Error = "ID supplied is incorrect"
	}
	return
}

func isValidBool(f *field, inp ...string) {
	f.Checked = len(strings.TrimSpace(inp[0])) >= 1
	if f.Required && !f.Checked {
		f.Error = "Please check this field"
	}
	return
}

func plural(length int) string {
	if length != 1 {
		return "s"
	}
	return ""
}
