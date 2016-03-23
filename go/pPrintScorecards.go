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
	var shooter EventShooter
	var ranges []Range
	event, err := getEvent(ids[0])

	//getParameters() checks that parameters contains a "/" otherwise it would fail.
	shooterID, sErr := b36tou(ids[1])
	if sErr == nil && int(shooterID) < len(event.Shooters) {
		shooter = event.Shooters[shooterID]
	}

	if err == nil {
		ranges = event.Ranges
		//Display the error converting shooter ID to a number (if any) when event was retrieved OK.
		err = sErr
	}
	templater(w, page{
		Title:  "Print Scorecards",
		Menu:   urlEvents,
		MenuID: event.ID,
		Error:  err,
		Data: map[string]interface{}{
			"EventID":   event.ID,
			"EventName": event.Name,
			"Ranges":    ranges,
			"Shooter":   shooter,
		},
	})
}
