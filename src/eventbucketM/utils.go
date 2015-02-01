package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/boombuler/barcode/qr"
	"github.com/boombuler/barcode"
	"os"
	"image/png"
//	"io"
	"io/ioutil"
	"encoding/base64"
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
func export(input interface{}) {
	fmt.Printf("\n%#v\n\n", input) //can copy and declare new variable with it. Most ouput available
}

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
func strToInt2(input interface{})int{
	output, _ := strconv.Atoi(fmt.Sprintf("%v", input))
	return output
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

func qrBarcode(value string)string{
	f, _ := os.Create("temp_barcode.png")
	defer f.Close()
	qrcode, err := qr.Encode(value,  qr.L, qr.Auto)
	if err == nil {
		qrcode, err = barcode.Scale(qrcode, 100, 100)
		if err == nil {
			png.Encode(f, qrcode)
			data, err := ioutil.ReadFile("temp_barcode.png")
			if err == nil {
				return "data:image/png;base64," + base64.StdEncoding.EncodeToString(data)
			}
		}
	}
	fmt.Println(err)
	return ""
}
