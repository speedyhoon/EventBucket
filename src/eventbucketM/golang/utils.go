package main

import (
	//"bytes"
	//"encoding/base64"
	//"errors"
	"fmt"
	//"github.com/boombuler/barcode"
	//"github.com/boombuler/barcode/datamatrix"
	//"github.com/boombuler/barcode/qr"
	//"image/png"
	"net"
	"os"
	"strconv"
	"strings"
)

/*func exists(dict M, key string) string {
	if val, ok := dict[key]; ok {
		return fmt.Sprintf("%v", val)
	}
	return ""
}*/

func str2Int(input interface{}) int {
	number, err := strconv.Atoi(fmt.Sprintf("%v", input))
	if err != nil {
		return number
	}
	return 0
}

func str2float(input interface{}) float64 {
	float, err := strconv.ParseFloat(fmt.Sprintf("%v", input), 64)
	if err != nil {
		return float
	}
	return 0
}

//research http://net.tutsplus.com/tutorials/client-side-security-best-practices/
func addQuotes(input string) string {
	if strings.Contains(input, " ") {
		return "\"" + input + "\""
	}
	return input
}

func isValueInSlice(needle int, haystack []int) bool {
	for _, value := range haystack {
		if needle == value {
			return true
		}
	}
	return false
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

/*func ternary(condition bool, True, False string) string {
	if condition {
		return True
	}
	return False
}*/

/*func imgBarcode(width, height int, barcodeType, value string) string {
	var data barcode.Barcode
	var err error
	switch barcodeType {
	case QRCODE:
		data, err = qr.Encode(value, qr.L, qr.Auto)
		break
	case DATAMATRIX:
		data, err = datamatrix.Encode(value)
		break
	default:
		err = errors.New("barcode type " + barcodeType + " is not implemented!")
		break
	}
	if err == nil {
		data, err = barcode.Scale(data, width, height)
		if err == nil {
			var buf bytes.Buffer
			err = png.Encode(&buf, data)
			if err == nil {
				return fmt.Sprintf("<img src=\"data:image/png;base64,%v\" width=%v height=%v alt=%v/>", base64.StdEncoding.EncodeToString(buf.Bytes()), width, height, value)
			}
		}
	}
	Error.Println(err)
	return ""
}*/

// dirExists returns a bool whether the given directory exists or not
func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		return true
	}
	if !os.IsNotExist(err) {
		Error.Printf("folder does not exist: %v", err)
	}
	return false
}

func hostnameIPAddresses() (string, []string) {
	hostname, _ := os.Hostname()
	var ipAddress []string
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, i := range interfaces {
			addrs, err2 := i.Addrs()
			if err2 == nil {
				for _, addr := range addrs {
					ipAddress = append(ipAddress, fmt.Sprintf("%v", addr))
				}
			}
		}
	}
	return hostname, ipAddress
}
