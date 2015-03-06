package main

func about() Page {
	hostname, ipAddresses := HostnameIpAddresses()
	return Page{
		TemplateFile: "about",
		Title:        "About",
		Theme:        TEMPLATE_HOME,
		Data: M{
			"Version":     VERSION,
			"Menu":        homeMenu(URL_about, HOME_MENU_ITEMS),
			"Hostname":    hostname,
			"IpAddresses": ipAddresses,
			"BuildDate":   BUILDDATE,
			"IconHeight":  30,
		},
	}
}
