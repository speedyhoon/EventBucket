

//type ErrorMsg struct {
//	Title, Message string
//	Info           bool
//}

/*var error_queue []ErrorMsg

func error_message(info bool, index, title, message string) {
	for _, queue_item := range error_queue {
		if queue_item.Title == title {
			return
		}
	}
	error_queue = append(error_queue, ErrorMsg{
		//	error_queue[index] = ErrorMsg{
		Title:   title,
		Message: message,
		Info:    info,
	})
	//	}
}
func remove_error(title string) {
	for index, queue_item := range error_queue {
		if queue_item.Title == title {
			error_queue[index] = ErrorMsg{}
		}
	}
}*/

//func render_errors() string {
//	var output string
//	if len(error_queue) >= 1 {
//		for _, error := range error_queue {
//			if error.Title != "" && error.Message != "" {
//				var class = "error"
//				if error.Info {
//					class = "info"
//				}
//				output += fmt.Sprintf("<div class=%v><h2>%v:</h2>%v</div>", class, error.Title, error.Message)
//			}
//		}
//		error_queue = []ErrorMsg{}
//	}
//	return output
//}