package main

import (
	"fmt"
)

func main(){
	var inter interface{} = []string{"a", "b", "c"}
	sa, ok := inter.([]string)
	if !ok {
		fmt.Println("not a []string")
	}else{
		for ta, s := range(sa) {
			fmt.Println(ta)
			fmt.Println(s)
		}
	}
	fmt.Print(sa[0])
	fmt.Print("\n\n\n\n")

	try_its := map[int]interface{}{
		0:"name",	//event name
		1:123456,	//date time stamp
		2:2,			//club id
		3: map[int]interface{}{		//shooters
			0:map[int]interface{}{			//shooters id
				0:"grade",
				1:"class",
				2:map[int]interface{}{		//scores
					0:map[int]interface{}{		//range id
						0: 40,		//total
						1: 5,		//centers
						2: 7,		//x's
						3: 2345678976,	//countback
						4: 987456788,	//xcountback
					},
					1:map[int]interface{}{		//range id
						0: 40,		//total
						1: 5,		//centers
						2: 7,		//x's
						3: 2345678976,	//countback
						4: 987456788,	//xcountback
					},
				},
			},
		},
		4:"settings",
		5:"team",
		6:"teamcat",
		7:"handicap?",//or should this go under shooter?
	}
//	temp := [3]string{"40","5","8754453657"}
//	green := [3]int{32, 6565, 32432}
//	my_map := map[int]interface{}{0:"a", 1:"4", 2:"c", 3:temp, 4:green}
	var inter2 interface{} = try_its
	ff, ok := inter2.(map[int]interface{})
	if ok {
		for gg, d := range(ff) {
			fmt.Println(gg)
			fmt.Println(d)
		}
	}



}
