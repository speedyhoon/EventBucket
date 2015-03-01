package main

import(
	"fmt"
)
func main(){
	a := 0
	for a < 256 {
		fmt.Println(id_suffix(a))
		a += 1
	}
}

func id_suffix(id int)string{
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789~!*()_-."
	if id <= 0 {
		//error_message(false, "Invalid id number supplied.", fmt.Sprintf("Id \"%v\" is out of range", id))
		return fmt.Sprintf("%c", charset[0])
	}
	id = id - 1
//	fmt.Printf("charset lenght = %v", len(charset))
	charset_length := 70
	temp := ""
	for id >= charset_length {
		temp = fmt.Sprintf("%c%v", charset[id % charset_length], temp)
		id = id/charset_length-1
	}
	return fmt.Sprintf("%c%v", charset[id % charset_length], temp)
}
