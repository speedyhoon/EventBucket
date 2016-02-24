package main

import (
	"bytes"
	"fmt"
	"image/png"
	"net/http"
	"strings"

	"github.com/boombuler/barcode/qr"
)

func base64QrH(w http.ResponseWriter, r *http.Request, parameters string) {
	IDs := strings.Split(parameters, "/")
	if len(IDs) != 3 {
		warn.Println(fmt.Errorf("Three parameters were not in the url /eventID/shooterID/rangeID"))
		return
	}
	eventID := IDs[0]
	shooterID := IDs[1]
	rangeID := IDs[2]

	qrcode, err := qr.Encode(strings.ToUpper(eventID+"/"+shooterID+"/"+rangeID), qr.L, qr.Auto)
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
	fmt.Fprintf(w, string(buf.Bytes()))
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
		Data: M{
			"EventID":   eventID,
			"ShooterID": shooterID,
			"Ranges":    ranges,
			"FirstName": "Cam",
		},
	})
}
