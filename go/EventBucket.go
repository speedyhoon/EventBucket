// build+ debug

//go:generate goversioninfo -icon=favicon.ico

package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
)

const (
	//Logging & Database directory name
	subDir = `\EventBucket`

	//HTTP settings
	address = "http://localhost"
	dirRoot = "./"
	dirCSS  = "dirCSS"
	dirJS   = "dirJS"
	dirGzip = "dirGzip"
	dirPNG  = "dirPNG"
	dirGIF  = "dirGIF"
	robots  = "robots.txt"
	favicon = "favicon.ico"

	//Date formats
	formatGMT = "Mon, 02 Jan 2006 15:04:05 GMT"
)

var (
	//Database open connection
	db *bolt.DB

	debug bool

	//HTTP settings
	expiresTime, currentYear, portAddr string

	//Logging
	//Output destinations can be os.Stdout, os.Stderr, ioutil.Discard.
	//Flags can be log.Lshortfile|log.Ltime
	trace = log.New(os.Stdout, "TRACE: ", log.Lshortfile)
	info  = log.New(os.Stdout, "INFO:  ", log.Lshortfile)
	warn  = log.New(os.Stderr, "WARN:  ", log.Lshortfile)
)

func init() {
	go maintainExpiresTime()
	go maintainSessions()

	port := flag.Uint("port", 80, "Assign a differnet port number for the http server. Range: 0 through 65535.")
	flag.BoolVar(&debug, "debug", false, "Turn on debugging.")
	flag.Parse()

	if *port > math.MaxUint16 || *port < 0 {
		warn.Println("Port number must be between 0 and", math.MaxUint16)
		return
	}

	portAddr = fmt.Sprintf(":%v", *port)
	fullAddr := address
	if *port != 80 {
		fullAddr += portAddr
	}

	if !debug && exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", fullAddr).Start() != nil {
		warn.Print("Unable to open a web browser for " + fullAddr)
	}

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
	dbPath := filepath.Join(os.Getenv("ProgramData"), subDir)
	err := mkDir(dbPath)
	dbPath = filepath.Join(dbPath, "EventBucket.db")
	if err != nil {
		return
	}
	db, err = bolt.Open(dbPath, 0644, nil)
	if err != nil {
		warn.Println(err)
	}
	defer db.Close()

	pages()
	info.Print("Starting EventBucket HTTP server...")
	warn.Printf("ListenAndServe: %v", http.ListenAndServe(portAddr, nil))
	info.Println("EvenBucket server stopped.")
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
	expiresTime = time.Now().UTC().AddDate(1, 0, 0).Format(formatGMT)
	//w3.org says: "All HTTP date/time stamps MUST be represented in Greenwich Mean Time" under 3.3.1 Full Date //www.w3.org/Protocols/rfc2616/rfc2616-sec3.html
	masterTemplate.CurrentYear = time.Now().Format("2006")
}
