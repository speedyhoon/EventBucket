package main

import (
	"fmt"
	"strconv"
	"strings"
)

//research http://net.tutsplus.com/tutorials/client-side-security-best-practices/
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
func dump(input interface{}) {
	fmt.Printf("\n%v\n", input)
}
//func vardump(input interface{}) {
//	fmt.Printf("%+v\n", input) //map field names included
//}
//func export(input interface{}) {
//	fmt.Printf("\n%#v\n\n", input) //can copy and declare new variable with it. Most ouput available
//}

func exists(dict M, key string) string {
	if val, ok := dict[key]; ok {
		return fmt.Sprintf("%v", val)
	}
	return ""
}

/*func echo(input interface{}) string {
	return fmt.Sprintf("%v", input)
}*/
func strToInt(input string)(int, bool){
	output, err := strconv.Atoi(input)
	if err != nil {
		return -1, false
	}
	return output, true
}
/*func strToInt64(input string) int64 {
	output, err := strconv.ParseInt(input, 10, 64)
	checkErr(err)
	return output
}*/

func addQuotes(input string) string {
	if strings.Contains(input, " ") || input == "/" {	// strings.Contains(input, "/") {
		return "\"" + input + "\""
	}
	return input
}
func addQuotesEquals(input string) string {
	if input != "" {
		if strings.Contains(input, " ") {
			return "=\"" + input + "\""
		}
		return "=" + input
	}
	fmt.Println("addQuotesEquals had an empty parameter!")
	return ""
}

// Ordinal gives you the input number in a rank/ordinal format.
// Ordinal(3) -> 3rd
//author github.com/dustin/go-humanize/blob/master/ordinals.go
func ordinal(x int) string {
	suffix := "th"
	switch x % 10 {
	case 1:
		if x%100 != 11 {
			suffix = "st"
		}
	case 2:
		if x%100 != 12 {
			suffix = "nd"
		}
	case 3:
		if x%100 != 13 {
			suffix = "rd"
		}
	}
	return strconv.Itoa(x) + suffix
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
