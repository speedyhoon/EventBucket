package main
import (
	"net/http"
	"html/template"
	"fmt"
	"mgo"
//	"mgo/bson"
//	"reflect"
	"net/url"
//	"strings"
)
var conn *mgo.Database//global variables that can't be a constant
const table_style = ".table{display:table}.tr{display:table-row}.td,.th{display:table-cell}.th{font-weight:900;text-align:center}"
func main(){
	conn = DB()
//	http.ResponseWriter.Header().Set("Expires", "Thu, 4 Apr 2013 20:00:00 GMT")
	http.HandleFunc("/", home)
	http.HandleFunc("/clubs", clubs)
	http.HandleFunc("/startShooting", startShooting)
	http.HandleFunc("/clubInsert", clubInsert)
	http.HandleFunc("/eventInsert", eventInsert)
	http.HandleFunc("/organisers", organisers)
	http.HandleFunc("/organiser", organisers)
	http.HandleFunc("/eventShow", eventShow)
	http.HandleFunc("/eventSetup", eventSetup)
	http.ListenAndServe(":80", nil)
}




func eventSetup(w http.ResponseWriter, r *http.Request){
	templator("admin", eventSetupTemplate(), eventSetupData(), w)
}
func eventSetupData()map[string]interface{}{
	clubs := getCollection("club")
	club_list := "<table><tr><th>Name</th></tr>"
	for _, row := range clubs{
		club_list += fmt.Sprintf("<tr><td><a href=%v>%v</a></td></tr>",row["url"],row["name"])
	}
	club_list += "</table>"
	return map[string]interface{}{
		"Title": "Event Setup",
		"Ranges": generateForm(getRangesForm()),
		"CreateRange": generateForm(getCreateRangesForm()),
		//		"Events": generateForm("eventInsert", getCreateEventForm(clubs)),
		//		"Clubs": generateForm("clubInsert", getOrganisersForm()),
		//		"ClubList": club_list,
	}
}
func eventSetupTemplate()string{
	return`<h1>Settings{{if .Title}} - {{ .Title}}{{end}}</h1>
		{{if .Ranges}}` + panel("Ranges",
			`<div class=table><div class=tr>
			<span class=th>Id</span>
			<span class=th>Range Name</span>
			<span class=th></span>
			<span class=th>Modify</span>
			<span class=th></span>
			<span class=th>Agg Ranges <abbr class="help" title="A list of range and/or agg Id's seperated by commas (without spaces), to calculate the sum for each shooter.">?</abbr></span>
			<span class=th>Show on<br>Scoreboard <abbr class="help" title="When ticked, the range or agg will be displayed on the scoreboard. When unticked the range or agg will be absent from the scoreboard, but all information will still remain in the system.">?</abbr></span>
			<span class=th>Enabled <abbr class="help" title="Unticked checkboxes will disable the ability to insert or update any shooters scores for that range. This functionality has been temporarily disabled.">?</abbr></span></tr>
			</div>
			{{XTC .Ranges}}</div>`) + `{{end}}
		{{if .CreateRange}}` + pane("Create Range", `{{XTC .CreateRange}}`) + `{{end}}`
}
func getRangesForm()map[string]InputTypers{
	return map[string]InputTypers{
		"form":InputTypers{
			action:	"eventInsert",
			method:	"post",
			table:	true,
		},
		"range_name":InputTypers{
			Type:		"name",
			Html:		"text",
		},
		"up":InputTypers{
			Type:		"submit",
			Html:		"submit",
			Value:	"Up",
		},
		"down":InputTypers{
			Type:		"submit",
			Html:		"submit",
			Value:	"Down",
		},
		"delete":InputTypers{
			Type:		"submit",
			Html:		"submit",
			Value:	"Delete",
		},
		"agg_ranges":InputTypers{
			Type:		"name",
			Html:		"text",
		},
		"scoreboard":InputTypers{
			Type:		"checkbox",
			Html:		"checkbox",
			Disabled:	true,
			Checked:		true,
		},
		"enabled":InputTypers{
			Type:		"checkbox",
			Html:		"checkbox",
			Disabled:	true,
			Checked:		true,
		},
	}
}
func getCreateRangesForm()map[string]InputTypers{
	return map[string]InputTypers{
		"form":InputTypers{
			action:	"rangeCreate",
			method:	"post",
			table:	false,
		},
//		"name":InputTypers{
//			Type:		"name",
//			Html:		"text",
//		},
//		"type":InputTypers{
//			Type:		"submit",
//			Html:		"submit",
//			Value:	"Up",
//		},
//		"position":InputTypers{
//			Type:		"submit",
//			Html:		"submit",
//			Value:	"Down",
//		},
//		"ranges":InputTypers{
//			Type:		"submit",
//			Html:		"submit",
//			Value:	"Delete",
//		},
//		"agg_ranges":InputTypers{
//			Type:		"name",
//			Html:		"text",
//		},
//		"scoreboard":InputTypers{
//			Type:		"checkbox",
//			Html:		"checkbox",
//			Disabled:	true,
//			Checked:		true,
//		},
//		"enabled":InputTypers{
//			Type:		"checkbox",
//			Html:		"checkbox",
//			Disabled:	true,
//			Checked:		true,
//		},
		"name":InputTypers{
			Type:		"name",
			Html:		"text",
			Label:	"Range Name",
		},
		"type":InputTypers{
//			Type:		"name",
			Html:		"select",
			Label:	"Range Type",
			SelectValues:	map[string]string{"1":"Scoring Range", "2":"Aggregate Total"},
		},
		"position":InputTypers{
			Label:	"Position",
			Html:		"select",
			SelectValues:	map[string]string{"1":"Before range1", "2":"Between range1 and range2", "3":"After range2"},
		},
		"ranges":InputTypers{
//			Type:		"name",
			Html:		"select",
			Label:	"Aggregate Ranges",
			SelectValues:	map[string]string{"1":"Before range1", "2":"Between range1 and range2", "3":"After range2"},
			MultiSelect:	true,
		},
	}
}



