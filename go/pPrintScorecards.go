package main

import (
	"bytes"
	"fmt"
	"image/png"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/datamatrix"
	"github.com/boombuler/barcode/qr"
)

func barcodeDM(w http.ResponseWriter, _ *http.Request, parameters string) {
	dmCode, err := datamatrix.Encode(strings.ToUpper(parameters))
	barcode2D(w, dmCode, err)
}

func barcodeQR(w http.ResponseWriter, _ *http.Request, parameters string) {
	qrCode, err := qr.Encode(strings.ToUpper(parameters), qr.H, qr.Auto)
	barcode2D(w, qrCode, err)
}

func barcode2D(w io.Writer, code barcode.Barcode, err error) {
	if err != nil {
		log.Println(err)
		return
	}
	buf := new(bytes.Buffer)
	err = png.Encode(buf, code)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = fmt.Fprint(w, buf.String())
	if err != nil {
		log.Println(err)
	}
}

func printScorecards(w http.ResponseWriter, r *http.Request, event Event, shooterID sID) {
	if len(event.Ranges) < 1 {
		errorHandler(w, r, "range")
		return
	}
	render(w, page{
		Title:  "Print Scorecards",
		Menu:   urlEvents,
		MenuID: event.ID,
		Data: map[string]interface{}{
			"Event":   event,
			"Shooter": event.Shooters[shooterID],
		},
	})
}
