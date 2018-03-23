package main

import (
	"net"
	"net/http"
	"os"
	"strings"
	"github.com/speedyhoon/session"
	"github.com/speedyhoon/forms"
)

func about(w http.ResponseWriter, r *http.Request) {
	f, _ := session.Forms(w, r, getFields, settings)
	render(w, page{
		Title: "About",
		Data: map[string]interface{}{
			"settings": f[0],
			"Network":  localIPs(),
		},
	})
}

func settingsUpdate(f forms.Form) (string, error) {
	mainTheme++
	if mainTheme == 2 {
		mainTheme = 0
	}
	return "", nil
}

//localIPs returns the non loopback local IPv4 of the host
func localIPs() map[string]interface{} {
	if isPrivate {
		return map[string]interface{}{}
	}

	var localIPs []string
	addresses, err := net.InterfaceAddrs()
	if err == nil {
		var ipNet *net.IPNet
		var ok bool
		for _, address := range addresses {
			//Check the address type is not localhost or a loopback address
			ipNet, ok = address.(*net.IPNet)
			if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil && !strings.HasPrefix(ipNet.IP.String(), "169.254.") {
				localIPs = append(localIPs, ipNet.IP.String()+portAddr)
			}
		}
	}
	hostname, _ := os.Hostname()
	return map[string]interface{}{
		"hostname":    hostname + portAddr,
		"ipAddresses": localIPs,
	}
}
