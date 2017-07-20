package main

import (
	"bytes"
	"fmt"
	"image/png"
	"io"
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

func barcode2D(w io.Writer, code barcode.Barcode, err error) {
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

func printScorecards(w http.ResponseWriter, r *http.Request, eventID, shooterID string) {
	event, err := getEvent(eventID)
	if err != nil {
		errorHandler(w, r, "event")
		return
	}

	uShooterID, err := stoU(shooterID)
	if err != nil || uShooterID >= uint(len(event.Shooters)) {
		errorHandler(w, r, "shooter")
		return
	}

	if len(event.Ranges) < 1 {
		errorHandler(w, r, "range")
		return
	}
	templater(w, page{
		Title:  "Print Scorecards",
		Menu:   urlEvents,
		MenuID: event.ID,
		Data: map[string]interface{}{
			"Ranges":    event.Ranges,
			"Shooter":   event.Shooters[uShooterID],
			"EventID":   event.ID,
			"EventName": event.Name,
		},
	})
}