//
//<div class=table>
//	<form class=tr method="post" action="blah.html">
//		<span class=td><input type="text"></span>
//		<span class=td><input type="text"></span>
//	</form>



//
//<div><h4>Create Range</h4><a href="addRange=1~" class="submit" data-id="addRange">Create new Range</a> <input class="addRange toggle" data-text="Range Name - Distance" value="Range Name - Distance">
//
//<abbr class="help" title="A set distance that shooters fire from and scores are recored shot by shot. If you are shooting from more than one range, use &quot;Create new Agg&quot; to calculate the totals.">?</abbr>
//<br><br>
//<a href="addAgg=1~" data-id="addAgg" class="submit">Create new Agg</a> <input class="addAgg toggle" data-text="Agg Name" value="Agg Name"> <abbr class="help" title="An Aggregate column (or agg for short) is used to sum up each shooters range scores. These are great as Total columns, Day Aggs and Prizemeeting Grand Aggregates. You can select which ranges and/or other aggs to add together with the &quot;Agg Ranges&quot; options.
//
//Championships can be setup by using several aggs and disabling previous ranges on the scoreboard that you no longer wish to display.">?</abbr> place new Range / Agg: <select id="position" class="addRange addAgg"><option value="0">before "300"</option><option value="1">between "300" and "500"</option><option value="2">between "500" and "Total"</option><option value="3" selected="">after "Total"</option></select> <abbr class="help" title="Create each range with a name like &quot;300&quot;, &quot;900yards&quot; or &quot;Sponsors Name&quot;">?</abbr><br></div>








func eventShow(w http.ResponseWriter, r *http.Request){
	templator("admin", eventShowTemplate(), eventShowData(), w)
}
func eventShowData()map[string]interface{}{
	clubs := getCollection("club")
	club_list := "<table><tr><th>Name</th></tr>"
	for _, row := range clubs{
		club_list += fmt.Sprintf("<tr><td><a href=%v>%v</a></td></tr>",row["url"],row["name"])
	}
	club_list += "</table>"
	return map[string]interface{}{
		"Title": "Event Show",
		"CurrentRange": generateForm(getCurrentRangeForm()),
//		"Events": generateForm("eventInsert", getCreateEventForm(clubs)),
//		"Clubs": generateForm("clubInsert", getOrganisersForm()),
//		"ClubList": club_list,
	}
}
func eventShowTemplate()string{
	return`<h1>{{ .Title}}</h1>
	{{if .CurrentRange}}` + pane("Current Range", `{{XTC .CurrentRange}}`) + `{{end}}` //+ `
//	{{if .Clubs}}{{if .ClubList}}` +
//			  panel("Clubs", pane("Create Club", `{{XTC .Clubs}}`) + pane("Create Club",`{{XTC .ClubList}}`)) + `
//	{{end}}{{end}}`
}
func getCurrentRangeForm()map[string]InputTypers{
	return map[string]InputTypers{
		"form":InputTypers{
			action:	"eventInsert",
			table:	false,
		},
		"range":InputTypers{
			Required:true,
			Min:		1,
			Max:		2,
			Type:		"range",
			Html:		"select",
			SelectValues:	getSelectValues("location"),
			Label:	"Current Range",
		},
	}
}





