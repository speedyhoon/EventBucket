package main

import (
	"strings"
)

//AddQuotes returns value with or without surrounding single or double quote characters suitable for a [[//dev.w3.org/html5/html-author/#attributes][HTML5 attribute]] value.
func addQuotes(value string) string {
	//Contains a single quote and a double quote character.
	if strings.Contains(value, "'") && strings.Contains(value, `"`) {
		warn.Printf("HTML attribute value %v contains both single & double quotes", value)
	}
	//Space, single quote, accent, equals, less-than sign, greater-than sign.
	if strings.ContainsAny(value, " '`=<>") {
		return `"` + value + `"`
	}
	//Double quote
	if strings.Contains(value, `"`) {
		return "'" + value + "'"
	}
	return value
}
