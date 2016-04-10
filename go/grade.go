package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Discipline separates different types of shooting so the number of shots & sighters can be easily changed while still using the same targets and Mark as another Discipline, e.g. Target rifles and Match rifles are vastly different disciplines but use the same scoring standard.
type Discipline struct {
	Name        string  `json:"name,omitempty"`
	ID          uint    `json:"id,omitempty"`
	QtySighters uint    `json:"sightersQty,omitempty"`
	QtyShots    uint    `json:"shotsQty,omitempty"`
	Grades      []Grade `json:"grades,omitempty"`
	Marking     Mark    `json:"marking,omitempty"`
	ShootOff    bool    `json:"shootOff,omitempty"`
}

// Mark is a group of settings associated with possible shooter scores on a target also known as "marking". Each type of target scoring standard can be specified by a Mark and be reused within several Disciplines.
type Mark struct {
	Buttons      string          `json:",buttons      omitempty"`
	Shots        map[string]Shot `json:",shots        omitempty"`
	DoCountBack2 bool            `json:",doCountBack2 omitempty"`
}

// Grade are subcategories of each discipline that shooters can be grouped together by similar skill levels.
type Grade struct {
	Name    string `json:"name,omitempty"` //Target A, F Class B, Match Reserve
	Abbr    string `json:"abbr,omitempty"` //Abbreviation of Name: A,B,C,FA,FB,FO,MO,MR
	ClassID uint   `json:"classID,omitempty"`
}

// Shot is exported
type Shot struct {
	Value      uint   `json:"value"`
	Center     uint   `json:"center,omitempty"`
	Shot       string `json:"shot,omitempty"`
	Sighter    string `json:"sighter,omitempty"`
	CountBack  string `json:"countBack,omitempty"`
	CountBack2 string `json:"countBack2,omitempty"`
}

var (
	globalDisciplines    []Discipline
	globalGrades         []Grade
	globalGradesDataList []option
)

func redoGlobals(disciplines []Discipline) {
	if len(disciplines) > 0 {
		globalDisciplines = disciplines
	} else {
		globalDisciplines = defaultGlobalDisciplines()
	}
	globalGrades = defaultGrades(globalDisciplines)
	globalGradesDataList = dataListGrades(globalGrades)
}

func defaultGrades(classes []Discipline) []Grade {
	var grades []Grade
	for classID, class := range classes {
		for _, grade := range class.Grades {
			grade.ClassID = uint(classID)
			grades = append(grades, grade)
		}
	}
	return grades
}

func dataListGrades(grades []Grade) []option {
	options := []option{{}}
	for id, grade := range grades {
		options = append(options, option{Value: fmt.Sprintf("%d", id), Label: grade.Name})
	}
	return options
}

func dataListGrades2(grades []Grade) []option {
	options := []option{}
	for id, grade := range grades {
		options = append(options, option{Value: fmt.Sprintf("%d", id), Label: grade.Name, Selected: true})
	}
	return options
}

