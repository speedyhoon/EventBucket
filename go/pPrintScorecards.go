package main

import (
	"bytes"
	"fmt"
	"image/png"
	"net/http"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/datamatrix"
	"github.com/boombuler/barcode/qr"
)

func barcode2D(w http.ResponseWriter, r *http.Request, parameters string) {
	var qrcode barcode.Barcode
	var err error

	switch parameters[0] {
	case 68, 100: //D for datamatrix
		qrcode, err = datamatrix.Encode(strings.ToUpper(parameters[1:]))
	default: //qrcode
		qrcode, err = qr.Encode(strings.ToUpper(parameters[1:]), qr.H, qr.Auto)
	}

	if err != nil {
		warn.Println(err)
		return
	}
	buf := new(bytes.Buffer)
	err = png.Encode(buf, qrcode)
	if err != nil {
		warn.Println(err)
		return
	}
	fmt.Fprintf(w, "%v", buf.String())
}

func printScorecards(w http.ResponseWriter, r *http.Request, parameters string) {
	//eventID/shooterID
	ids := strings.Split(parameters, "/")
	event, err := getEvent(ids[0])
	if err != nil {
		errorHandler(w, r, http.StatusNotFound, "event")
		return
	}

	var shooterID uint
	shooterID, err = stoU(ids[1])
	if err != nil || shooterID >= uint(len(event.Shooters)) {
		errorHandler(w, r, http.StatusNotFound, "shooter")
		return
	}

	if len(event.Ranges) < 1 {
		errorHandler(w, r, http.StatusNotFound, "range")
		return
	}
	templater(w, page{
		Title:  "Print Scorecards",
		Menu:   urlEvents,
		MenuID: event.ID,
		Data: map[string]interface{}{
			"Event":     event,
			"ShooterID": shooterID,
		},
	})
}
