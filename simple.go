package main

import (
	"net/http"
	"html/template"
	"fmt"
	"mgo"
	"mgo/bson"
//	"reflect"
	"net/url"
)
var conn *mgo.Database//global variables that can't be a constant
func main(){
	conn = DB()
	http.HandleFunc("/", home)
	http.HandleFunc("/clubs", clubs)
	http.HandleFunc("/startShooting", startShooting)
	http.HandleFunc("/clubInsert", clubInsert)
	http.HandleFunc("/organisers", organisers)
	http.HandleFunc("/organiser", organisers)
	http.ListenAndServe(":80", nil)
}
func DB()*mgo.Database{
	session, err := mgo.Dial("localhost")
	checkErr(err)
//	defer session.Close()
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



func getClub()bson.M{
	//	var result interface{}
	//	err := conn.C("club").Find(bson.M{"name": "Ale"}).One(&result)
//	query := conn.C("club").Find(bson.M{"location": "VIC"})
//	checkErr(err)
////	fmt.Println("Phone:", result)
//return query

	var m bson.M
	err := conn.C("club").Find(nil).One(&m)
	checkErr(err)
	return m
}
func clubs(w http.ResponseWriter, r *http.Request){
	temp := getClub()
	fmt.Print(temp)
//	for value := range temp{
//		fmt.Print("\n\n\n")
//		fmt.Print(value)
//	}
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
		//todo"Location": map[]interface{}{"VIC","TAS","SA","QLD","NSW","ACT","NT","WA"},
		//todo"Discipline": map[]interface{}{"Rifle","Pistol","Shotgun"},
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
	return pane("Create Club",source)
}


func InsertDoc(data interface{}, collection string){
//	temp := map[string]interface{}{
//		"_id":"qwerty232",
//		"0":"name",	//event name
//		"1":123456,	//date time stamp
//		"2":2,			//club id
//		"3": map[string]interface{}{		//shooters
//			"0":map[string]interface{}{			//shooters id
//				"0":"grade",
//				"1":"class",
//				"2":map[string]interface{}{		//scores
//					"0":map[string]interface{}{		//range id
//						"0": 40,		//total
//						"1": 5,		//centers
//						"2": 7,		//x's
//						"3": 2345678976,	//countback
//						"4": 987456788,	//xcountback
//					},
//					"1":map[string]interface{}{		//range id
//						"0": 40,		//total
//						"1": 5,		//centers
//						"2": 7,		//x's
//						"3": 2345678976,	//countback
//						"4": 987456788,	//xcountback
//					},
//				},
//			},
//		},
//		"4":"settings",
//		"5":"team",
//		"6":"teamcat",
//		"7":"handicap?",//or should this go under shooter?
//	}
//	data = temp
//	data["_id"] = bson.M{}

//	temp := map[string]interface{}{
//		"_id": setAutoInc(collection),
//		"name": "Bob D.",
//	}
	if data{
		err := conn.C(collection).Insert(data)
		checkErr(err)
	}
}
//Basic validation to get only the items listed in options, all others are ignored
//func validate(data map[string]interface{}, options map[int]string)interface{}{
//func validate(data url.Values, options map[string]string)interface{}{
func validate(form url.Values, options map[string][]int)interface{}{
	tryThis :=  map[string]string{}
	for option, min_max  := range options{
		array, ok := form[option]
		if ok{
			length := len(array[0])
			if length >= min_max[1] && length <= min_max[2] {
				tryThis[option] = array[0]
			}else if min_max[0]{
				return false
			}
		}else if min_max[0]{
			return false
		}
	}

//	fmt.Print(data)
//	temp_type := make([]string)
//	output := make(map[string]string)
//	for key,check := range data{
////		if value,ok := data[check]; ok {
////			if len(value) == 1{//} && reflect.TypeOf(value){
//		fmt.Print("\nkey=")
//		fmt.Print(key)
//		fmt.Print(" check=")
//		fmt.Print(check)
////				output[key] = check[0]
////			}
////		}
//	}
	return tryThis
}
//http://net.tutsplus.com/tutorials/client-side-security-best-practices/

func clubInsert(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST" {
		options := map[string][]int{
			//Required, minLength, maxLength
			"name":[]int{1,3,99},
			"url":[]int{1,3,26},
			"discipline":[]int{1,5,8},
			"location":[]int{1,3,3},
			"lat":[]int{0,3,15},
			"long":[]int{0,3,15},
		}
//		fmt.Print("received a http POST")
		r.ParseForm()
//		tryThis :=  map[string]string{}
		tryThis := validate(r.Form, options)
//		fmt.Print(r.Form)
//		fmt.Print("\n")
//		fmt.Print(r.Form["name"][0])
//		fmt.Print("\n")
//		fmt.Print(r.Form["location"][0])

//		for _, value := range options{
//			if "" != r.Form[value][0] && 100 > len(r.Form[value][0]){
//				tryThis[value] = r.Form[value][0]
//			}
////		for key, value := range r.Form{
//
//			if len(r.Form[value]) == 1{

//			}
//			fmt.Print("\nvalue=")
//			fmt.Print(r.Form[value][0])
////			fmt.Print("\nkey=")
////			fmt.Print(key)
//			fmt.Print("\t value=")
//			fmt.Print(value)
//			fmt.Print("\t length=")
//			fmt.Print(len(r.Form[value]))
//			fmt.Print("\t NORMAL value=")
//			fmt.Print(value[0])
//			fmt.Print("\t type is ")
//			fmt.Print(reflect.TypeOf(value))
//		}
//		options := map[string]string{
//			"name"		:"[]string",
//			"url"			:"[]string",
//			"discipline":"[]string",
//			"location"	:"[]string",
//			"lat"			:"[]string",
//			"long"		:"[]string",
//		}
//		save := validate(r.Form, options)
//		fmt.Print(save)
//		InsertClub(save)
		InsertDoc(tryThis, "club")
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
<style>body{background:#000}*{color:#555}</style>
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
<style>body{background:#000}*{color:#555}</style>
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
