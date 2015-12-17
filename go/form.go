package main

import "time"

type form struct {
	action, title string
	fields        []field
	inputs        []input
}

type field struct {
	name, Error, Value string
	Required           bool
	Options            []option
	maxLen, minLen     int
	min, max, step     int
	kind               interface{}
	defValue           func() []string
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

//250
var GlobalForms = []form{
	{
		title: "Insert Event",
		fields: []field{
			{
				name:     schemaClub,
				Required: true,
				//				autoFocus: true,
				maxLen: 50,
				minLen: 1,
				kind:   "",
				//				options: true,
			},
			{
				name:     schemaName,
				Required: true,
				maxLen:   50,
				minLen:   1,
				kind:     "",
				//				options: true,
			},
			{
				name:     schemaDate,
				defValue: defaultDate,
				minLen:   10,
				maxLen:   10,
				kind:     "",
			},
			{

				kind:     "",
				name:     schemaTime,
				defValue: defaultTime,
				minLen:   5,
				maxLen:   5,
			},
		},
	},
	/*{
		title:  "Insert Shooter",
		display: vertical,
		fields: []field{
			search{
				name:      schemaName,
				label:     "Event Name",
				autoFocus: true,
				required:  true,
				options: true,
			},
			search{
				name:     schemaClub,
				label:    "Club Name",
				required: true,
				autoFocus: true,
				maxLen:   50,
				options: true,
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
	{
		title: "Club Settings",
		fields: []field{
			{
				name:     schemaName,
				Required: true,
			},
			{
				name: schemaClubDefault,
			},
			{
				name: schemaClub,
				kind: "_id",
			},
		},
	},
}
