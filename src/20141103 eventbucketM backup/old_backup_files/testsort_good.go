package main

import (
	"fmt"
	"sort"
)

type Ranges struct {
	total int
	centers int
	cbs string
}

// A Change is a record of source code changes, recording user, language, and delta size.
type Change struct {
	clas string
	grade string
	ranges []Ranges
	first string
	surname string
}

type lessFunc func(p1, p2 *Change) bool

// multiSorter implements the Sort interface, sorting the changes within.
type multiSorter struct {
	changes []Change
	less    []lessFunc
}

// Sort sorts the argument slice according to the less functions passed to OrderedBy.
func (ms *multiSorter) Sort(changes []Change) {
	ms.changes = changes
	sort.Sort(ms)
}

// OrderedBy returns a Sorter that sorts using the less functions, in order.
// Call its Sort method to sort the data.
func OrderedBy(less ...lessFunc) *multiSorter {
	return &multiSorter{
		less: less,
	}
}

// Len is part of sort.Interface.
func (ms *multiSorter) Len() int {
	return len(ms.changes)
}

// Swap is part of sort.Interface.
func (ms *multiSorter) Swap(i, j int) {
	ms.changes[i], ms.changes[j] = ms.changes[j], ms.changes[i]
}

// Less is part of sort.Interface. It is implemented by looping along the
// less functions until it finds a comparison that is either Less or
// !Less. Note that it can call the less functions twice per call. We
// could change the functions to return -1, 0, 1 and reduce the
// number of calls for greater efficiency: an exercise for the reader.
func (ms *multiSorter) Less(i, j int) bool {
	p, q := &ms.changes[i], &ms.changes[j]
	// Try all but the last comparison.
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		switch {
		case less(p, q):
			// p < q, so we have a decision.
			return true
		case less(q, p):
			// p > q, so we have a decision.
			return false
		}
		// p == q; try the next comparison.
	}
	// All comparisons to here said "equal", so just return whatever
	// the final comparison reports.
	return ms.less[k](p, q)
}

// ExampleMultiKeys demonstrates a technique for sorting a struct type using different
// sets of multiple fields in the comparison. We chain together "Less" functions, each of
// which compares a single field.
func main() {
	// Closures that order the Change structure.
	total := func(c1, c2 *Change) bool {
		return c1.ranges[0].total > c2.ranges[0].total
	}
	centa := func(c1, c2 *Change) bool {
		return c1.ranges[0].centers > c2.ranges[0].centers
	}
	cb := func(c1, c2 *Change) bool {
		return c1.ranges[0].cbs > c2.ranges[0].cbs
	}
	clas := func(c1, c2 *Change) bool {
		return c1.clas < c2.clas
	}
	grade := func(c1, c2 *Change) bool {
		return c1.grade < c2.grade
	}


	var changes = []Change{}

	sortByRange := "4"
	for _, shooterList := range event.Shooters {
		temp := Change{
			clas: shooterList.Class,
			grade: shooterList.Grade,
			ranges: []Ranges{{
				total: shooterList.Score[sortByRange].Total,
				centers: shooterList.Score[sortByRange].Centers,
				cbs: shooterList.Score[sortByRange].Countback,
			}},
			first: shooterList.FirstName,
			surname: shooterList.MiddleName,
		}
		changes = append(changes, temp)
	}

	OrderedBy(clas, grade, total, centa, cb).Sort(changes)
	fmt.Println("By language,<lines,user:")
	for _, item := range changes{
		fmt.Println(fmt.Sprintf("%v",item))
	}
}

