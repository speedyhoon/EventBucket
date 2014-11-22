package main

import (
	"fmt"
)

func generateForm(action string, form map[string]Inputs) string {
	output := fmt.Sprintf("<form action=%v method=post>", addQuotes(action))
	for inputName, inputData := range form {
		if inputData.Html != "submit" && inputData.Label != "" {
			output += " <label>"+inputData.Label+": "
		}
		if inputData.Html == "select" {
			output += "<select"
			if inputName != "" {
				output += " name="+addQuotes(inputName)
			}
			if inputData.MultiSelect {
				output += " multiple"
				if len(inputData.Select) > 1{
					output += fmt.Sprintf(" size=%d", len(inputData.Select))
				}else if len(inputData.SelectValues) > 1{
					output += fmt.Sprintf(" size=%d", len(inputData.SelectValues))
				}
			}
			if inputData.Required {
				output += " required"
			}
			if inputData.Help != "" {
				output += fmt.Sprintf("title=%v", addQuotes(inputData.Help))
			}
			output += ">"
			//TODO <option value="" disabled selected>Select your option</option>
			if inputData.Placeholder != "" {
				output+= fmt.Sprintf("<option disabled selected>%v</option>", inputData.Label)
			}
			for _, option := range inputData.Select {
				output += fmt.Sprintf("<option>%v</option>", option)
			}

			for value, option := range inputData.SelectValues {
				if value != ""{
					output += fmt.Sprintf("<option value=%v>%v</option>", addQuotes(value), option)
				}else {
					output += fmt.Sprintf("<option>%v</option>", option)
				}
			}
			output += draw_list_box(inputData.SelectedValues)
			output += "</select>"
		}else if inputData.Html == "submit" {
			output += "<button>"+inputData.Label+"</button>"
		}else {
			output += "<input"
			if inputData.Html != "text" {
				output += " type="+inputData.Html
			}
			if inputData.Html != "submit" {
				output += " name="+inputName
			}
			if inputData.Help == "number" || inputData.Help == "range"{
				if inputData.Min > -1 {
					output += " min="+echo(inputData.Min)
				}
				if inputData.Max > -1 {
					output += " max="+echo(inputData.Max)
				}
			}
			if inputData.Checked {
				output += " checked"
			}
			if inputData.Size > 0{
				output += fmt.Sprintf(" size=%d", inputData.Size)
			}
			if inputData.Placeholder != "" {
				output += " placeholder="+inputData.Placeholder
			}
			if inputData.Required {
				output += " required"
			}
			if inputData.Help != "" {
				output += fmt.Sprintf("title=%v", addQuotes(inputData.Help))
			}
			if inputData.Value != "" {
				output += fmt.Sprintf(" value=%v", addQuotes(inputData.Value))
			}
			output += ">"
		}
		if inputData.Html != "submit" && inputData.Label != "" {
			output += "</label>"
		}
	}
	output += "</form>"
	return output
}
