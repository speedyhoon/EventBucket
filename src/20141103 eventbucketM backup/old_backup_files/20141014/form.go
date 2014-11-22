package main

import (
	"fmt"
)

func _required(attribute bool)string{
	if attribute {
		return " required"
	}
	return ""
}

func generateForm2(form Form) string {
	output := ""
	for inputName, inputData := range form.Inputs {
		if inputData.Html != "submit" && inputData.Label != "" {
			output += " <label>" + inputData.Label + ": "
		}
		if inputData.Html == "select" {
			output += "<select"
			if inputName != "" {
				output += " name=" + addQuotes(inputName)
			}
			if inputData.MultiSelect {
				output += " multiple"
				if len(inputData.SelectedValues) > 1 {
					output += fmt.Sprintf(" size=%d", len(inputData.SelectedValues))
				} else if len(inputData.SelectValues) > 1 {
					output += fmt.Sprintf(" size=%d", len(inputData.SelectValues))
				} else if len(inputData.Select) > 1 {
					output += fmt.Sprintf(" size=%d", len(inputData.Select))
				}
			}
			output += _required(inputData.Required)
			if inputData.Help != "" {
				output += fmt.Sprintf("title=%v", addQuotes(inputData.Help))
			}
			output += ">"
			//TODO <option value="" disabled selected>Select your option</option>
			options, selected_options := draw_list_box(inputData.SelectedValues)
			if inputData.Placeholder != "" && !selected_options {
				output += fmt.Sprintf("<option disabled selected value>%v</option>", inputData.Placeholder)
			}
			for _, option := range inputData.Select {
				output += fmt.Sprintf("<option>%v</option>", option)
			}
			output += build_options_deprecated(inputData.SelectValues)

			output += options
			output += "</select>"
		} else if inputData.Html == "submit" {
			//TODO change all the submit labels to values
			if inputData.Value != "" && inputData.Label != "" {
				output += inputData.Label + " <button>" + inputData.Value + "</button>"
			} else {
				output += "<button>" + inputData.Label + "</button>"
			}
		} else if inputData.Html == "datalist" {
			output += "<input"
			if inputName != "" {
				output += " name=" + addQuotes(inputName)
			}
			if inputData.Required {
				output += " required"
			}
			if inputData.Value != "" {
				output += fmt.Sprintf(" value=%v", addQuotes(inputData.Value))
			}
			if inputData.AutoComplete != ""{
				output += " autocomplete=" + inputData.AutoComplete
			}
			if inputData.Placeholder != "" {
				output += " placeholder=" + addQuotes(inputData.Placeholder)
			}
			if len(inputData.Options) > 0 {
				output += " list=" + inputName + "><datalist id=" + inputName + ">"
				for _, option := range inputData.Options {
					output += fmt.Sprintf("<option value=%v>", addQuotes(option))
				}
				output += "</datalist>"
			}else if len(inputData.SelectedValues) > 0 {
				//TODO remove this old datalist generator
				output += " list=" + inputName + ">"

				output += "<datalist id=" + inputName + ">"
				for value, option := range inputData.SelectValues {
					if value != "" {
						output += fmt.Sprintf("<option value=%v>%v</option>", addQuotes(value), option)
					} else {
						output += fmt.Sprintf("<option>%v</option>", option)
					}
				}
				output += "</datalist>"
			} else {
				output += ">"
			}
		} else {
			output += "<input"
			if inputData.Html != "text" {
				output += " type=" + inputData.Html
			}
			if inputData.Html != "submit" {
				output += " name=" + inputName
			}
			if inputData.AutoComplete != ""{
				output += " autocomplete=" + inputData.AutoComplete
			}
			if inputData.Html == "number" || inputData.Html == "range" {
				if inputData.Min > -1 {
					output += " min=" + echo(inputData.Min)
				}
				if inputData.Max > -1 {
					output += " max=" + echo(inputData.Max)
				}
				if inputData.Step > 0 {
					output += " step=" + echo(inputData.Step)
				}
			}
			if inputData.Checked {
				output += " checked"
			}
			if inputData.Size > 0 {
				output += fmt.Sprintf(" size=%d", inputData.Size)
			}
			if inputData.Placeholder != "" {
				output += " placeholder=" + addQuotes(inputData.Placeholder)
			}
			if inputData.Required {
				if inputData.Html != "hidden" {
					output += " required"
				} else {
					fmt.Println("\nhidden inputs are not allowed to have required attributes\n")
				}
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
	if form.Title != "" {
		output = field_set(form.Title) + output + "</fieldset>"
	}
	return fmt.Sprintf("<form action=%v method=post>%v</form>", addQuotes(form.Action), output)
}

//TODO change map to slice of SelectedValues
//func build_options(options []SelectedValues)string{
func build_options_deprecated(options map[string]string) string {
	output := ""
	for value, option := range options {
		output += fmt.Sprintf("<option value%v>%v</option>", addQuotesEquals(value), option)
	}
	return output
}
func build_options(options []SelectedValues) string {
	output := ""
	selected := ""
	for _, option := range options {
		if option.Selected {
			selected = " selected"
		}
		output += fmt.Sprintf("<option%v value%v>%v</option>", selected, addQuotesEquals(option.Value), option.Display)
		selected = ""
	}
	return output
}

type SelectedValues struct {
	Value    string `json:"v,omitempty"`
	Display  string `json:"d,omitempty"`
	Selected bool   `json:"s,omitempty"`
}

type Inputs struct {
	Html, Label, Help, Value       string
	Placeholder                    string
	Select                         []string
	SelectValues                   map[string]string
	SelectedValues                 []SelectedValues
	Options								 []string
	Checked, MultiSelect, Required bool
	Min, Max                       int
	Size                           int
	Step                           float64
	Pattern                        string

	//TODO maybe use bool? but how to detect if it is not set?
	AutoComplete						string	//0="off", 1="on"
}
type Form struct {
	Action string
	Title  string
	Inputs map[string]Inputs
}

func draw_list_box(options []SelectedValues) (string, bool) {
	output := ""
	selected_option := false
	for _, option := range options {
		output += "<option"
		if option.Selected {
			output += " selected"
			selected_option = true
		}
		if option.Value != "" {
			output += " value=" + addQuotes(option.Value)
		}
		output += ">" + option.Display + "</option>"
	}
	return output, selected_option
}

func field_set(title string) string {
	return fmt.Sprintf("<fieldset><legend>%v</legend>", title)
}
//6,493 bytes
