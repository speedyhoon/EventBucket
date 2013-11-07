
func controller(input string){
	letter := htmltemp()
	recipients := model()
	generater(letter, recipients )
}

func model() ([]Page){
	return []Page{
		{"Aunt Mildred", "bone china tea set", true},
		{"Uncle John", "moleskin pants", false},
		{"Cousin Rodney", "", false},
	}
}

func templater() string{
return `
Dear {{.Name}},
{{if .Attended}}
It           was a pleasure to see you at the wedding.{{else}}
It is a shame you couldn't make it to the wedding.{{end}}
{{with .Gift}}Thank you for the lovely {{.}}.
{{end}}
Best wishes,
Josie
`
}

func htmltemp() string{
	return `
<!doctype html>
<html>
<head>
  <title>EventBucket{{if .Name}} - {{.Name}}{{else}} Gooo Web Framework v0.1{{end}}</title>
</head>
<body>
{{if .Attended}}
It        was a pleasure to see you at the wedding.{{else}}
It is a shame you couldn't make it to the wedding.{{end}}
{{with .Gift}}Thank you for the lovely {{.}}.
{{end}}
</body>
</html>
`
}