var event = Event{
	Shooters:	map[string]Shooter{
		"0": Shooter{
			FirstName: "Cam",
			Surname: "W",
			MiddleName: "F",
			Age: "",
			Club: 1,
			Class: "target",
			Grade: "B",
			Score: map[string]Score{
				"4":	Score{
					Total: 160,
					Centers: 18,
					Countback: "1234567890",
				},
			},
		},
		"1": Shooter{
			FirstName: "Leesa",
			Surname: "N",
			MiddleName: "M",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "B",
			Score: map[string]Score{
				"4":	Score{
					Total: 160,
					Centers: 19,
					Countback: "1234567890",
				},
			},
		},
		"2": Shooter{
			FirstName: "Sam",
			Surname: "N",
			MiddleName: "Smith",
			Age: "",
			Club: 2,
			Class: "fclass",
			Grade: "FB",
			Score: map[string]Score{
				"4":	Score{
					Total: 200,
					Centers: 25,
					Countback: "1234567890",
				},
			},
		},
		"3": Shooter{
			FirstName: "Larry",
			Surname: "N",
			MiddleName: "Perkins",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "A",
			Score: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 13,
					Countback: "1234567890",
				},
			},
		},
		"4": Shooter{
			FirstName: "Peter",
			Surname: "N",
			MiddleName: "Brock",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "A",
			Score: map[string]Score{
				"4":	Score{
					Total: 160,
					Centers: 14,
					Countback: "1234567890",
				},
			},
		},
		"5": Shooter{
			FirstName: "Darrel",
			Surname: "N",
			MiddleName: "Nepia",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "B",
			Score: map[string]Score{
				"4":	Score{
					Total: 158,
					Centers: 6,
					Countback: "1234567890",
				},
			},
		},
		"6": Shooter{
			FirstName: "Raymond",
			Surname: "N",
			MiddleName: "Tapps",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "C",
			Score: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 5,
					Countback: "1234567890",
				},
			},
		},
		"7": Shooter{
			FirstName: "Leo",
			Surname: "N",
			MiddleName: "Stigger",
			Age: "U21",
			Club: 2,
			Class: "fclass",
			Grade: "FA",
			Score: map[string]Score{
				"4":	Score{
					Total: 201,
					Centers: 3,
					Countback: "1234567890",
				},
			},
		},
		"8": Shooter{
			FirstName: "Tim",
			Surname: "N",
			MiddleName: "Murray",
			Age: "U21",
			Club: 2,
			Class: "fclass",
			Grade: "FA",
			Score: map[string]Score{
				"4":	Score{
					Total: 199,
					Centers: 16,
					Countback: "1234567890",
				},
			},
		},
		"9": Shooter{
			FirstName: "Garry",
			Surname: "N",
			MiddleName: "Green",
			Age: "U21",
			Club: 2,
			Class: "fclass",
			Grade: "FB",
			Score: map[string]Score{
				"4":	Score{
					Total: 95,
					Centers: 40,
					Countback: "1234567890",
				},
			},
		},
		"10": Shooter{
			FirstName: "Fred",
			Surname: "N",
			MiddleName: "Webb",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "B",
			Score: map[string]Score{
				"4":	Score{
					Total: 187,
					Centers: 60,
					Countback: "1234567890",
				},
			},
		},
		"11": Shooter{
			FirstName: "Leesa",
			Surname: "N",
			MiddleName: "M",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "B",
			Score: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					Countback: "1034567894",
				},
			},
		},
		"12": Shooter{
			FirstName: "Leesa",
			Surname: "N",
			MiddleName: "M",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "B",
			Score: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					Countback: "1234567855",
				},
			},
		},
		"13": Shooter{
			FirstName: "Leesa",
			Surname: "N",
			MiddleName: "M",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "B",
			Score: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
				},
			},
		},
		"14": Shooter{
			FirstName: "Leesa",
			Surname: "N",
			MiddleName: "M",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "B",
			Score: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					Countback: "1234567990",
				},
			},
		},
		"15": Shooter{
			FirstName: "Leesa",
			Surname: "N",
			MiddleName: "M",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "B",
			Score: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					Countback: "1234567897",
				},
			},
		},
		"16": Shooter{
			FirstName: "Leesa",
			Surname: "N",
			MiddleName: "M",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "B",
			Score: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					Countback: "2134567890",
				},
			},
		},
		"17": Shooter{
			FirstName: "Leesa",
			Surname: "N",
			MiddleName: "M",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "B",
			Score: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					Countback: "1234567809",
				},
			},
		},
		"18": Shooter{
			FirstName: "Leesa",
			Surname: "N",
			MiddleName: "M",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "B",
			Score: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					Countback: "1234557890",
				},
			},
		},
		"19": Shooter{
			FirstName: "Leesa",
			Surname: "N",
			MiddleName: "M",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "B",
			Score: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					Countback: "1234567893",
				},
			},
		},
		"20": Shooter{
			FirstName: "Leesa",
			Surname: "N",
			MiddleName: "M",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "B",
			Score: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					Countback: "0234567890",
				},
			},
		},
		"21": Shooter{
			FirstName: "Leesa",
			Surname: "N",
			MiddleName: "M",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "B",
			Score: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					Countback: "1234567892",
				},
			},
		},
		"22": Shooter{
			FirstName: "Leesa",
			Surname: "N",
			MiddleName: "M",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "B",
			Score: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					Countback: "1234567891",
				},
			},
		},
		"23": Shooter{
			FirstName: "Leesa",
			Surname: "N",
			MiddleName: "M",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "B",
			Score: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					Countback: "1234567890",
				},
			},
		},
		"24": Shooter{
			FirstName: "Leesa",
			Surname: "N",
			MiddleName: "M",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "B",
			Score: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					Countback: "1234567890B",
				},
			},
		},
		"25": Shooter{
			FirstName: "Leesa",
			Surname: "N",
			MiddleName: "M",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "B",
			Score: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					Countback: "1234567890A",
				},
			},
		},
		"26": Shooter{
			FirstName: "Leesa",
			Surname: "N",
			MiddleName: "M",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "B",
			Score: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					Countback: "1234567890b",
				},
			},
		},
		"27": Shooter{
			FirstName: "Leesa",
			Surname: "N",
			MiddleName: "M",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "B",
			Score: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					Countback: "1234567890a",
				},
			},
		},
		"28": Shooter{
			FirstName: "Leesa",
			Surname: "N",
			MiddleName: "M",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "B",
			Score: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					Countback: "1",
				},
			},
		},
		"29": Shooter{
			FirstName: "Leesa",
			Surname: "N",
			MiddleName: "M",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "B",
			Score: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					Countback: "2",
				},
			},
		},
		"30": Shooter{
			FirstName: "Leesa",
			Surname: "N",
			MiddleName: "M",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "B",
			Score: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					Countback: "3",
				},
			},
		},
		"31": Shooter{
			FirstName: "Leesa",
			Surname: "N",
			MiddleName: "M",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "B",
			Score: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					Countback: "4",
				},
			},
		},
		"32": Shooter{
			FirstName: "Leesa",
			Surname: "N",
			MiddleName: "M",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "B",
			Score: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					Countback: "5",
				},
			},
		},
		"33": Shooter{
			FirstName: "Leesa",
			Surname: "N",
			MiddleName: "M",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "B",
			Score: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					Countback: "6",
				},
			},
		},
		"34": Shooter{
			FirstName: "Leesa",
			Surname: "N",
			MiddleName: "M",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "B",
			Score: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					Countback: "7",
				},
			},
		},
		"35": Shooter{
			FirstName: "Leesa",
			Surname: "N",
			MiddleName: "M",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "B",
			Score: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					Countback: "8",
				},
			},
		},
		"36": Shooter{
			FirstName: "Leesa",
			Surname: "N",
			MiddleName: "M",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "B",
			Score: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					Countback: "9",
				},
			},
		},
		"37": Shooter{
			FirstName: "Leesa",
			Surname: "N",
			MiddleName: "M",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "B",
			Score: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					Countback: "10",
				},
			},
		},
		"38": Shooter{
			FirstName: "Leesa",
			Surname: "N",
			MiddleName: "M",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "B",
			Score: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					Countback: "11",
				},
			},
		},
		"39": Shooter{
			FirstName: "Leesa",
			Surname: "N",
			MiddleName: "M",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "B",
			Score: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					Countback: "12",
				},
			},
		},
		"40": Shooter{
			FirstName: "Leesa",
			Surname: "N",
			MiddleName: "M",
			Age: "U21",
			Club: 2,
			Class: "target",
			Grade: "B",
			Score: map[string]Score{
				"4":	Score{
					Total: 159,
					Centers: 60,
					Countback: "13",
				},
			},
		},
	},
}

type Score struct {
	Total			int		`bson:"t,omitempty"`
	Centers		int		`bson:"c,omitempty"`
	Shots			string	`bson:"s,omitempty"`
	Countback	string	`bson:"b,omitempty"`
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
	DateTime	string					`bson:"d,omitempty"`
//	AutoInc	AutoInc					`bson:"U,omitempty"`
//	Range		map[string]Range		`bson:"R,omitempty"`
	Shooters	map[string]Shooter	`bson:"S,omitempty"`
	Club		int						`bson:"C,omitempty"`
}
