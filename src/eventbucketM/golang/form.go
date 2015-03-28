package main

import (
	"fmt"
	"strings"
)

func generateForm(form Form) string {
	var output, formID string
	var formElements []string
	if conn == nil {
		formElements = []string{"<p class=error>Unable to connect to the EventBucket database.</p>"}
	} else {
		var attributes, element, options string
		formAttr := ""
		if form.table && form.id != "" {
			formAttr = " form=" + addQuotes(form.id)
			formID = " id=" + addQuotes(form.id)
		}
		var inputValue, inputSnippet string
		for _, input := range form.inputs {
			element = ""
			attributes = formAttr
			options = ""
			output = ""
			//devModeCheckForm(input.html != "submit" || input.html != "number" || input.html != "text" || input.html != "range" || input.html != "datalist" || input.html != "select" || input.html != "date" || input.html != "hidden", "don't use element "+input.html)

			inputValue = fmt.Sprintf("%v", input.value)
			if input.html != "submit" {
				if input.name != "" {
					attributes += " name=" + input.name
					//devModeCheckForm(input.name == addQuotes(input.name), "names can't have spaces")
				}
				if input.value != nil && inputValue != "" {
					attributes += " value=" + addQuotes(inputValue)
					//devModeCheckForm(input.html != "select", "select boxes shouldn't have a value attribute")
				}
			} else {
				if input.html == "submit" && input.action != "" {
					attributes += " formaction=" + input.action
				}
				//devModeCheckForm(inputValue != "", "submits should have a value")
				if input.name != "" {
					attributes += " name=" + input.name
					//devModeCheckForm(input.name == addQuotes(input.name), "names can't have spaces")
				}
			}
			if input.action != "" {
				attributes += " formaction=" + input.action
			}
			inputValue = fmt.Sprintf("%v", input.value)
			if input.value != nil && inputValue != "" {
				attributes += " value=" + addQuotes(inputValue)
				//devModeCheckForm(input.html != "select", "select boxes shouldn't have a value attribute")
			}
			if input.required {
				attributes += " required"
				//devModeCheckForm(input.html == "number" || input.html == "text" || input.html == "range" || input.html == "datalist" || input.html == "date" || input.html == "select" || input.html == "tel", "this element shouldn't have required, type="+input.html)
			}
			if input.placeholder != "" && input.html != "select" {
				attributes += " placeholder=" + addQuotes(input.placeholder)
				//devModeCheckForm(input.html == "text" || input.html == "number" || input.html == "range" || input.html == "datalist", "placeholders are only allowed on text, datalist, number and ranges")
			}
			if input.min != nil {
				if input.html == "number" || input.html == "range" {
					attributes += fmt.Sprintf(" min=%v", *input.min)
					//devModeCheckForm(input.html == "number" || input.html == "range", "min is only allowed on type  number and range")
				} else if input.html == "text" || input.html == "email" || input.html == "search" || input.html == "password" || input.html == "tel" || input.html == "url" {
					attributes += fmt.Sprintf(" minlength=%v", *input.min)
				}
			}
			if input.max != nil {
				if input.html == "number" || input.html == "range" {
					attributes += fmt.Sprintf(" max=%v", *input.max)
					//devModeCheckForm(input.html == "number" || input.html == "range", "max is only allowed on type  number and range")
				} else if input.html == "text" || input.html == "email" || input.html == "search" || input.html == "password" || input.html == "tel" || input.html == "url" {
					attributes += fmt.Sprintf(" maxlength=%v", *input.min)
				}
			}
			if input.step != 0 {
				attributes += fmt.Sprintf(" step=%v", input.step)
				//devModeCheckForm(input.html == "number" || input.html == "range", "step is only allowed on type  number and range")
			}
			if input.checked {
				attributes += " checked"
				//devModeCheckForm(input.html == "radio" || input.html == "checkbox", "checked is only valid on radio buttons and checkboxes")
			}
			if input.autofocus == "on" {
				attributes += " autofocus"
			}
			if input.accessKey != "" {
				attributes += " accesskey=" + input.accessKey
			}
			if input.size > 0 {
				attributes += fmt.Sprintf(" size=%d", input.size)
				//devModeCheckForm(input.html == "select", "size is only allowed on select tags")
				//devModeCheckForm(input.size >= 4, "size should be >= 4")
			}
			if input.autoComplete != "" {
				attributes += " autocomplete=" + input.autoComplete
				//devModeCheckForm(input.html == "datalist", "autocomplete is only allowed on datalist tags")
			}
			if input.multiSelect {
				attributes += " multiple"
				if len(input.options) > 4 {
					attributes += fmt.Sprintf(" size=%d", len(input.options))
				}
				//devModeCheckForm(input.html == "select", "multiple is only available on select boxes")
				//devModeCheckForm(input.html != "submit", "buttons and submits shouldn't have multiple")
			}
			if len(input.options) > 0 {
				options = drawOptions(input)
			}
			if input.maxLength > 0 {
				attributes += fmt.Sprintf(" maxlength=%v", input.maxLength)
			}
			if input.minLength > 1 { //Only adds it if MinLength is 2 - browse doesn't enforce 1! Using Required attribute enforces MinLength of 1
				attributes += fmt.Sprintf(" minlength=%v", input.minLength)
			}
			if input.html == "select" {
				element += "<select" + attributes + ">" + options + "</select>"
			} else if input.html == "submit" {
				if input.value != nil && inputValue != "" && input.inner != "" {
					attributes += " value=" + addQuotes(inputValue)
				}
				output += "<button" + attributes + ">" + input.inner + "</button>"
			} else {
				if input.dataList && options != "" {
					if input.id == "" {
						Warning.Println("datalist needs a unique ID")
					}
					if input.html == "text" {
						Warning.Println("datalist type should be search")
					}
					attributes += " type=" + input.html + " list=" + input.id
				}
				if input.html != "text" {
					attributes += " type=" + input.html
				}
				if input.html != "" {
					element += "<input" + attributes + ">" + options
				}
			}
			if input.label != "" && !form.table {
				var errorClass string
				if input.error != "" {
					errorClass = " class=error"
					input.error = " " + input.error
				}
				output += "<label" + errorClass + ">" + input.label + ": " + element + input.error + "</label>"
				//devModeCheckForm(input.html != "submit" || input.html != "button", "submits and buttons shouldn't have lables")
			} else if element != "" {
				output += element
			}
			inputSnippet = fmt.Sprintf("%v", input.snippet)
			if input.snippet != nil && inputSnippet != "" {
				output += " " + inputSnippet
			}
			if input.help != "" {
				output += "<abbr class=help title=\"" + input.help + "\">?</abbr>"
			}
			formElements = append(formElements, output)
		}
	}
	if form.table {
		formElements[0] = fmt.Sprintf("<form%v action=%v method=post>%v</form>", formID, addQuotes(form.action), formElements[0])
		output = "<tr><td>" + strings.Join(formElements, "</td><td>") + "</td></tr>"
		return output
	}

	output = strings.Join(formElements, " ")
	if form.title != "" {
		output = fieldSet(form.title) + output + "</fieldset>"
	} else {
		//devModeCheckForm(false, "all forms should have a title")
	}
	return fmt.Sprintf("<form%v action=%v method=post>%v</form>", formID, addQuotes(form.action), output)
}

func drawOptions(input Inputs) string {
	//devModeCheckForm(len(input.options) > 0, "select should have at least one option to select from for element='"+input.name+"' type='"+input.html+"'")
	if input.required {
		//devModeCheckForm(len(input.options) > 0, "select shouldn't be required with no available options to select")
	}
	output := ""
	var optionValue string
	for _, option := range input.options {
		output += "<option"
		if option.Selected {
			output += " selected"
			//devModeCheckForm(input.html != "datalist", "datalist shouldn't have any selected values! change it to a value attribute")
			//devModeCheckForm(!(input.placeholder != "" && input.html != "datalist"), "shouldn't set a placeholder when options are already selected")
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
	if input.dataList {
		output = "<datalist id=" + input.id + ">" + output + "</datalist>"
	} else if input.placeholder != "" {
		output += "<option selected value disabled>" + input.placeholder + "</option>"
	}
	return output
}

func fieldSet(title string) string {
	return fmt.Sprintf("<fieldset><legend>%v</legend>", title)
}
