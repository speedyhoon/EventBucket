package main

func about() Page {
	hostname, ipAddresses := HostnameIpAddresses()
	return Page{
		TemplateFile: "about",
		Title:        "About",
		Theme:        TEMPLATE_HOME,
		Data: M{
			"Version":     VERSION,
			"Hostname":    hostname,
			"IpAddresses": ipAddresses,
			"BuildDate":   BUILDDATE,
			"IconHeight":  30,
		},
	}
}
