package main

import(
//	"fmt"
)

type ErrorMsg struct{
	Title, Message string
	Info bool
}

var error_queue []ErrorMsg

func main(){
	//	error_queue = append(error_queue, ErrorMsg{
//	error_queue[999] = ErrorMsg{
//	Title: "title",
//	Message: "message",
//		Info: true,
//	}
	error_queue[999] = make(ErrorMsg{})
}
