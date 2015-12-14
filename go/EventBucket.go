// build+ debug

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	//"os/signal"
	"path/filepath"
	"time"
)

const (
	//Logging & Database directory name
	subDir = "/EventBucket"

	//HTTP settings
	address = "http://localhost"
	dirPNG  = "dirPNG"
	dirJS   = "dirJS"
	dirCSS  = "dirCSS"
	robots  = "robots.txt"
	favicon = "favicon.ico"

	schemaClub        = "schemaClub"
	schemaEvent       = "schemaEvent"
	schemaMound       = "schemaMound"
	schemaRange       = "schemaRange"
	schemaShooter     = "schemaShooter"
	schemaAutoInc     = "schemaAutoInc"
	schemaID          = "schemaID"
	schemaAddress     = "schemaAddress"
	schemaClubDefault = "schemaClubDefault"
	schemaClose       = "schemaClose"
	schemaDate        = "schemaDate"

	schemaGrade          = "schemaGrade"
	schemaIsPrizeMeet    = "schemaIsPrizeMeet"
	schemaLongName       = "schemaLongName"
	schemaName           = "schemaName"
	schemaSortScoreboard = "schemaSortScoreboard"
	schemaPostcode       = "schemaPostcode"
	schemaSort           = "schemaSort"
	schemaTime           = "schemaTime"
	schemaURL            = "schemaURL"

	schemaTown      = "schemaTown"
	schemaLatitude  = "schemaLatitude"
	schemaLongitude = "schemaLongitude"
)

var (
	debug bool

	//HTTP settings
	expiresTime, currentYear, portAddr string

	//Logging
	tempPath    = os.Getenv("temp") + subDir
	logFileName = filepath.Join(tempPath, time.Now().Format("20060102")+".log")
	//Logging destinations are os.Stdout, os.Stderr, ioutil.Discard
	trace = log.New(ioutil.Discard, "TRACE: ", log.Lshortfile)
	info  = log.New(os.Stdout, "INFO:  ", log.Lshortfile|log.Ltime)
	warn  = log.New(os.Stderr, "WARN:  ", log.Lshortfile|log.Ltime)

	//EventBucket database
	databasePath = os.Getenv("ProgramData") + subDir

	//URL validation matching
	regexEventId = regexp.MustCompile(`^[a-z0-9]+$`)
)

type M map[string]interface{}

func init() {
	port := flag.Uint("port", 80, "Assign a differnet port number for the http server. Range: 0 through 65535.")
	flag.BoolVar(&debug, "debug", false, "Turn on debugging.")
	flag.Parse()

	if *port > math.MaxUint16 || *port < 0 {
		info.Println("Port number must be between 0 and", math.MaxUint16)
		return
	}

	portAddr = fmt.Sprintf(":%v", *port)
	fullAddr := address
	if *port != 80 {
		fullAddr += portAddr
	}
	info.Println(fullAddr)

	if !debug && exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", fullAddr).Start() != nil {
		warn.Print("Unable to open a web browser for " + fullAddr)
	}

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
}

func main() {
	serveFile(favicon)
	serveFile(robots)
	serveDir(dirCSS, false)
	serveDir(dirJS, false)
	serveDir(dirPNG, false)
	//BUG any url breaks when appending "&*((&*%"
	get404(urlHome, home)
	getRedirectPermanent(urlClubs, clubs)
	getRedirectPermanent(urlAbout, about)
	getRedirectPermanent(urlArchive, eventArchive)
	getRedirectPermanent(urlShooters, shooters)
	getRedirectPermanent(urlLicence, licence)
	getRedirectPermanent(urlEvents, events)
	getRedirectPermanent("/all", all)
	getParameters(urlEvent, event, regexEventId)
	http.HandleFunc("/0", insertEvent)

	info.Print("Starting EventBucket HTTP server...")
	warn.Printf("ListenAndServe: %v", http.ListenAndServe(portAddr, nil))
	info.Println("EvenBucket stopped.")
}

//Setup logging into temp directory
func startLogging() {
	err := mkDir(tempPath)
	if err == nil {
		var f *os.File
		f, err = os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err == nil {
			if debug {
				trace.SetOutput(f)
			}
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
