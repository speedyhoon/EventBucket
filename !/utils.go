package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

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

// localIP returns the non loopback local IPv4 of the host
func localIPs() []string {
	var localIPs []string
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		var ipnet *net.IPNet
		var ok bool
		for _, address := range addrs {
			// check the address type and if it is not a loopback the display it
			ipnet, ok = address.(*net.IPNet)
			if ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				localIPs = append(localIPs, ipnet.IP.String())
			}
		}
	}
	return localIPs
}

func defaultDate() string {
	return time.Now().Format("2006-01-02")
}

func defaultTime() string {
	return time.Now().Format("15:04")
}

func toB36(b uint64) string {
	return strconv.FormatUint(b, 36)
}

func B36toUint(b string) (uint64, error) {
	return strconv.ParseUint(b, 36, 64)
}

func trimFloat(num float32) string {
	return strings.TrimRight(strings.Trim(fmt.Sprintf("%f", num), "0"), ".")
}
