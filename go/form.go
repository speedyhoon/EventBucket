package main

import "time"

type form struct {
	title  string
	action int
	fields []field
	inputs []input
}

type field struct {
	name, Error, Value string
	Required           bool
	Options            []option
	maxLen, minLen     int
	min, max, step     int
	//	kind               interface{}
	internalValue interface{}
	kind          interface{}
	//	v8 string
	v8       func(string, field) (interface{}, string)
	defValue func() []string
}

type input struct {
	name, Error, Value string
	Required           bool
	Options            []option
	maxLen, minLen     int
	min, max, step     int
	kind               interface{}
}

type option struct {
	Label, Value string
	Selected     bool
}

func defaultDate() []string {
	return []string{time.Now().Format(formatYMD)}
}

func defaultTime() []string {
	return []string{time.Now().Format(formatTime)}
}

const (
	clubNew      = 0
	clubDetails  = 1
	clubMoundNew = 2
	eventNew     = 3
)

var GlobalForms = []form{
	clubNew: {
		fields: []field{
			{
				name:     schemaName,
				Required: true,
				v8:       isValidStr,
			},
			{
				name: schemaIsDefault,
				v8:   isValidBool,
			},
		},
	},
	clubDetails: {
		fields: []field{
			{
				name: schemaID,
				v8:   isValidStr,
			},
			{
				name:     schemaName,
				Required: true,
				v8:       isValidStr,
			},
			{name: schemaAddress, v8: isValidStr},
			{name: schemaTown, v8: isValidStr},
			{name: schemaPostCode, v8: isValidStr},
			{name: schemaLatitude, v8: isValidStr},
			{name: schemaLongitude, v8: isValidStr},
		},
	},
	clubMoundNew: {
		fields: []field{
			{
				//submit - Club ID
				name: schemaID,
				v8:   isValidStr,
			},
			{
				name: schemaName,
				v8:   isValidStr,
			},
			{
				name: schemaDistance,
				v8:   isValidInt,
			},
			{
				name: schemaUnit,
				v8:   isValidStr,
			},
		},
	},
	eventNew: {
		fields: []field{
			{
				name:     schemaClub,
				Required: true,
				//				autoFocus: true,
				maxLen: 64,
				minLen: 1,
				kind:   "",
				v8:     isValidStr,
				//				options: true,
			},
			{
				name:     schemaName,
				Required: true,
				maxLen:   64,
				minLen:   1,
				kind:     "",
				//				kind: isValidStr,
				v8: isValidStr,
				//				options: true,
			},
			{
				name:     schemaDate,
				defValue: defaultDate,
				minLen:   10,
				maxLen:   10,
				kind:     "",
				v8:       isValidStr,
			},
			{

				kind:     "",
				v8:       isValidStr,
				name:     schemaTime,
				defValue: defaultTime,
				minLen:   5,
				maxLen:   5,
			},
		},
	},
	/*{
		title:   "Insert Shooter",
		display: vertical,
		fields: []field{
			search{
				name:      schemaName,
				label:     "Event Name",
				autoFocus: true,
				required:  true,
				options:   true,
			},
			search{
				name:      schemaClub,
				label:     "Club Name",
				required:  true,
				autoFocus: true,
				maxLen:    50,
				options:   true,
			},
			dateTime{},
			dateTime{
				kind: "time",
			},
			checkbox{},
			hidden{},
			submit{
				name:  "eventId",
				value: "3",
				label: "Save",
			},
		},
	},*/
}
