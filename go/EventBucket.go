// build+ debug

package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	//"os/signal"
	"path/filepath"
	"time"
)

const (
	debug = true

	//Logging & Database directory name
	subDir = "/EventBucket"

	//HTTP settings
	address = "http://localhost"
	port    = "80"
	dirRoot = "./"
	//	dirCSS  = "/css/"
	//dirJS   = "dirJS"
	dirPNG  = "dirPNG"
	robots  = "robots.txt"
	favicon = "favicon.ico"

//	sitemap = "sitemap.xml"

//dirHTML = "/html/"
//dirSVG  = "/svg/"
//dirWOF2 = "/woff2/"
)

var (
	//HTTP settings
	expiresTime, currentYear string

	//Logging
	tempPath    = os.Getenv("temp") + subDir
	logFileName = filepath.Join(tempPath, time.Now().Format("20060102")+".log")
	//Logging destinations are os.Stdout, os.Stderr, ioutil.Discard
	trace = log.New(os.Stdout, "TRACE: ", log.Lshortfile)
	info  = log.New(os.Stdout, "INFO:  ", log.Lshortfile|log.Ltime)
	warn  = log.New(os.Stderr, "WARN:  ", log.Lshortfile|log.Ltime)

	//EventBucket database
	databasePath     = os.Getenv("ProgramData") + subDir
	regexWeakEventId = regexp.MustCompile(`^[[:alnum:]]+$`)
	regexEventId     = regexp.MustCompile(`^[a-z0-9]+$`)
)

type M map[string]interface{}

func init() {
	startLogging()
	go maintainExpiresTime()
	go mkDir(databasePath)
	setExpiresTime()

	//Display message during shutdown.
	/*osChan := make(chan os.Signal, 1)
	signal.Notify(osChan, os.Interrupt)
	go func(){
		for _ = range osChan{
			fmt.Println("Shutting down EventBucket")
			<- osChan
			signal.Stop(osChan)
			break
		}
		os.Exit(0)
	}()*/

	/*my_temp := map[string]Field2{
		"eventId": submit{
			//				name:      ,
			autoFocus: true,
			Value:     "9",
			label:     "Button!",
		},
		"shooterId": selectBox{
			//				name:     ,
			required: true,
			options: []option{
				{Label: "5"},
				{Label: "6"},
			},
			label: "selectBox",
		},

		"moundId": multiSelect{
			//				name:     ,
			required: true,
			options: []option{
				{Label: "500"},
				{Label: "600"},
			},
			size:  4,
			label: "multiSelect",
		},
		"grren": hidden{
			//				name:  ,
			Value: "3t3",
		},
	}

	myForm := form2{
		action: "startshooting",
		title:  "Start shooting FOrm",
		fields: my_temp,
	}
	fmt.Println(myForm.html(""))*/
}

/*
func tempers(eventID string) form3 {
	return form3{
		action: "/EventRangeInsert",
		title:  "Add Range",
		fields: []Field2{
			textbox{
				name:      "name",
				label:     "Range Name",
				error:     "whoops this seems to be an unexpected error :(",
				autoFocus: true,
				required:  true,
			}, submit{
				label: "Create Range",
				name:  "eventid",
				Value: eventID,
			},
		},
	}
}*/

func main() {
	serveFile(favicon)
	serveFile(robots)
	serveDir(dirPNG)
	//BUG any url breaks when appending "&*((&*%"
	get404(urlHome, home)
	getParameters(urlEvent, event, regexEventId, regexWeakEventId)
	getRedirectPermanent(urlClubs, clubs)
	getRedirectPermanent(urlAbout, about)
	getRedirectPermanent(urlArchive, eventArchive)
	getRedirectPermanent(urlShooters, shooters)
	getRedirectPermanent(urlLicence, licence)
	getRedirectPermanent(urlEvents, events)
	getRedirectPermanent("/all", all)

	http.HandleFunc("/0", insertEvent)

	//	fmt.Println(tempers("4"))
	//	fmt.Printf("\n\n\n\n%#v", tempers("4"))
	//	fmt.Printf("\n\n\n\n%+v", tempers("4"))
	//	fmt.Printf("\n\n\n\n%T", tempers("4"))
	//	my_temp := tempers("4")
	//	for i, t := range my_temp.fields {
	//		fmt.Println(i, t, "\n\n")
	//	}

	if !debug && exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", address).Start() != nil {
		warn.Print("Unable to open a web browser for " + address)
	}
	info.Print("Starting EventBucket server...")
	warn.Printf("ListenAndServe: %v", http.ListenAndServe(":"+port, nil))
	info.Println("EvenBucket stopped.")
}

//Setup logging into temp directory
func startLogging() {
	err := mkDir(tempPath)
	if err == nil {
		var f *os.File
		f, err = os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err == nil {
			trace.SetOutput(f)
			info.SetOutput(f)
			warn.SetOutput(f)
		}
		//Not calling defer to close file pointer "f" because we need it open after init has finished. Not witnessing any adverse affects keeping the file open.
	}
}

// Attempt to create the path supplied if it doesn't exist.
func mkDir(path string) error {
	stat, err := os.Stat(path)
	if err != nil || !stat.IsDir() {
		err = os.Mkdir(path, os.ModeDir)
		if err != nil {
			warn.Printf("Unable to create directory %v %v", path, err)
		}
	}
	return err
}

//Update the expires http header time, every 15 minutes rather than recalculating it on every http request.
func maintainExpiresTime() {
	ticker := time.NewTicker(time.Minute * 15)
	for range ticker.C {
		//Can't directly change global variables in a go routine, so call an external function.
		setExpiresTime()
	}
}

//Set expiry date 1 year, 0 months & 0 days in the future.
func setExpiresTime() {
	//Date format is the same as Go`s time.RFC1123 but uses "GMT" timezone instead of "UTC" time standard.
	expiresTime = time.Now().UTC().AddDate(1, 0, 0).Format("Mon, 02 Jan 2006 15:04:05 GMT")
	//w3.org says: "All HTTP date/time stamps MUST be represented in Greenwich Mean Time" under 3.3.1 Full Date //www.w3.org/Protocols/rfc2616/rfc2616-sec3.html
	masterTemplate.CurrentYear = time.Now().Format("2006")
}
