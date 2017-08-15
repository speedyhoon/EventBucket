package main

import (
	"net"
	"net/http"
	"os"
	"strings"
)

func about(w http.ResponseWriter, r *http.Request) {
	templater(w, page{
		Title: "About",
		Data:  localIPs(),
	})
}

//localIPs returns the non loopback local IPv4 of the host
func localIPs() map[string]interface{} {
	if isPrivate {
		return map[string]interface{}{}
	}

	var localIPs []string
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		var ipnet *net.IPNet
		var ok bool
		for _, address := range addrs {
			//Check the address type is not localhost or a loopback address
			ipnet, ok = address.(*net.IPNet)
			if ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil && !strings.HasPrefix(ipnet.IP.String(), "169.254.") {
				localIPs = append(localIPs, ipnet.IP.String()+portAddr)
			}
		}
	}
	hostname, _ := os.Hostname()
	return map[string]interface{}{
		"hostname":    hostname + portAddr,
		"ipAddresses": localIPs,
	}
}
