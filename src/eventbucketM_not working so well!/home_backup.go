

func home(w http.ResponseWriter, r *http.Request) {
	templator("home", home_temp, homeData(getCollection("event")), w)

}
func homeData(event []map[string]interface{}) map[string]interface{} {

	//	datetime := "201302021252"

	//	test, err := time.Parse("2006/02/01 15:04", datetime)
	//	test, err := time.Parse("20060201", fmt.Sprintf("%v", time.Now()))


	this := fmt.Sprintf("%v", time.Now())

	all := []HomeCalendar{}
	for _, row := range event {
		all = append(all, HomeCalendar{
				Name:      row["name"].(string),
				Club:      row["clubId"].(string),
				Day:      row["datetime"].(string),
				Date:     tryTime(row["datetime"]),
				//				Date:      fmt.Sprintf("%v", test),
				//				Date:      fmt.Sprintf("%v", test),
				//				Date:      row["datetime"].(string),
				//				Month:  row["datetime"].(string),
				Month:  this,
				Time:      row["datetime"].(string),
			})
	}

	return map[string]interface{}{
		"Source": all,
		"PageName": "Calendar",
		"Menu": "Menu is not built yet",
	}

}
func tryTime(datetime interface{}) string {
	test, err := time.Parse("200601021504", fmt.Sprintf("%v", datetime))
	//	test, err := time.Parse("2006-01-02 15:04", "2011-01-19 22:15")
	checkError(err)
	//	return test.Month
	return fmt.Sprintf("%v", test.Month)
}

const (
	stdLongMonth      = "January"
	stdMonth          = "Jan"
	stdNumMonth       = "1"
	stdZeroMonth      = "01"
	stdLongWeekDay    = "Monday"
	stdWeekDay        = "Mon"
	stdDay            = "2"
	stdUnderDay       = "_2"
	stdZeroDay        = "02"
	stdHour           = "15"
	stdHour12         = "3"
	stdZeroHour12     = "03"
	stdMinute         = "4"
	stdZeroMinute     = "04"
	stdSecond         = "5"
	stdZeroSecond     = "05"
	stdLongYear       = "2006"
	stdYear           = "06"
	stdPM             = "PM"
	stdpm             = "pm"
	stdTZ             = "MST"
	stdISO8601TZ      = "Z0700"  // prints Z for UTC
	stdISO8601ColonTZ = "Z07:00" // prints Z for UTC
	stdNumTZ          = "-0700"  // always numeric
	stdNumShortTZ     = "-07"    // always numeric
	stdNumColonTZ     = "-07:00" // always numeric
)

const home_temp = `<h1>{{ .PageName}}</h1>
{{if .Menu}} {{XTC .Menu}}{{end}}
{{with .Source}}
    {{range .}}
        Name is {{.Name}} <br>
        Club is {{.Club}} <br>
        Day: {{.Day}}<br>
        Date {{.Date}}<br>
        Time {{.Time}}<br>
        Month {{.Month}}<br>
    {{end}}
{{end}}
================================
`

//<a href=organisers>Organisers</a>
/*
Thursday
3rd
May
EditTest Event
Test Club
3.00am
*/

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

type HomeCalendar struct {
	Name, Club, Day, Date, Month, Time string
}

