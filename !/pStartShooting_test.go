package main

import (
	"fmt"
	"testing"
		"net/http"
	"io/ioutil"
	"net/url"
	"strings"




	//	"net"
	//	"net/http/httptest"
	//	"tls"
	//	"log"
	//	"time"
)

//func BenchmarkHello(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		fmt.Sprintf("hello")
//	}
//}

func TestUpdateShotScores2(t *testing.T) {
	//	t.Parallel()
	host := "http://localhost:81/"
	target := "eventInsert"
	eventName := "Tree"
	clubName := "Goose"
	dateValue := "12.12.2016"
	//	resp, err := http.PostForm(host, url.Values{"key": {"Value"}, "id": {"123"}})
	resp, err := http.PostForm(host+target, url.Values{"name": {eventName}, "club": {clubName}, "date": {dateValue}})
	defer resp.Body.Close()
	if err != nil {
		t.Error("unable to complete request")
	}
	body, err := ioutil.ReadAll(resp.Body)
	var html string
	if err != nil {
		t.Error("Something is wrong with the body")
	}
	html = fmt.Sprintf("%s", body)

	if !strings.Contains(html, "<h1>Event Settings - "+eventName+"</h1>") {
		t.Error("event page doesn't contain the event name")
	}

//	t.Logf("boo %v", html)
	resp, err = http.Get(host)
	defer resp.Body.Close()
	if err != nil{
		t.Error("unable to complete request")
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil{
		t.Error("Something is wrong with the body")
	}
	html = fmt.Sprintf("%s", body)


	if !strings.Contains(html, ">" + eventName + "</a>") {
		t.Error("home page doesn't contain the event name")
	}
	if !strings.Contains(html, ">" + clubName + "</a>") {
		t.Error("home page doesn't contain the club name")
	}
}

//type MockDd struct {}
//
//function (db MockDb) GetBacon() {
//	return "bacon"
//}

func TestHome(t *testing.T) {
//	mockDb := MockDb{}
	homeHandle := homeHandler(mockDb)
	req, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()
	homeHandle.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Home page didn't return %v", http.StatusOK)
	}
}








/*
type Server struct {
	URL string
	Listtener net.Listener
	TLS *tls.Config
	Config *http.Server
}
func NewServer(handler http.Handler) *Server {

}
func (*Server) Close() error {

}
*/
/*
func TestUpdateShotScores2(t *testing.T){
//			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
//				w.WriteHeader(404)
//				fmt.Fprintln(w, "hello")
//			}))
//			defer testServer.Close()
//			fmt.Println(testServer.URL)
//
//		//	var v float64
//		//	v = Average([]float64{1,2})
//		//	if v != 1.5 {
//		//		t.Error("Expected 1.5, got ", v)
//		//	}
//			v := 2 + 3
//			if v != 5{
//				t.Error("Expected 5, got ", v)
//			}


//	req, err := http.NewRequest("GET", "http://example.com/foo", nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer resp.Body.Close()


	url := "updateShotScores?eventid=Mk&rangeid=0&shooterid=138&shots=$%255555655555"
	eventName := fmt.Sprintf("%b", time.Now())
	clubName := "Booya"
//	resp, err := http.Post("http://localhost:81/")
	resp, err := http.Post("http://localhost:81/eventInsert", "application/x-www-form-urlencoded", "name=g&club=d&date=2015-06-25&time=22%3A17")

//	resp, err := http.PostForm(host + "eventInsert", url.Values{
//		"name": {eventName},
//		"club": {clubName},
//		"date": {},
//		"time": {},
//	})
//	resp, err := http.Post("http://localhost:81/eventInsert", "application/x-www-form-urlencoded", "name=GDF&club=87&date=2015-06-25&time=22%3A17")
//	resp, err := http.PostForm("http://localhost:81/eventInsert",
//		url.Values{"key": {"Value"}, "id": {"123"}})

	//	resp, err := http.PostForm(host + url, nil)
	defer resp.Body.Close()

	if err != nil {
		t.Error("failed to complete/retrieve new event body" + err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	html := fmt.Sprintf("%s", body)

	if !strings.Contains(html, ">" + eventName + "</a>") {
		t.Error("home page doesn't contain the event name")
	}
	if !strings.Contains(html, ">" + clubName + "</a>") {
		t.Error("home page doesn't contain the club name")
	}

//	event, err := getEvent(eventID)
//	if err != nil{
//		t.Error("failed to retrieve updated event")
//	}

}
*/
