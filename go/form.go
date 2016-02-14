package main

import (
	"regexp"
	"time"
)

type form struct {
	action uint8
	Fields []field
	Error  string
	expiry time.Time
}

type field struct {
	name, Error, Value string
	Required           bool
	Options            []option
	maxLen, minLen     uint
	min, max, step     float64
	AutoFocus          bool
	Checked            bool //only used by checkboxes
	regex              *regexp.Regexp
	internalValue      interface{}
	v8                 func([]string, field) (interface{}, string)
	defValue           func() []string
}

type option struct {
	Label, Value string
	Selected     bool
}

const (
	fieldMaxLen = 64

	clubNew              = 0
	clubDetails          = 1
	clubMoundNew         = 2
	eventNew             = 3
	eventDetails         = 4
	eventRangeNew        = 5
	eventAggNew          = 6
	eventShooterExisting = 7
	eventShooterNew      = 8
)

func getForm(id uint8) []field {
	switch id {
	case 0:
		return []field{{
			name:     "n",
			Required: true,
			maxLen:   0x40,
			min:      0,
			Checked:  false,
		}, {
			name:    "b",
			min:     0,
			Checked: false,
		}}
	}
	return []field{}
}