func organisers(w http.ResponseWriter, r *http.Request){
	templator("admin", organisersView(), organisersData(), w)
}
func organisersData()map[string]interface{}{
	clubs := getCollection("club")
	club_list := "<table><tr><th>Name</th></tr>"
	for _, row := range clubs{
		club_list += fmt.Sprintf("<tr><td><a href=%v>%v</a></td></tr>",row["url"],row["name"])
	}
	club_list += "</table>"
	return map[string]interface{}{
		"Title": "Organisers",
		"Events": generateForm(getCreateEventForm(clubs)),
		"Clubs": generateForm(getOrganisersForm()),
		"ClubList": club_list,
	}
}
func organisersView()string{
	return`<h1>{{ .Title}}</h1>
	{{if .Events}}` + panel("Events", pane("Create Event", `{{XTC .Events}}`)) + `{{end}}` + `
	{{if .Clubs}}{{if .ClubList}}` +
	panel("Clubs", pane("Create Club", `{{XTC .Clubs}}`) + pane("Create Club",`{{XTC .ClubList}}`)) + `
	{{end}}{{end}}`
}

func getCreateEventForm(club_list []map[string]interface{})map[string]InputTypers{
	return map[string]InputTypers{
		"form":InputTypers{
			action:	"eventInsert",
			table: false,
		},
		"name":InputTypers{
			Required:true,
			Min:		3,
			Max:		99,
			Type:		"name",//type of validation used upon
			Html:		"text",//input field type
			Label:	"Club Name",
		},
		"clubId":InputTypers{
			Required:true,
			Min:		1,
			Max:		24,
			Type:		"id",
			SelectValues:	getClubSelectBox(club_list),
			Html:		"select",
			Label:	"Club",
		},
		"datetime":InputTypers{
			Required:true,
			Min:		12,        //201302020412 = 2013/02/02-16:15
			Max:		12,
			Type:		"datetime",
			Html:		"datetime",
			//			Select:	[]string{"Rifle","Pistol","Shotgun"},//TODO change these to functions eventually!
			Label:	"Date Time",
		},
		"submit":InputTypers{
			Required:false,
			Type:		"submit",
			Html:		"submit",
			Label:	"Add Event",
		},
	}
}

func getOrganisersForm()map[string]InputTypers{
	return map[string]InputTypers{
		"form":InputTypers{
			action:	"clubInsert",
			table:	false,
		},
		"name":InputTypers{
			Required:true,
			Min:		3,
			Max:		99,
			Type:		"name",//type of validation used upon
			Html:		"text",//input field type
			Label:	"Club Name",
			PlaceHolder:	"Club Name",
		},
		"url":InputTypers{
			Required:true,
			Min:		3,
			Max:		26,
			Type:		"url",
			Html:		"text",
			Label:	"URL Key",
			Help:		"Description to\nhelp!",
		},
		"discipline":InputTypers{
			Required:true,
			Min:		5,
			Max:		8,
			Type:		"discipline",
			Html:		"select",
			Select:	getSelect("discipline"),
			Label:	"Discipline",
		},
		"location":InputTypers{
			Required:true,
			Min:		1,
			Max:		2,
			Type:		"location",
			Html:		"select",
			SelectValues:	getSelectValues("location"),
			Label:	"Location",
			PlaceHolder:"Temp Select",
		},
		"lat":InputTypers{
			Required:false,
			Min:		3,
			Max:		15,
			Type:		"decimal",
			Html:		"text",
			Label:	"Latitude",
		},
		"long":InputTypers{
			Required:false,
			Min:		3,
			Max:		15,
			Type:		"decimal",
			Html:		"text",
			Label:	"Longitude",
		},
		"submit":InputTypers{
			Required:false,
			Type:		"submit",
			Html:		"submit",
			Label:	"Add Club",
		},
	}
}


















