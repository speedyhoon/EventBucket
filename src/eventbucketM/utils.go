package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/boombuler/barcode/qr"
	"github.com/boombuler/barcode"
	"os"
	"image/png"
	"io/ioutil"
	"encoding/base64"
)

/*func exists(dict M, key string) string {
	if val, ok := dict[key]; ok {
		return fmt.Sprintf("%v", val)
	}
	return ""
}*/

func strToInt(input interface{})(int, error){
	return strconv.Atoi(fmt.Sprintf("%v", input))
}

//research http://net.tutsplus.com/tutorials/client-side-security-best-practices/
func addQuotes(input string) string {
	if strings.Contains(input, " "){ //}|| input == "/" {	// strings.Contains(input, "/") {
		return "\"" + input + "\""
	}
	return input
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

func qrBarcode(width, height int, value string)string{
	f, err := os.Create("temp_barcode.png")
	if err != nil {
		Warning.Println(err)
	}
	defer f.Close()
	var qrCode barcode.Barcode
	qrCode, err = qr.Encode(value,  qr.L, qr.Auto)
	if err == nil {
		qrCode, err = barcode.Scale(qrCode, width, height)
		if err == nil {
			png.Encode(f, qrCode)
			data, err := ioutil.ReadFile("temp_barcode.png")
			if err == nil {
				return fmt.Sprintf("<img src=\"data:image/png;base64,%v\" width=%v height=%v alt=%v/>", base64.StdEncoding.EncodeToString(data), width, height, value)
			}
		}
	}
	Warning.Println(err)
	return ""
}
