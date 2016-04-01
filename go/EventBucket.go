//go:generate goversioninfo -icon=favicon.ico

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
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
)

var (
	//Database connection.
	db *bolt.DB

	//Command line flags.
	portAddr string
	debug    bool

	//Used for every HTTP request with cache headers set.
	cacheExpires string

	//Logging
	trace = log.New(ioutil.Discard, "TRACE: ", log.Lshortfile) //Flags can be log.Lshortfile|log.Ltime
	info  = log.New(os.Stdout, "INFO:  ", log.Lshortfile)
	warn  = log.New(os.Stderr, "WARN:  ", log.Lshortfile)
)

func init() {
	go maintainExpiresTime()
	go maintainSessions()

	//Add support for changing the port number as a command line flag
	port := flag.Uint("port", 80, "Assign a differnet port number for the http server. Range: 0 through 65535.")
	flag.BoolVar(&debug, "debug", false, "Turn on debugging and turn off HTML file caching.")
	flag.Parse()

	if *port >= math.MaxUint16 || *port < 2 {
		warn.Printf("Port number must be between 2 and %d. Default port number is 80.", math.MaxUint16-1)
		os.Exit(-3)
	}

	portAddr = fmt.Sprintf(":%v", *port)
	fullAddr := "http://localhost"
	if *port != 80 {
		fullAddr += portAddr
	}

	if debug {
		trace.SetOutput(os.Stdout)
	} else if exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", fullAddr).Start() != nil {
		warn.Println("Unable to open a web browser for", fullAddr)
	}

	setExpiresTime()
}

func main() {
	//Database save location
	dbPath := filepath.Join(os.Getenv("ProgramData"), `\EventBucket`)
	err := mkDir(dbPath)
	dbPath = filepath.Join(dbPath, "EventBucket.db")
	if err != nil {
		return
	}

	//Open database connection
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
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
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
	cacheExpires = time.Now().UTC().AddDate(1, 0, 0).Format(formatGMT)
	//w3.org: "All HTTP date/time stamps MUST be represented in Greenwich Mean Time" under 3.3.1 Full Date //www.w3.org/Protocols/rfc2616/rfc2616-sec3.html
	masterTemplate.CurrentYear = time.Now().Format("2006")
}
