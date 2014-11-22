package main

import (
	"fmt"
	"sort"
)

type dataSlice []*data

type data struct {
	count int64
	size  int64
}
type sortIt []shoota
type shoota struct {
	Class	string
	Grade	string
	Total int
	Centa	int
}

// Len is part of sort.Interface.
func (d dataSlice) Len() int {
	return len(d)
}

// Swap is part of sort.Interface.
func (d dataSlice) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

// Less is part of sort.Interface. We use count as the value to sort by
func (d dataSlice) Less(i, j int) bool {
	return d[i].count < d[j].count
}
//////////////////////////////////////////
//////////////////////////////////////////
//////////////////////////////////////////
//////////////////////////////////////////
// Len is part of sort.Interface.
func (d sortIt) Len() int {
	return len(d)
}

// Swap is part of sort.Interface.
func (d sortIt) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

// Less is part of sort.Interface. We use count as the value to sort by
func (d sortIt) Less(i, j int) bool {
	return d[i].Total < d[j].Total
}

func main() {
	event := Event{
		Shooters:	map[string]Shooter{
			"0": Shooter{
//				Id:	"fd",
//				FirstName: "Cam",
//				Surname: "W",
//				MiddleName: "F",
//				Age: "",
//				Club: 1,
				Class: "target",
				Grade: "B",
				Score: map[string]Score{
					"0":	Score{
						Total: 40,
						Centers: 5,
					},
					"1":	Score{
						Total: 40,
						Centers: 4,
					},
					"2":	Score{
						Total: 40,
						Centers: 3,
					},
					"3":	Score{
						Total: 40,
						Centers: 6,
					},
					"4":	Score{
						Total: 160,
						Centers: 18,
					},
				},
			},
			"1": Shooter{
//				Id:	"sa",
//				FirstName: "Leesa",
//				Surname: "N",
//				MiddleName: "M",
//				Age: "U21",
//				Club: 2,
				Class: "target",
				Grade: "B",
				Score: map[string]Score{
					"0":	Score{
						Total: 40,
						Centers: 5,
					},
					"1":	Score{
						Total: 40,
						Centers: 4,
					},
					"2":	Score{
						Total: 40,
						Centers: 3,
					},
					"3":	Score{
						Total: 40,
						Centers: 7,
					},
					"4":	Score{
						Total: 160,
						Centers: 19,
					},
				},
			},
			"2": Shooter{
//				Id:	"sa",
//				FirstName: "Leesa",
//				Surname: "N",
//				MiddleName: "M",
//				Age: "U21",
//				Club: 2,
				Class: "target",
				Grade: "B",
				Score: map[string]Score{
					"0":	Score{
						Total: 40,
						Centers: 5,
					},
					"1":	Score{
						Total: 40,
						Centers: 4,
					},
					"2":	Score{
						Total: 40,
						Centers: 3,
					},
					"3":	Score{
						Total: 40,
						Centers: 7,
					},
					"4":	Score{
						Total: 160,
						Centers: 10,
					},
				},
			},
			"3": Shooter{
//				Id:	"sa",
//				FirstName: "Leesa",
//				Surname: "N",
//				MiddleName: "M",
//				Age: "U21",
//				Club: 2,
				Class: "target",
				Grade: "B",
				Score: map[string]Score{
					"0":	Score{
						Total: 40,
						Centers: 5,
					},
					"1":	Score{
						Total: 40,
						Centers: 4,
					},
					"2":	Score{
						Total: 40,
						Centers: 3,
					},
					"3":	Score{
						Total: 40,
						Centers: 7,
					},
					"4":	Score{
						Total: 159,
						Centers: 60,
					},
				},
			},
		},
	}





	M := map[string]*data {
		"x": {0, 0},
		"y": {2, 9},
		"z": {1, 7},
	}

	s := make(dataSlice, 0, len(M))

	for _, d := range M {
		s = append(s, d)
	}

	sort.Sort(s)

	for _, d := range s {
		fmt.Printf("%+v\n", *d)
	}



	sortByRange := "4"
	shooters := make(sortIt, 0, 2)
	for _, shooterList := range event.Shooters {
		temp := shoota{
			Class: shooterList.Class,
			Grade: shooterList.Grade,
			Total: shooterList.Score[sortByRange].Total,
			Centa: shooterList.Score[sortByRange].Centers,
		}
		shooters = append(shooters, temp)
	}
	fmt.Printf("%v", shooters)
	fmt.Println("--")
	sort.Sort(shooters)
	fmt.Printf("%v", shooters)
}

type Score struct {
	Total			int		`bson:"t,omitempty"`
	Centers		int		`bson:"c,omitempty"`
	Shots			string	`bson:"s,omitempty"`
	Countback	uint64	`bson:"b,omitempty"`
	Xs				int		`bson:"x,omitempty"`
}
type Shooter struct {
//	Id		bson.ObjectId		`bson:"_id,omitempty"`
	FirstName	string		`bson:"f,omitempty"`
	Surname		string		`bson:"s,omitempty"`
	MiddleName	string		`bson:"m,omitempty"`
	Age			string		`bson:"a,omitempty"`
	Club			int			`bson:"C,omitempty"`
	Class			string		`bson:"c,omitempty"`
	Grade			string		`bson:"g,omitempty"`
	Score	map[string]Score	`bson:"r,omitempty"`
}
type Event struct {
	//		Id			bson.ObjectId			`bson:"_id,omitempty"`
	//		DateTime	string					`bson:"d,omitempty"`
	//		AutoInc	AutoInc					`bson:"U,omitempty"`
	//		Range		map[string]Range		`bson:"R,omitempty"`
	Shooters	map[string]Shooter	`bson:"S,omitempty"`
	//		Club		int						`bson:"C,omitempty"`
}
