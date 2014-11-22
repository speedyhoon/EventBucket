package main
import (
	"net/http"
//	"mgo"


	"Eventbucket"


//	"html/template"
//	"fmt"
//	"net/url"

	//	"mgo/bson"
	//	"reflect"
	//	"strings"
)
var conn *mgo.Database//global variables that can't be a constant
func main(){
	conn = DB()
	//	http.ResponseWriter.Header().Set("Expires", "Thu, 4 Apr 2013 20:00:00 GMT")



	http.HandleFunc("/", Eventbucket.home)
//	http.HandleFunc("/clubs", clubs)
//	http.HandleFunc("/startShooting", startShooting)
//	http.HandleFunc("/clubInsert", clubInsert)
//	http.HandleFunc("/eventInsert", eventInsert)
//	http.HandleFunc("/organisers", organisers)
//	http.HandleFunc("/organiser", organisers)
//	http.HandleFunc("/eventShow", eventShow)
//	http.HandleFunc("/eventSetup", eventSetup)
	http.ListenAndServe(":80", nil)
}
















