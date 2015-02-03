package main

import (
	"fmt"
	"strings"
)

func generateForm2(form Form) string {
	var output string
	var formElements []string
	if database_status == true {
		var attributes, element, options string
		for _, input := range form.Inputs {
			element = ""
			attributes = ""
			options = ""
			output = ""
			dev_mode_check_form(input.Html != "submit" || input.Html != "number" || input.Html != "text" || input.Html != "range" || input.Html != "datalist" || input.Html != "select" || input.Html != "date" || input.Html != "hidden", "don't use element "+input.Html)

			if input.Html != "submit" {
				if input.Name != "" {
					attributes += " name="+input.Name
					dev_mode_check_form(input.Name == addQuotes(input.Name), "names can't have spaces")
				}
				if input.Value != "" {
					attributes += " value="+addQuotes(input.Value)
					dev_mode_check_form(input.Html != "select", "select boxes shouldn't have a value attribute")
				}
			}else {
				dev_mode_check_form(input.Value != "", "submits should have a value")
			}
			if input.Required {
				attributes += " required"
				dev_mode_check_form(input.Html == "number" || input.Html == "text" || input.Html == "range" || input.Html == "datalist" || input.Html == "date" || input.Html == "select" || input.Html == "tel", "this element shouldn't have required, type="+input.Html)
			}
			if input.Placeholder != "" && input.Html != "select" {
				attributes += " placeholder="+addQuotes(input.Placeholder)
				dev_mode_check_form(input.Html == "text" || input.Html == "number" || input.Html == "range" || input.Html == "datalist", "placeholders are only allowed on text, datalist, number and ranges")
			}
			if input.Min != "" {
				attributes += fmt.Sprintf(" min=%v", input.Min)
				dev_mode_check_form(input.Html == "number" || input.Html == "range", "min is only allowed on type  number and range")
			}
			if input.Max != "" {
				attributes += fmt.Sprintf(" max=%v", input.Max)
				dev_mode_check_form(input.Html == "number" || input.Html == "range", "max is only allowed on type  number and range")
			}
			if input.Step != 0 {
				attributes += fmt.Sprintf(" step=%v", input.Step)
				dev_mode_check_form(input.Html == "number" || input.Html == "range", "step is only allowed on type  number and range")
			}
			if input.Checked {
				attributes += " checked"
				dev_mode_check_form(input.Html == "radio" || input.Html == "checkbox", "checked is only valid on radio buttons and checkboxes")
			}
			if input.Size > 0 {
				attributes += fmt.Sprintf(" size=%d", input.Size)
				dev_mode_check_form(input.Html == "select", "size is only allowed on select tags")
				dev_mode_check_form(input.Size >= 4, "size should be >= 4")
			}
			if input.AutoComplete != "" {
				attributes += " autocomplete="+input.AutoComplete
				dev_mode_check_form(input.Html == "datalist", "autocomplete is only allowed on datalist tags")
			}

			if input.MultiSelect {
				attributes += " multiple"
				if len(input.Options) > 4 {
					attributes += fmt.Sprintf(" size=%d", len(input.Options))
				}
				dev_mode_check_form(input.Html == "select", "multiple is only available on select boxes")
				dev_mode_check_form(input.Html != "submit", "buttons and submits shouldn't have multiple")
			}
			if len(input.Options) > 0 {
				options = draw_options(input, input.Name)
			}
			if input.Help != "" {
				attributes += "title="+addQuotes(input.Help)
			}
			if input.Html == "select" {
				element += "<select"+attributes+">"+options+"</select>"
			}else if input.Html == "submit" {
				output += "<button"+attributes+">"+input.Value+"</button>"
			}else {
				if input.Html == "datalist" && options != "" {
					attributes += " type=datalist id="+input.Name
				}
				if input.Html != "text" {
					attributes += " type="+input.Html
				}
				element += "<input"+attributes+">"+options
			}
			if input.Label != "" {
				output += "<label>"+input.Label+": "+element+"</label>"
				dev_mode_check_form(input.Html != "submit" || input.Html != "button", "submits and buttons shouldn't have lables")
			}else {
				output += element
			}
			formElements = append(formElements, output)
		}
	}else {
		formElements = []string{"<p>Unable to connect to the EventBucket database</p>"}
	}
	output = strings.Join(formElements, " ")
	if form.Title != "" {
		output = field_set(form.Title) + output + "</fieldset>"
	}else {
		dev_mode_check_form(false, "all forms should have a title")
	}
	return fmt.Sprintf("<form action=%v method=post>%v</form>", addQuotes(form.Action), output)
}

type Form struct {
	Action string
	Title  string
	Inputs []Inputs
}
type Inputs struct {
	Name, Html, Label, Help, Value, Pattern, Placeholder, Min, Max, AutoComplete string	//AutoComplete values can be: "off" or "on"
	Checked, MultiSelect, Required bool
	Size    int
	Options []Option
	Step    float64
	VarType string		//the type of variable to return
	VarMaxLen int		//the length of variable to return
	VarMinLen int		//the length of variable to return
}
type Option struct {
	Value    string `json:"v,omitempty"`
	Display  string `json:"d,omitempty"`
	Selected bool   `json:"s,omitempty"`
}

func draw_options(input Inputs, name string)string{
	dev_mode_check_form(len(input.Options) > 0, "select should have at least one option to select from for element='"+name+"' type='"+input.Html+"'")
	if input.Required {
		dev_mode_check_form(len(input.Options) > 0, "select shouldn't be required with no available options to select")
	}
	output := ""
	for _, option := range input.Options {
		output += "<option"
		if option.Selected {
			output += " selected"
			dev_mode_check_form(input.Html != "datalist", "datalist shouldn't have any selected values! change it to a value attribute")
			dev_mode_check_form(!(input.Placeholder != "" && input.Html != "datalist"), "shouldn't set a placeholder when options are already selected")
		}
		if option.Value != "" {
			output += " value=" + addQuotes(option.Value)
		}else {
			dev_mode_check_form(false, "option values shouldn't be empty")
		}
		output += ">" + option.Display + "</option>"
		dev_mode_check_form(!(option.Display==""&&option.Value==""&&option.Selected==false), "option must have display text")
	}
	if input.Html == "datalist"{
		output = "<datalist id=" + name + ">"+output+"</datalist>"
		//dev_mode_check_form(false,"make sure datalist id='"+name+"' is unique!")
	}else if input.Placeholder != ""{
		output += "<option selected value disabled>"+input.Placeholder+"</option>"
	}
	return output
}

func field_set(title string) string {
	return fmt.Sprintf("<fieldset><legend>%v</legend>", title)
}