func defaultGlobalDisciplines() []Discipline {
	XV5 := Mark{Buttons: "012345VX",
		DoCountBack2: true,
		Shots: map[string]Shot{
			"-": {Value: 0, Center: 0, CountBack: "0", CountBack2: "0", Shot: "-", Sighter: "-"},
			"0": {Value: 0, Center: 0, CountBack: "0", CountBack2: "0", Shot: "0", Sighter: "a"},
			"1": {Value: 1, Center: 0, CountBack: "1", CountBack2: "1", Shot: "1", Sighter: "b"},
			"2": {Value: 2, Center: 0, CountBack: "2", CountBack2: "2", Shot: "2", Sighter: "c"},
			"3": {Value: 3, Center: 0, CountBack: "3", CountBack2: "3", Shot: "3", Sighter: "d"},
			"4": {Value: 4, Center: 0, CountBack: "4", CountBack2: "4", Shot: "4", Sighter: "e"},
			"5": {Value: 5, Center: 0, CountBack: "5", CountBack2: "5", Shot: "5", Sighter: "f"},
			"V": {Value: 5, Center: 1, CountBack: "6", CountBack2: "6", Shot: "V", Sighter: "v"},
			"6": {Value: 5, Center: 1, CountBack: "6", CountBack2: "6", Shot: "V", Sighter: "g"},
			"X": {Value: 5, Center: 1, CountBack: "6", CountBack2: "7", Shot: "X", Sighter: "x"},
		}}
	return []Discipline{{
		ID:          0,
		Name:        "Target Rifle",
		QtySighters: 2,
		QtyShots:    10,
		Marking:     XV5,
		Grades: []Grade{{Abbr: "A", Name: "Target A"},
			{Abbr: "B", Name: "Target B"},
			{Abbr: "C", Name: "Target C"}},
	}, {
		ID:          1,
		Name:        "F Class",
		QtyShots:    12,
		QtySighters: 3,
		Marking: Mark{
			Buttons: "0123456X",
			Shots: map[string]Shot{
				"-": {Value: 0, Center: 0, CountBack: "0", Shot: "-", Sighter: "-"},
				"0": {Value: 0, Center: 0, CountBack: "0", Shot: "0", Sighter: "a"},
				"1": {Value: 1, Center: 0, CountBack: "1", Shot: "1", Sighter: "b"},
				"2": {Value: 2, Center: 0, CountBack: "2", Shot: "2", Sighter: "c"},
				"3": {Value: 3, Center: 0, CountBack: "3", Shot: "3", Sighter: "d"},
				"4": {Value: 4, Center: 0, CountBack: "4", Shot: "4", Sighter: "e"},
				"5": {Value: 5, Center: 0, CountBack: "5", Shot: "5", Sighter: "f"},
				"V": {Value: 6, Center: 0, CountBack: "6", Shot: "6", Sighter: "g"},
				"6": {Value: 6, Center: 0, CountBack: "6", Shot: "6", Sighter: "g"},
				"X": {Value: 6, Center: 1, CountBack: "7", Shot: "X", Sighter: "x"},
			}},
		Grades: []Grade{{Abbr: "FA", Name: "F Class A"},
			{Abbr: "FB", Name: "F Class B"},
			{Abbr: "FO", Name: "F Class Open"},
			{Abbr: "F/TR", Name: "F/TR"}},
	}, {
		ID:          2,
		Name:        "Match Rifle",
		QtySighters: 3,
		QtyShots:    15,
		Marking:     XV5,
		Grades: []Grade{{Abbr: "Open", Name: "Match Open"},
			{Abbr: "Reserve", Name: "Match Reserve"}},
	}, {
		ID:          3,
		Name:        "Service Rifle",
		QtySighters: 1,
		QtyShots:    8,
		Marking: Mark{
			Buttons: "012345V",
			Shots: map[string]Shot{
				"-": {Value: 0, Center: 0, CountBack: "0", Shot: "-", Sighter: "-"},
				"0": {Value: 0, Center: 0, CountBack: "0", Shot: "0", Sighter: "a"},
				"1": {Value: 1, Center: 0, CountBack: "1", Shot: "1", Sighter: "b"},
				"2": {Value: 2, Center: 0, CountBack: "2", Shot: "2", Sighter: "c"},
				"3": {Value: 3, Center: 0, CountBack: "3", Shot: "3", Sighter: "d"},
				"4": {Value: 4, Center: 0, CountBack: "4", Shot: "4", Sighter: "e"},
				"5": {Value: 5, Center: 0, CountBack: "5", Shot: "5", Sighter: "f"},
				"V": {Value: 5, Center: 1, CountBack: "6", Shot: "V", Sighter: "v"},
				"6": {Value: 5, Center: 1, CountBack: "6", Shot: "V", Sighter: "v"},
				"X": {Value: 5, Center: 1, CountBack: "6", Shot: "V", Sighter: "v"},
			}},
		Grades: []Grade{{Abbr: "303", Name: "303 Rifle"}},
	}}
}

func loadGrades(filePath string) error {
	if filePath == "" {
		return nil
	}
	contents := readFile(filePath)
	//If file is empty, try to write a new JSON file to storage.
	if len(contents) < 1 {
		warn.Println("empty file contents")
		//Generate JSON from globalDisciplines
		src, err := json.MarshalIndent(globalDisciplines, "", "\t")
		//Output marshal errors
		if err != nil {
			warn.Println("error:", err)
		}
		writeFile(filePath, src)
		info.Println("Created grades settings JSON file:", filePath)
		//Return an error because EventBucket was unable to load the original settings file specified.
		return fmt.Errorf("Unable to load settings file: %v", filePath)
	}
	var disciplines []Discipline
	err := json.Unmarshal(contents, &disciplines)
	if err != nil {
		//Unable to unmarshal settings from JSON file.
		warn.Println("error:", err)
		return fmt.Errorf("Error: %v, File: %v", err, filePath)
	}
	info.Println("Loaded grade settings from:", filePath)
	redoGlobals(disciplines)
	return nil
}

func readFile(filename string) []byte {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		warn.Println(err)
	}
	return src
}

func writeFile(destination string, src []byte) {
	err := ioutil.WriteFile(destination, src, 0777)
	if err != nil {
		warn.Printf("\nUnable to write to file %v -- %v", destination, err)
	}
}
