package main

import "time"

type form struct {
	action int
	Fields []field
	Error  string
}

type field struct {
	name, Error, Value string
	Required           bool
	Options            []option
	maxLen, minLen     int
	min, max, step     int
	AutoFocus          bool
	//	kind               interface{}
	internalValue interface{}
	kind          interface{}
	//	v8 string
	v8       func(string, field) (interface{}, string)
	v9       func([]string, field) (interface{}, string)
	defValue func() []string
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

var GlobalForms = [][]field{
	clubNew: {
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
	clubDetails: {
		{
			name:     schemaName,
			Required: true,
			v8:       isValidStr,
		},
		{name: schemaAddress, v8: isValidStr},
		{name: schemaTown, v8: isValidStr},
		{name: schemaPostcode, v8: isValidStr},
		{name: schemaLatitude, v8: isValidFloat64,
			min: -90,
			max: 90,
		},
		{name: schemaLongitude, v8: isValidFloat64,
			min: -180,
			max: 180,
		},
		{name: schemaClub, v8: isValidID},
	},
	clubMoundNew: {
		{
			name: schemaDistance,
			//			v8:       isValidInt,
			v8:       isValidUint64,
			Required: true,
		},
		{
			name:     schemaUnit,
			v8:       isValidStr,
			Required: true,
		},
		{
			//submit - Club ID
			name:     schemaClub,
			v8:       isValidID,
			Required: true,
		},
	},
	eventNew: {
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
	eventDetails: {},
	eventRangeNew: {
		{
			name:     schemaName,
			Required: true,
			v8:       isValidStr,
		}, {
			name:     schemaEvent,
			Required: true,
			v8:       isValidID,
		},
	},
	eventAggNew: {
		{
			name:     schemaName,
			Required: true,
			v8:       isValidStr,
		}, {
			name:     schemaRange,
			Required: true,
			v9:       isValidRangeIDs,
		}, {
			name:     schemaEvent,
			Required: true,
			v8:       isValidID,
		},
	},
	eventShooterNew: {
		{
			name:     schemaFirstName,
			Required: true,
			v8:       isValidStr,
		}, {
			name:     schemaSurname,
			Required: true,
			v8:       isValidStr,
		}, {
			name:     schemaClub,
			Required: true,
			v8:       isValidStr,
		}, {
			name:     schemaGrade,
			Required: true,
			//			v8:       isValidGrade,
			v8:      isValidUint64,
			min:     1,
			max:     len(dataListGrades()),
			step:    1,
			Options: dataListGrades(),
		}, {
			name: schemaAgeGroup,
			//			v8:   isValidAgeGroup,
			v8:      isValidUint64,
			min:     1,
			max:     len(dataListAgeGroup()),
			step:    1,
			Options: dataListAgeGroup(),
		}, {
			name:     schemaEvent,
			Required: true,
			v8:       isValidID,
		},
	},
	eventShooterExisting: {
		{
			name:     schemaShooter,
			Required: true,
			v8:       isValidID,
		}, {
			name:     schemaGrade,
			Required: true,
			v8:       isValidUint64,
			//			v8:       isValidGrade,
			min:     1,
			max:     len(dataListGrades()) - 1,
			step:    1,
			Options: dataListGrades(),
		}, {
			name: schemaAgeGroup,
			//			v8:   isValidAgeGroup,
			v8:      isValidUint64,
			min:     1,
			max:     len(dataListAgeGroup()) - 1,
			step:    1,
			Options: dataListAgeGroup(),
		}, {
			name:     schemaEvent,
			Required: true,
			v8:       isValidID,
		},
	},
}
