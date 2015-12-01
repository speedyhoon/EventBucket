package main

import "strings"

//AddQuotes returns value with or without surrounding single or double quote characters suitable for a [[//dev.w3.org/html5/html-author/#attributes][HTML5 attribute]] value.
func addQuotes(value string) string {
	//TODO escape any rune character over X code point
	//	value = html.EscapeString(value)
	//	escaper:= = strings.NewReplacer(
	//		`&`, "&amp;",
	//`'`, "&#39;", // "&#39;" is shorter than "&apos;" and apos was not in HTML until HTML5.
	//		`<`, "&lt;",
	//		`>`, "&gt;",
	//`"`, "&#34;", // "&#34;" is shorter than "&quot;".
	//	)
	value = strings.Replace(value, `&`, "&amp;", -1) //will destroy any existing escaped characters like &#62;
	double := strings.Count(value, `"`)
	single := strings.Count(value, `'`)
	if single > 0 && single >= double {
		return `"` + strings.Replace(value, `"`, "&#34;", -1) + `"`
	}
	if double > 0 || strings.ContainsAny(value, " `=<>") {
		return `'` + strings.Replace(value, `'`, "&#39;", -1) + `'`
	}
	/*//Contains a single quote and a double quote character.
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
	}*/
	return value
}
