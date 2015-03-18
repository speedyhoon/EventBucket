package main

import (
	"fmt"
	"strings"
)

func generateForm(form Form) string {
	var output, formId string
	var formElements []string
	if conn == nil {
		formElements = []string{"<p class=error>Unable to connect to the EventBucket database.</p>"}
	} else {
		var attributes, element, options string
		formAttr := ""
		if form.Type != "" && form.Id != "" {
			formAttr = " form=" + addQuotes(form.Id)
			formId = " id=" + addQuotes(form.Id)
		}
		var inputValue, inputSnippet string
		for _, input := range form.Inputs {
			element = ""
			attributes = formAttr
			options = ""
			output = ""
			//devModeCheckForm(input.Html != "submit" || input.Html != "number" || input.Html != "text" || input.Html != "range" || input.Html != "datalist" || input.Html != "select" || input.Html != "date" || input.Html != "hidden", "don't use element "+input.Html)

			inputValue = fmt.Sprintf("%v", input.Value)
			if input.Html != "submit" {
				if input.Name != "" {
					attributes += " name=" + input.Name
					//devModeCheckForm(input.Name == addQuotes(input.Name), "names can't have spaces")
				}
				if input.Value != nil && inputValue != "" {
					attributes += " value=" + addQuotes(inputValue)
					//devModeCheckForm(input.Html != "select", "select boxes shouldn't have a value attribute")
				}
			} else {
				if input.Action != "" {
					attributes += " formaction=" + input.Action
				}
				//devModeCheckForm(inputValue != "", "submits should have a value")
			}

			if input.Name != "" {
				attributes += " name=" + input.Name
				//devModeCheckForm(input.Name == addQuotes(input.Name), "names can't have spaces")
			}
			if input.Action != "" {
				attributes += " formaction=" + input.Action
			}
			inputValue = fmt.Sprintf("%v", input.Value)
			if input.Value != nil && inputValue != "" {
				attributes += " value=" + addQuotes(inputValue)
				//devModeCheckForm(input.Html != "select", "select boxes shouldn't have a value attribute")
			}
			if input.Required {
				attributes += " required"
				//devModeCheckForm(input.Html == "number" || input.Html == "text" || input.Html == "range" || input.Html == "datalist" || input.Html == "date" || input.Html == "select" || input.Html == "tel", "this element shouldn't have required, type="+input.Html)
			}
			if input.Placeholder != "" && input.Html != "select" {
				attributes += " placeholder=" + addQuotes(input.Placeholder)
				//devModeCheckForm(input.Html == "text" || input.Html == "number" || input.Html == "range" || input.Html == "datalist", "placeholders are only allowed on text, datalist, number and ranges")
			}
			if input.Min != nil {
				if input.Html == "number" || input.Html == "range" {
					attributes += fmt.Sprintf(" min=%v", *input.Min)
					//devModeCheckForm(input.Html == "number" || input.Html == "range", "min is only allowed on type  number and range")
				} else if input.Html == "text" || input.Html == "email" || input.Html == "search" || input.Html == "password" || input.Html == "tel" || input.Html == "url" {
					attributes += fmt.Sprintf(" minlength=%v", *input.Min)
				}
			}
			if input.Max != nil {
				if input.Html == "number" || input.Html == "range" {
					attributes += fmt.Sprintf(" max=%v", *input.Max)
					//devModeCheckForm(input.Html == "number" || input.Html == "range", "max is only allowed on type  number and range")
				} else if input.Html == "text" || input.Html == "email" || input.Html == "search" || input.Html == "password" || input.Html == "tel" || input.Html == "url" {
					attributes += fmt.Sprintf(" maxlength=%v", *input.Min)
				}
			}
			if input.Step != 0 {
				attributes += fmt.Sprintf(" step=%v", input.Step)
				//devModeCheckForm(input.Html == "number" || input.Html == "range", "step is only allowed on type  number and range")
			}
			if input.Checked {
				attributes += " checked"
				//devModeCheckForm(input.Html == "radio" || input.Html == "checkbox", "checked is only valid on radio buttons and checkboxes")
			}
			if input.Autofocus == "on" {
				attributes += " autofocus"
			}
			if input.AccessKey != "" {
				attributes += " accesskey=" + input.AccessKey
			}
			if input.Size > 0 {
				attributes += fmt.Sprintf(" size=%d", input.Size)
				//devModeCheckForm(input.Html == "select", "size is only allowed on select tags")
				//devModeCheckForm(input.Size >= 4, "size should be >= 4")
			}
			if input.AutoComplete != "" {
				attributes += " autocomplete=" + input.AutoComplete
				//devModeCheckForm(input.Html == "datalist", "autocomplete is only allowed on datalist tags")
			}
			if input.MultiSelect {
				attributes += " multiple"
				if len(input.Options) > 4 {
					attributes += fmt.Sprintf(" size=%d", len(input.Options))
				}
				//devModeCheckForm(input.Html == "select", "multiple is only available on select boxes")
				//devModeCheckForm(input.Html != "submit", "buttons and submits shouldn't have multiple")
			}
			if len(input.Options) > 0 {
				options = drawOptions(input, input.Name)
			}
			if input.Html == "select" {
				element += "<select" + attributes + ">" + options + "</select>"
			} else if input.Html == "submit" {
				if inputValue != "" && input.Inner != "" {
					attributes += " value=" + addQuotes(inputValue)
				}
				if input.Inner != "" {
					inputValue = input.Inner
				}
				output += "<button" + attributes + ">" + inputValue + "</button>"
			} else {
				if input.Html == "datalist" && options != "" {
					attributes += " type=datalist id=" + input.Name
				}
				if input.Html != "text" {
					attributes += " type=" + input.Html
				}
				if input.Html != "" {
					element += "<input" + attributes + ">" + options
				}
			}
			if input.Label != "" && form.Type != "table" {
				var errorClass string
				if input.Error != "" {
					errorClass = " class=error"
					input.Error = " " + input.Error
				}
				output += "<label" + errorClass + ">" + input.Label + ": " + element + input.Error + "</label>"
				//devModeCheckForm(input.Html != "submit" || input.Html != "button", "submits and buttons shouldn't have lables")
			} else if element != "" {
				output += element
			}
			inputSnippet = fmt.Sprintf("%v", input.Snippet)
			if input.Snippet != nil && inputSnippet != "" {
				output += " " + inputSnippet
			}
			if input.Help != "" {
				output += "<abbr class=help title=\"" + input.Help + "\">?</abbr>"
			}
			formElements = append(formElements, output)
		}
	}
	if form.Type == "table" {
		formElements[0] = fmt.Sprintf("<form%v action=%v method=post>%v</form>", formId, addQuotes(form.Action), formElements[0])
		output = "<tr><td>" + strings.Join(formElements, "</td><td>") + "</td></tr>"
		return output
	}

	output = strings.Join(formElements, " ")
	if form.Title != "" {
		output = fieldSet(form.Title) + output + "</fieldset>"
	} else {
		//devModeCheckForm(false, "all forms should have a title")
	}
	return fmt.Sprintf("<form%v action=%v method=post>%v</form>", formId, addQuotes(form.Action), output)
}

