package main

type form struct {
	action, title string
	fields        []field
}

type field struct {
	name, error, value string
	required           bool
	options            []option
	maxLen, minLen     int
	min, max, step     int
	kind               interface{}
}

type input struct {
	name, Error, Value string
	required           bool
	Options            []option
	maxLen, minLen     int
	min, max, step     int
	kind               interface{}
}

type option struct {
	Label, Value string
	Selected     bool
}
