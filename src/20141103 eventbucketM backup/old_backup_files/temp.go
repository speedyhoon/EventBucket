package main

import (
	"fmt"
	"html/template"
	"os"
)

type HomeCalendar struct {
	Name, Club, Day, Date, Month, Time string
}

const templ2 = `
{{with .Jobs}}
    {{range .}}
        An employer is {{.Name}}
        and the role is {{.Club}}
    {{end}}
{{end}}
`

func main() {
	all := []HomeCalendar{}
	for _, key := range []string{"4", "1"} {
		all = append(all, HomeCalendar{
				Name:      "eeee:"+key,
				Club:      "dddd"+key,
			})
	}

	person := map[string]interface{}{
		"Jobs": all,
	}


	t := template.New("Person template")
	t, err := t.Parse(templ2)
	checkError(err)

	err = t.Execute(os.Stdout, person)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