func drawOptions(input Inputs, name string) string {
	//devModeCheckForm(len(input.Options) > 0, "select should have at least one option to select from for element='"+name+"' type='"+input.Html+"'")
	if input.Required {
		//devModeCheckForm(len(input.Options) > 0, "select shouldn't be required with no available options to select")
	}
	output := ""
	var optionValue string
	for _, option := range input.Options {
		output += "<option"
		if option.Selected {
			output += " selected"
			//devModeCheckForm(input.Html != "datalist", "datalist shouldn't have any selected values! change it to a value attribute")
			//devModeCheckForm(!(input.Placeholder != "" && input.Html != "datalist"), "shouldn't set a placeholder when options are already selected")
		}
		optionValue = fmt.Sprintf("%v", option.Value)
		if optionValue != "" {
			output += " value=" + addQuotes(optionValue)
		} else {
			//devModeCheckForm(false, "option values shouldn't be empty")
		}
		output += ">" + option.Display + "</option>"
		//devModeCheckForm(!(option.Display == "" && option.Value == "" && option.Selected == false), "option must have display text")
	}
	if input.Html == "datalist" {
		output = "<datalist id=" + name + ">" + output + "</datalist>"
	} else if input.Placeholder != "" {
		output += "<option selected value disabled>" + input.Placeholder + "</option>"
	}
	return output
}

func fieldSet(title string) string {
	return fmt.Sprintf("<fieldset><legend>%v</legend>", title)
}
