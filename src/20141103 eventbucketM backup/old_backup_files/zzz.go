package main

import (
	"fmt"
)

func generateForm2(form Form) string {
	var output, attributes, element, options string
	for name, input := range form.Inputs {
		element = ""
		attributes = ""
		options = ""

		if name != "" {
			attributes += " name=" + name
		}
		if input.Value != ""{
			if input.Html != "submit" {
				attributes += " value=" + addQuotes(input.Value))
			}
		}
		if input.Required {
			attributes += " required"
		}
		if input.Placeholder != "" {
			attributes += " placeholder="+addQuotes(input.Placeholder)
		}
		if input.Min > -1 {
			attributes += fmt.Sprintf(" min=%v", input.Min)
		}
		if input.Max > -1 {
			attributes += fmt.Sprintf(" max=%v", input.Max)
		}
		if input.Step > 0 {
			attributes += fmt.Sprintf(" step=%v", input.Step)
		}
		if input.Checked {
			attributes += " checked"
		}
		if input.Size > 0 {
			attributes += " size=%d" + input.Size
		}
		if input.AutoComplete != "" {
			attributes += " autocomplete="+input.AutoComplete
		}

		if input.MultiSelect {
			attributes += " multiple"
			if len(input.Options) > 4 {
				attributes += fmt.Sprintf(" size=%d", len(input.Options))
			}
		}
		if input.Html == "datalist"{
			attributes += " id=" + name
		}
		options = draw_options(input, name)
		if input.Help != "" {
			attributes += "title=" + addQuotes(input.Help)
		}


		if input.Html == "select" {
			element += "<select"+attributes+">"+options+"</select>"
		}else if input.Html == "submit" {
			element += "<button"+attributes+">"+input.Value+options+"</button>"
		}else {
			if input.Html != "text" {
				attributes += " type="+input.Html
			}
			element += "<input"+attributes+">"
		}
		if input.Label != "" {
			output += "<label>"+input.Label+": "+element+"</label>"
		}
	}
	if form.Title != "" {
		output = field_set(form.Title) + output + "</fieldset>"
	}
	return fmt.Sprintf("<form action=%v method=post>%v</form>", addQuotes(form.Action), output)
}

type Form struct {
	Action string
	Title  string
	Inputs map[string]Inputs
}
type Inputs struct {
	//AutoComplete values can be: "off" or "on"
	Html, Label, Help, Value, Pattern, Placeholder, AutoComplete  string
	Checked, MultiSelect, Required bool
	Min, Max, Size                 int
	Options                        []Options
	Step                           float64
}
type Options struct {
	Value    string `json:"v,omitempty"`
	Display  string `json:"d,omitempty"`
	Selected bool   `json:"s,omitempty"`
}

func draw_options(input Inputs, name string)string{
	if len(input.Options) <= 0 {
		return ""
	}
	output := ""
	if input.Placeholder != "" && input.Html != "datalist"{
		output += "<option selected value disabled>"+input.Placeholder+"</option>"
	}
	for _, option := range input.Options {
		output += "<option"
		if option.Selected {
			output += " selected"
		}
		if option.Value != "" {
			output += " value" + addQuotesEquals(option.Value)
		}
		output += ">" + option.Display + "</option>"
	}
	if input.Html == "datalist"{
		output = "<datalist id=" + name + ">"+output+"</datalist>"
	}
	return output
}

func field_set(title string) string {
	return fmt.Sprintf("<fieldset><legend>%v</legend>", title)
}
//6,493 bytes
