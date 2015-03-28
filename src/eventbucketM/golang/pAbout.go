package main

func about() Page {
	hostname, ipAddresses := hostnameIPAddresses()
	return Page{
		TemplateFile: "about",
		Title:        "About",
		Theme:        templateHome,
		Data: M{
			"Version":     versionNumber,
			"Hostname":    hostname,
			"IpAddresses": ipAddresses,
			"BuildDate":   buildDate,
			"IconHeight":  30,
		},
	}
}
