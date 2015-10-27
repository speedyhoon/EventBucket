package main

import (
	"strings"
)

//dev.w3.org/html5/html-author/#attributes
func addQuotes(value string) string {
	//TODO add support for a value with single & double quotes
	//Space, single quote, accent, equals, less-than sign, greater-than sign.
	if strings.ContainsAny(value, "'\"") {
		warn.Printf("HTML attribute value %v contains single & double quotes", value)
	}
	if strings.ContainsAny(value, " '`=<>") {
		return "\"" + value + "\""
	}
	if strings.Contains(value, "\"") {
		return "'" + value + "'"
	}
	return value
}
