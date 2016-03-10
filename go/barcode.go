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

func base64QrH(w http.ResponseWriter, r *http.Request, parameters string) {
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
	IDs := strings.Split(parameters, "/")
	if len(IDs) != 2 {
		return
	}
	eventID := IDs[0]
	shooterID := IDs[1]
	ranges := []uint64{1, 2, 3, 4, 5, 6}

	templater(w, page{
		Title:  "Print Scorecards",
		Menu:   urlEvents,
		MenuID: eventID,
		Data: map[string]interface{}{
			"EventID":   eventID,
			"ShooterID": shooterID,
			"Ranges":    ranges,
			"FirstName": "Cam",
		},
	})
}
