package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func listUint(f *field, inp ...string) {
	if len(inp) < f.minLen {
		f.Error = fmt.Sprintf("Not enough items selected. At least %v items are needed.", f.minLen) //TODO plural
		return
	}

	check := make(map[uint]bool, len(inp))
	var list []uint

	//Use a temporary field as a pointer so isValidUint can assign values & errors.
	g := *f

	for _, in := range inp {
		//t.Println("unvalidated rangeID", in, "Isn't empty", in != "")

		g.Error = ""
		isValidUint(&g, in)
		if g.Error != "" {
			f.Error = g.Error
			return
		}

		_, ok := check[g.valueUint]
		if ok {
			f.Error = "Duplicate values found in list."
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

	var err error
	f.Value = inp[0]
	f.valueUint, err = strToUint(f.Value)
	if err != nil {
		//Return error if input string failed to convert.
		f.Error = err.Error()
		return
	}

	if !f.Required && f.valueUint == 0 {
		//f.valueUint is zero by default so assigning zero isn't required
		return
	}
	if f.valueUint < uint(f.min) || f.valueUint > uint(f.max) {
		f.Error = fmt.Sprintf("Must be between %v and %v.", f.min, f.max)
		return
	}
	if f.valueUint%uint(f.step) != 0 {
		below := f.valueUint - (f.valueUint % uint(f.step))
		f.Error = fmt.Sprintf("Please enter a valid value. The two nearest values are %d and %d.", below, below+uint(f.step))
		return
	}
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
	if num < f.min || num > f.max {
		f.Error = fmt.Sprintf("Must be between %v and %v.", f.min, f.max)
		return
	}

	if rem := toFixed(math.Mod(f64, float64(f.step))); rem != 0 {
		f.Error = fmt.Sprintf("Please enter a valid value. The two nearest values are %v and %v.", num-rem, num-rem+f.step)
		return
	}
	f.Value = fmt.Sprintf("%v", num)
	f.valueFloat32 = num
	return
}

func toFixed(num /*, precision*/ float64) float32 {
	//output := math.Pow(10, precision)
	//return float64(int(num * output)) / output
	return float32(int(num*1000000)) / 1000000
}

func isValidStr(f *field, inp ...string) {
	//Developer checks
	if f.maxLen == 0 {
		t.Println("f.maxLen should be set: isValidStr", f.name)
	}
	if f.minLen == 0 && f.Required {
		t.Println("f.minLen should be set: isValidStr", f.name)
	}

	f.Value = strings.TrimSpace(inp[0])
	length := len(f.Value)

	//Check value matches regex
	if f.regex != nil && !f.regex.MatchString(f.Value) {
		f.Error = "Failed pattern"
		return
	}

	if length < f.minLen {
		f.Error = fmt.Sprintf("Please lengthen this text to %d characters or more (you are currently using %d character%v).", f.minLen, length, plural(length))
		return
	}
	if length > f.maxLen {
		//Truncate string instead of raising an error
		f.Value = f.Value[:f.maxLen]
		//f.Error = fmt.Sprintf("Please shorten this text to %d characters or less (you are currently using %d character%v).", f.maxLen, length, plural(length))
		//return
	}

	//Check value matches one of the options (optional).
	/*if len(f.Options) > 0 {
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
	}*/
	return
}

func isValidRegex(f *field, inp ...string) {
	//TODO remove developer check
	if f.regex == nil {
		t.Println("missing regex for field:", f.name)
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
		f.Error = "Please check this field."
	}
	return
}

func plural(length int) string {
	if length != 1 {
		return "s"
	}
	return ""
}
