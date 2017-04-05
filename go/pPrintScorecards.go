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

func barcodeDM(w http.ResponseWriter, r *http.Request, parameters string) {
	dmCode, err := datamatrix.Encode(strings.ToUpper(parameters))
	barcode2D(w, dmCode, err)
}

func barcodeQR(w http.ResponseWriter, r *http.Request, parameters string) {
	qrCode, err := qr.Encode(strings.ToUpper(parameters), qr.H, qr.Auto)
	barcode2D(w, qrCode, err)
}

func barcode2D(w http.ResponseWriter, code barcode.Barcode, err error) {
	if err != nil {
		warn.Println(err)
		return
	}
	buf := new(bytes.Buffer)
	err = png.Encode(buf, code)
	if err != nil {
		warn.Println(err)
		return
	}
	fmt.Fprint(w, buf.String())
}

func printScorecards(w http.ResponseWriter, r *http.Request, eventID, shooterId string) {
	event, err := getEvent(eventID)
	if err != nil {
		errorHandler(w, r, http.StatusNotFound, "event")
		return
	}

	var shooterID uint
	shooterID, err = stoU(shooterId)
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
		JS:     []string{"print"},
		Data: map[string]interface{}{
			"Ranges":    event.Ranges,
			"Shooter":   event.Shooters[shooterID],
			"EventID":   event.ID,
			"EventName": event.Name,
		},
	})
}
