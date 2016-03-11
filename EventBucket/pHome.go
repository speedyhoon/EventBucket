package main

import (
	"net"
	"net/http"
	"os"
)

func home(w http.ResponseWriter, r *http.Request) {
	listEvents, err := getEvents(onlyOpen)
	templater(w, page{
		Title: "Home",
		Error: err,
		Data: map[string]interface{}{
			"NewEvent":   getFormSession(w, r, eventNew),
			"ListEvents": listEvents,
		},
	})
}

func onlyOpen(event Event) bool {
	return !event.Closed
}

func about(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	templater(w, page{
		Title: "About",
		Data: map[string]interface{}{
			"Hostname":    hostname,
			"IpAddresses": localIPs(),
		},
	})
}

func licence(w http.ResponseWriter, r *http.Request) {
	templater(w, page{
		Title: "Licence",
	})
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