func InsertDoc(data interface{}, collection string){
	if data != false{
		err := conn.C(collection).Insert(data)
		checkErr(err)
	}
}
//Basic validation to get only the items listed in options, all others are ignored
//func validate(form url.Values, options map[string][]int)interface{}{
func validate(form url.Values, options map[string]InputTypers)interface{}{
	tryThis :=  map[string]string{}
	for option, min_max  := range options{
		if options[option].Html != "submit" && options[option].Html != "form"{
			array, ok := form[option]
			if ok{
				length := len(array[0])
				if length >= min_max.Min && length <= min_max.Max {
					tryThis[option] = array[0]
				}else if min_max.Required{
					fmt.Print("\nlength is not within range")
					return false
				}else{
					fmt.Print("\nELSE length is not within range "+option)
				}
			}else if min_max.Required{
				fmt.Print("\nform[option] not in array")
				return false
			}else{
				fmt.Print("\nELSE form[option] not in array "+option)
			}
		}
	}
	return tryThis
}
//research http://net.tutsplus.com/tutorials/client-side-security-best-practices/
func checkErr(err error){
	if err != nil {
		panic(err)
	}
}
type InputTypers struct {
	Required bool
	Min, Max, RangeMin, RangeMax int
	Name, Type string
	Html string
	Select []string
	SelectValues map[string]string
	Label, Help, PlaceHolder string
	AutoCorrect, AutoCapitalize bool
	Value string
	Disabled, Checked, MultiSelect bool


	method, action string
	table bool
}
func addQuotes(input string)string{
	return "\""+input+"\""
}
func generateForm(formData map[string]InputTypers)string{
	form := "<form"
	if form_data, ok := formData["form"]; ok {
		form += fmt.Sprintf(" method=%v action=%v", addQuotes(form_data.method), addQuotes(form_data.action))
		if form_data.table == true {
			form += " class=tr"
		}
	}else{
		panic(fmt.Sprintf("A form element is not set in the form object: %v", formData))
	}
	display_as_table := formData["form"].table

//	for attribute, value := range form_attributes{
//		if attribute != "table"{
//			form += fmt.Sprintf(" %v=%v",attribute, addQuotes(value))
//		}
//	}
	form += ">"
//	output := "<form method=post action="+addQuotes(action)+">"
	for inputName, inputData := range formData{
		output := "\n"

		if inputData.Html != "submit" && inputData.Label != ""{
			output += " "+inputData.Label+":"
		}
		if inputData.Html == "text" || inputData.Html == "submit" || inputData.Html == "number" || inputData.Html == "url" || inputData.Html == "datetime"|| inputData.Html == "checkbox"|| inputData.Html == "radio"{
			output += "<input"
			if inputData.Html == "submit" || inputData.Html == "number" || inputData.Html == "url" || inputData.Html == "datetime"|| inputData.Html == "checkbox"|| inputData.Html == "radio"{
				output += " type="+inputData.Html
			}
			if inputData.Html != "submit"{
				output += " name="+inputName
			}
			if inputData.Required{
				output += " required"
			}
			if inputData.Disabled{
				output += " disabled"
			}
			if inputData.Checked{
				output += " checked"
			}
			if inputData.AutoCorrect == false{
				output += " autocorrect=off"
			}
			if inputData.AutoCorrect == true{
				output += " autocorrect=on"
			}
			if inputData.AutoCapitalize == false{
				output += " autocapitalize=off"
			}
			if inputData.AutoCapitalize == true{
				output += " autocapitalize=on"
			}
			if inputData.PlaceHolder != ""{
				output += " placeHolder="+addQuotes(inputData.PlaceHolder)
			}
			if inputData.RangeMin > 0{
				output += fmt.Sprintf(" min=%d", inputData.Min)
			}
			if inputData.RangeMax > 0{
				output += fmt.Sprintf(" max=%d", inputData.Max)
			}
//			if  inputData.Html == "submit"{
//				output += fmt.Sprintf(" value=%v",addQuotes(inputData.Label))
//			}
			if inputData.Value != ""{
				output += fmt.Sprintf(" value=%v",addQuotes(inputData.Value))
			}
			output += ">"
		}else if inputData.Html == "select"{
			output += "<select name="+inputName
			if inputData.Required{
				output += " required"
			}
			if inputData.MultiSelect{
				output += " multiple"
			}
			if inputData.PlaceHolder != ""{
				output += " placeHolder="+addQuotes(inputData.PlaceHolder)
			}
			if inputData.AutoCorrect == false{
				output += " autocorrect=off"
			}
			if inputData.AutoCorrect == true{
				output += " autocorrect=on"
			}
			output += ">"
			for _, option := range inputData.Select{
				output += fmt.Sprintf("\n\t<option>%v</option>",option)
			}
			for value, option := range inputData.SelectValues{
				output += fmt.Sprintf("\n\t<option value=%v>%v</option>",addQuotes(value), option)
			}
			output += "</select>"
		}
		if inputData.Help != "" {
			output += fmt.Sprintf("<abbr title=%v>?</abbr>",addQuotes(inputData.Help))
		}

		if display_as_table{
			form += fmt.Sprintf("<span class=td>%v</span>", output)
		}else{
			form += output
		}
	}
	form += "</form>"
	return form
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










func getCollection(collection_name string)[]map[string]interface{}{
	//TODO add in support to select only the columns required
	var result []map[string]interface{}
	checkErr(conn.C(collection_name).Find(nil).All(&result))
	return result
}



func clubs(w http.ResponseWriter, r *http.Request){
	temp := getCollection("event")
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













































func getSelect(name string)[]string{
	var returned []string
	if name == "discipline"{
		returned = []string{
			"Rifle",
			"Pistol",
			"Shotgun",
		}
	}
	return returned
}
func getSelectValues(name string)map[string]string{
	var returned map[string]string
	if name == "location"{
		returned = map[string]string{
			"7": "VIC",
			"5": "SA",
			"6": "TAS",
			"1": "ACT",
			"2": "NSW",
			"3": "NT",
			"4": "QLD",
			"8": "WA",
		}
	}else if name == "clubs"{
		var result []map[string]interface{}
		c := conn.C("club")
		err := c.Find(nil).All(&result)
		checkErr(err)
		//		var selectOptions []string
		for _, name := range result{
			clubName, Cok := name["name"].(string)
			//			if ok{
			//				selectOptions = append(selectOptions, clubName)
			//			}
			clubId, Iok := name["_id"].(string)
			//			if Iok {
			//				clubId = clubId[13:len(clubId)-2]
			//			}
			fmt.Print("\n CId ")
			//			fmt.Print(        String(name["_id"])          )
			fmt.Print("\n CN ")
			fmt.Print(clubName)
			//			fmt
			if Cok && Iok{
				//				returned[clubId] = append(returned[clubId], clubName)
				returned[ clubId ] = clubName
			}
		}

		fmt.Print("\n\n\n")
		fmt.Print(returned)
	}
	return returned
}
func getClubSelectBox(club_list []map[string]interface{})map[string]string{
	drop_down := make(map[string]string)
	for _, row := range club_list{
		drop_down[fmt.Sprintf("%v",row["_id"])[13:37]] = fmt.Sprintf("%v",row["name"])
	}
	return drop_down
}


func clubInsert(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST" {
		fmt.Print("\ndefinatly a POST")
		r.ParseForm()
		tryThis := validate(r.Form, getOrganisersForm())
		fmt.Print("\nInsert Data")
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

func eventInsert(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST" {
		r.ParseForm()
		clubs := getCollection("club")
		tryThis := validate(r.Form, getCreateEventForm(clubs))
		InsertDoc(tryThis, "event")
		organisers(w, r)
//		templator("admin", eventInsertView(), eventInsertData(), w)
	}else{
		fmt.Print("did not receive a http POST")
	}
}





func home(w http.ResponseWriter, r *http.Request){
	templator("home", homeView(), homeData(getCollection("event")), w)
}
func homeData(event []map[string]interface{})map[string]interface{}{
	var source string
	for _, row := range event{
		source += fmt.Sprintf("Name: %v, Club:%v, Date:%v, Time:%v<br>",row["name"],row["clubId"],row["datetime"],row["datetime"])
	}
	return map[string]interface{}{
		"Source": source,
		"PageName": "Calendar",
		"Menu": "Menu is not built yet",
	}
}
func homeView()string{
	return `<div id=topBar>
		<h1>{{ .PageName}}</h1>{{if .Menu}} {{XTC .Menu}}{{end}}
	</div>
	{{if .Source}}{{XTC .Source}}{{end}}
	<a href=organisers>Organisers</a>`
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
<style>body{background:#000}*{color:#777}`+ table_style +`</style>
</head>
<body>
	` + body + `
	{{if .Js}}<script src={{.Js}}></script>{{end}}
</body>
</html>`
}
func adminTemplate(body string)string{
	return `<!doctype html><html><head>
	{{if .Css}}<link rel=stylesheet href={{.Css}}>{{end}}
	{{if .Ico}}<link rel=icon type={{.IcoType}} href={{.Ico}}>{{end}}
	<title>EventBucket{{if .Title}} - {{.Title}}{{end}}</title>
	<style>body{background:#000}*{color:#777}`+ table_style +`</style>
</head><body>
	` + body + `
	{{if .Js}}<script src={{.Js}}></script>{{end}}
</body></html>`
}
func panel(title, source string)string{
	return `<h2>`+title+`</h2><div>`+source+`</div>`
}
func pane(title, body string)string{
	return`<div><h4>`+title+`</h4>`+body+`</div>`
}
