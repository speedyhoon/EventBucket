package main

import (
	"fmt"
)

func main() {
	var number uint64 = 1
	fmt.Println(number)
	number *= 2
	fmt.Println(number)
	for{
		number *= 2
		fmt.Println(number)
		if number == 0 { break }
	}
	number = 9223372036854775807
	
	
	for{
		number += 1
		fmt.Println(number)
		if number == 0 { break }
	}
}