package main

import (
	"net/http"
	"html/template"
	"fmt"
	"mgo"
)
//global variables that can't be a constant
var conn *mgo.Database

func main(){
	conn = DB()
	http.HandleFunc("/", home)
	http.HandleFunc("/clubs", clubs)
	http.HandleFunc("/startShooting", startShooting)
	http.HandleFunc("/clubInsert", clubInsert)
	http.HandleFunc("/organisers", organisers)
	http.ListenAndServe(":80", nil)
}
func DB()*mgo.Database{
	session, err := mgo.Dial("localhost")
	checkErr(err)
	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	//	session.SetMode(mgo.Monotonic, true)
	session.SetMode(mgo.Eventual, true)//this is supposed to be faster
	return session.DB("eb")
}
func checkErr(err error){
	if err != nil {
		panic(err)
	}
}
func templator(template string, fill_in string, data map[string]interface{}, w http.ResponseWriter){
	switch(template){
	case "home":
		generator(w, homeTemplate(fill_in), data)
		break
	case "admin":
		generator(w, adminTemplate(fill_in), data)
		break
	case "ajax":
		generator(w, fill_in, data)
	}
}
func generator(w http.ResponseWriter, fillin string, data map[string]interface{}){
	my_html := template.New("fillin").Funcs(template.FuncMap {
		"XTC": func(x string) template.HTML {
			return template.HTML(x)
	}})
	t := template.Must(my_html.Parse(fillin))
	err := t.Execute(w, data)
	if err != nil {
		fmt.Print("executing template:", err)
	}
}



func clubs(w http.ResponseWriter, r *http.Request){
	templator("home", clubsView(), clubsData(), w)
}
func clubsData()map[string]interface{}{
	return map[string]interface{}{
		"Source": "Clubs page!",
	}
}
func clubsView()string{
	return `<div id=topBar>
		<h1>Clubs</h1>{{if .Menu}} {{XTC .Menu}}{{end}}
	</div>
	{{if .Source}}{{XTC .Source}}{{end}}`
}



func organisers(w http.ResponseWriter, r *http.Request){
	templator("admin", organisersView(), organisersData(), w)
}
func organisersData()map[string]interface{}{
	return map[string]interface{}{
		"Title": "Organisers",
//		"Source": "Insert Club",
//		"Location": map[]interface{}{"VIC","TAS","SA","QLD","NSW","ACT","NT","WA"},
	}
}
func organisersView()string{
	source := `<form method=post action=clubInsert>
	Club Name:<input name=name>
	Discipline:<select name=discipline>
	<option>Rifle</option>
	<option>Pistol</option>
	<option>Shotgun</option>
	</select>
	Location:<select name=location>
	<option>VIC</option>
	<option>ACT</option>
	<option>NSW</option>
	<option>NT</option>
	<option>QLD</option>
	<option>SA</option>
	<option>TAS</option>
	<option>WA</option>
	</select>

	URL Key:<input name=url>
	Latitude:<input name=lat>
	Longitude:<input name=long>
	<input type=submit value="Add Club">
	</form>`


//	<div id=topBar>
//	<h1>Club Insert</h1>{{if .Menu}} {{XTC .Menu}}{{end}}
//</div>
//Club Name:<input>
//{{if .Source}}{{XTC .Source}}{{end}}




//	<input class="toggle clubNew" data-text="Club Name" value="Club Name">
//	Discipline:<select class="clubNew"><option>Rifle</option><option disabled="">Pistol</option><option disabled="">Shotgun</option></select>
//	Location:<select class="clubNew">
//	<option>VIC</option>
//	<option disabled="">ACT</option>
//	<option disabled="">NSW</option>
//	<option disabled="">NT</option>
//	<option disabled="">QLD</option>
//	<option disabled="">SA</option>
//	<option disabled="">TAS</option>
//	<option disabled="">WA</option>
//	</select> Build subsite:<input type="checkbox" checked="" class="clubNew"> URL Key:<input class="clubNew">
//	<abbr class="help" title="Each club must use an unique URL Key to access their website address.">?</abbr>
//	<a href="clubNew=" class="submit" data-id="clubNew">Add Club</a>
	return pane("Create Club",source)
}


func InsertClub(data interface{}){
	c := conn.C("people")
//	err := c.Insert(&Person{"Ale", "+55 53 8116 9639"},&Person{"Cla", "+55 53 8402 8510"})
	err := c.Insert(data)
	checkErr(err)

//	result := Person{}
//	err = c.Find(bson.M{"name": "Ale"}).One(&result)
//	checkErr(err)
//	fmt.Println("Phone:", result.Phone)
}


func clubInsert(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST" {
		fmt.Print("received a http POST")
//		r.ParseForm()
//		value := r.FormValue("value")
//		for w := range ch {
		for key, value := range r.Form{
//		for index, value := {
			fmt.Print("\nkey=")
			fmt.Print(key)
			fmt.Print(" value=")
			fmt.Print(value)
		}
		templator("admin", clubInsertView(), clubInsertData(), w)
	}else{
		fmt.Print("did not receive a http POST")
	}
}
func clubInsertData()map[string]interface{}{
	return map[string]interface{}{
		"Title": "Insert Club",
		"Source": "Insert Club",
	}
}
func clubInsertView()string{
	return `<div id=topBar>
		<h1>Club Insert</h1>{{if .Menu}} {{XTC .Menu}}{{end}}
	</div>
	Club Name:<input>
	{{if .Source}}{{XTC .Source}}{{end}}`
}






func home(w http.ResponseWriter, r *http.Request){
	templator("home", homeView(), homeData(), w)
}
func homeData()map[string]interface{}{
	return map[string]interface{}{
		"Source": "<select><option>G'day Mate</option></select>",
	}
}
func homeView()string{
	return `<div id=topBar>
		<h1>{{ .PageName}}</h1>{{if .Menu}} {{XTC .Menu}}{{end}}
	</div>
	{{if .Source}}{{XTC .Source}}{{end}}`
}



func startShooting(w http.ResponseWriter, r *http.Request){
	templator("home", startShootingView(), startShootingData(), w)
}
func startShootingData()map[string]interface{}{
	return map[string]interface{}{
		"Source": "<select><option>G'day Mate</option></select>",
	}
}
func startShootingView()string{
	return `<div id=topBar>
		<h1>{{ .PageName}}</h1>{{if .Menu}} {{XTC .Menu}}{{end}}
	</div>
	{{if .Source}}{{XTC .Source}}{{end}}`
}



func homeTemplate(body string)string{
	return `<!doctype html>
<html>
<head>
	{{if .Css}}<link rel=stylesheet href={{.Css}}>{{end}}
	{{if .Ico}}<link rel=icon type={{.IcoType}} href={{.Ico}}>{{end}}
	<title>EventBucket{{if .Title}} - {{.Title}}{{end}}</title>
</head>
<body>
	` + body + `
	{{if .Js}}<script src={{.Js}}></script>{{end}}
</body>
</html>`
}
func adminTemplate(body string)string{
	return `<!doctype html>
<html>
<head>
	{{if .Css}}<link rel=stylesheet href={{.Css}}>{{end}}
	{{if .Ico}}<link rel=icon type={{.IcoType}} href={{.Ico}}>{{end}}
	<title>EventBucket{{if .Title}} - {{.Title}}{{end}}</title>
</head>
<body>
	` + body + `
	{{if .Js}}<script src={{.Js}}></script>{{end}}
</body>
</html>`
}
func pane(title, source string)string{
	return `<h2>`+title+`</h2>
<div>
	`+source+`
</div>`
}
