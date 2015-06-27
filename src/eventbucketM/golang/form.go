package main

import (
	"fmt"
	"strconv"
	"strings"
)

var on = 1
var formAttrs = map[string]Inputs{
	//submit, button, reset
	"submit": {
		autofocus: true,
		//disabled
		//form
		action: "1",
		//formenctype
		//formmethod
		//formnovalidate
		//formtarget
		name:  "1",
		value: "1",
	},
	"datalist": {
		id: "1",
	},
	"meter": {
		min: &on,
		max: &on,
		//low
		//high
		//optimum
		//form
		value: "1",
	},
	"output": {
		//for
		//form
		name: "1",
	},
	"select": {
		autofocus: true,
		//disabled
		//form
		multiSelect: true,
		name:        "1",
		required:    true,
		size:        1,
	},
	"textarea": {
		autoComplete: "1",
		//cols
		//disabled
		//form
		maxLength:   1,
		minLength:   1,
		name:        "1",
		placeholder: "1",
		//readonly
		required: true,
		//rows
		//selectionDirection
		//selectionEnd
		//selectionStart
		//spellcheck
		//wrap
	},
	"text": {
		//accept
		//autocomplete
		autofocus: true,
		//autosave
		//checked: true,
		//disabled
		//form
		//action: "1",
		//formmethod
		//formnovalidate
		//formtarget
		//inputmode
		dataList: true,
		//max: 1,
		//min: 1,
		//step: 1,
		maxLength: 1,
		minLength: 1,
		//multiSelect: true,
		name:        "1",
		pattern:     "1",
		placeholder: "1",
		//readonly
		required: true,
		size:     1,
		//spellcheck
		value: "1",
	},
	"checkbox": { //checkbox, radio
		//accept
		//autocomplete
		autofocus: true,
		//autosave
		checked: true,
		//disabled
		//form
		//action: "1",
		//formmethod
		//formnovalidate
		//formtarget
		//inputmode
		//dataList: true,
		//max: 1,
		//min: 1,
		//maxLength: 1,
		//minLength: 1,
		//multiSelect: true,
		name: "1",
		//pattern: "1",
		//placeholder: "1",
		//readonly
		required: true,
		//size: 1,
		//spellcheck
		value: "1",
		html:  "1",
	},
	//colour
	"date": {
		//accept
		//autocomplete
		autofocus: true,
		//autosave
		//checked: true,
		//disabled
		//form
		//action: "1",
		//formmethod
		//formnovalidate
		//formtarget
		//inputmode
		//dataList: true,
		//max: 1,
		//min: 1,
		//maxLength: 1,
		//minLength: 1,
		//multiSelect: true,
		name: "1",
		//pattern: "1",
		//placeholder: "1",
		//readonly
		required: true,
		//size: 1,
		//spellcheck
		value: "1",
		html:  "1",
	},
	//datetime
	//datetime-local
	//email
	//file
	"hidden": {
		//accept
		//autocomplete
		//autofocus: true,
		//autosave
		//checked: true,
		//disabled
		//form
		//action: "1",
		//formmethod
		//formnovalidate
		//formtarget
		//inputmode
		//dataList: true,
		//max: 1,
		//min: 1,
		//maxLength: 1,
		//minLength: 1,
		//multiSelect: true,
		name: "1",
		//pattern: "1",
		//placeholder: "1",
		//readonly
		//required: true,
		//size: 1,
		//spellcheck
		value: "1",
		html:  "1",
	},
	//image
	//month
	"number": { //number, range
		//accept
		//autocomplete
		autofocus: true,
		//autosave
		//checked: true,
		//disabled
		//form
		//action: "1",
		//formmethod
		//formnovalidate
		//formtarget
		//inputmode
		//dataList: true,
		max:  &on,
		min:  &on,
		step: 1,
		//maxLength: 1,
		//minLength: 1,
		//multiSelect: true,
		name: "1",
		//pattern: "1",
		//placeholder: "1",
		//readonly
		required: true,
		//size: 1,
		//spellcheck
		value: "1",
		html:  "1",
	},
	//password
	//radio
	//range
	//reset
	"search": { //number, range
		//accept
		//autocomplete
		autofocus: true,
		//autosave
		//checked: true,
		//disabled
		//form
		//action: "1",
		//formmethod
		//formnovalidate
		//formtarget
		//inputmode
		dataList: true,
		//max: 1,
		//min: 1,
		//step: 1,
		maxLength: 1,
		minLength: 1,
		//multiSelect: true,
		name:        "1",
		pattern:     "1",
		placeholder: "1",
		//readonly
		required: true,
		size:     1,
		//spellcheck
		value: "1",
		html:  "1",
	},
	//tel
	"time": { //number, range
		//accept
		//autocomplete
		autofocus: true,
		//autosave
		//checked: true,
		//disabled
		//form
		//action: "1",
		//formmethod
		//formnovalidate
		//formtarget
		//inputmode
		//dataList: true,
		max:  &on,
		min:  &on,
		step: 1,
		//maxLength: 1,
		//minLength: 1,
		//multiSelect: true,
		name: "1",
		//pattern: "1",
		//placeholder: "1",
		//readonly
		required: true,
		//size: 1,
		//spellcheck
		value: "1",
		html:  "1",
	},
	//url
	//week
}

func generateForm(form Form) string {
	if conn == nil {
		return "<p class=error>Unable to connect to the EventBucket database.</p>"
	}
	var formElements []string
	var formID, attributes, element, options string
	for _, input := range form.inputs {
		allowedAttrs := formAttrs[input.html]
		if input.snippet != nil {
			element = fmt.Sprintf("%v", input.snippet)
		}
		if allowedAttrs.html != "" {
			attributes += " type=" + input.html
		}
		if allowedAttrs.name != "" && input.name != "" {
			attributes += " name=" + input.name
		}
		if input.accessKey != "" {
			attributes += " accesskey=" + input.accessKey
		}
		if allowedAttrs.autoComplete != "" && input.autoComplete == "on" || input.autoComplete == "off" {
			attributes += " autocomplete=" + input.autoComplete
		}
		if allowedAttrs.autofocus && input.autofocus {
			attributes += " autofocus"
		}
		if allowedAttrs.checked && input.checked {
			attributes += " checked"
		}
		if allowedAttrs.action != "" && input.action != "" {
			attributes += " formaction=" + input.action
		}
		if allowedAttrs.dataList && input.dataList && input.id != "" {
			attributes += " list=" + input.id
			if len(input.options) > 0 {
				options = "<datalist id=" + input.id + ">" + drawOnlyOptions(input.options) + "</datalist>"
			}
		}
		if allowedAttrs.max != nil && input.max != nil {
			attributes += fmt.Sprintf(" max=%v", *input.max)
		}
		if allowedAttrs.maxLength > 0 && input.maxLength > 0 {
			attributes += fmt.Sprintf(" maxlength=%v", input.maxLength)
		}
		if allowedAttrs.min != nil && input.min != nil {
			attributes += fmt.Sprintf(" min=%v", *input.min)
		}
		if allowedAttrs.minLength > 0 && ((input.required && input.minLength > 1) || (!input.required && input.minLength > 0)) {
			attributes += fmt.Sprintf(" minlength=%v", input.minLength)
		}
		if allowedAttrs.multiSelect && input.multiSelect {
			attributes += " multiple"
			if input.html == "select" && len(input.options) > 0 {
				attributes += fmt.Sprintf(" size=%v", min(len(input.options), 12))
			}
		}
		if input.html == "select" {
			options = drawOnlyOptions(input.options)
		}
		if allowedAttrs.pattern != "" && input.pattern != "" {
			attributes += " pattern=" + addQuotes(input.pattern)
		}
		if allowedAttrs.placeholder != "" && input.placeholder != "" {
			attributes += " placeholder=" + addQuotes(input.placeholder)
		}
		if allowedAttrs.required && input.required {
			attributes += " required"
		}
		if allowedAttrs.size > 0 && input.size > 0 {
			if input.html == "select" && len(input.options) > 0 {
				attributes += fmt.Sprintf(" size=%v", min(len(input.options), 12))
			} else if input.html != "select" {
				attributes += fmt.Sprintf(" size=%v", input.size)
			}
		}
		if allowedAttrs.step > 0 && input.step > 0 {
			//Trim leading zeros & FormatFloat removes any trailing decimal places. e.g. 1.000000 = "1", 0.001000 = ".001".
			attributes += " step=" + strings.TrimLeft(strconv.FormatFloat(input.step, 'f', -1, 32), "0")
		}
		if allowedAttrs.value != nil && input.value != nil {
			value := fmt.Sprintf("%v", input.value)
			if value != "" {
				attributes += " value=" + addQuotes(value)
			}
		}

		if allowedAttrs.html != "" || input.html == "text" {
			element += fmt.Sprintf("<input%v>%v", attributes, options)
		} else if input.html == "select" {
			element += fmt.Sprintf("<select%v>%v</select>", attributes, options)
		} else if input.html == "submit" {
			element += fmt.Sprintf("<button%v>%v</button>", attributes, input.inner)
		} else if input.html != "" {
			element += fmt.Sprintf("<%v%v>%v</%v>", input.html, attributes, input.inner, input.html)
		}
		if input.help != "" {
			element += "<abbr class=help title=\"" + input.help + "\">?</abbr>"
		}
		if input.label != "" {
			element = "<label>" + input.label + element + "</label>"
		}
		formElements = append(formElements, element)
		attributes = ""
		options = ""
		element = ""
	}

	if form.id != "" {
		formID = " id=" + form.id
	}
	if form.table {
		formElements[0] = fmt.Sprintf("<form%v action=%v method=post>%v</form>", formID, addQuotes(form.action), formElements[0])
		return "<tr><td>" + strings.Join(formElements, "<td>")
	}
	return fmt.Sprintf("<form%v action=%v method=post>%v%v</fieldset></form>", formID, addQuotes(form.action), fieldSet(form.title), strings.Join(formElements, ""))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func drawOnlyOptions(options []Option) string {
	var output, optionValue string
	for _, option := range options {
		optionValue = fmt.Sprintf("%v", option.Value)
		if optionValue == "" && option.Display == "" {
			info.Println("options contains an empty option")
		}
		output += "<option"
		if option.Selected {
			output += " selected"
		}
		if optionValue != option.Display {
			output += " value=" + addQuotes(optionValue)
		} else {
			info.Println("option.Value shouldn't be included when it has the save value as option.Display")
		}
		output += ">" + option.Display
	}
	return output
}

/*func generateForm(form Form) string {
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
			if input.autofocus {
				attributes += " autofocus"
			}
			if input.accessKey != "" {
				attributes += " accesskey=" + input.accessKey
			}
			if input.size > 0 {
				attributes += fmt.Sprintf(" size=%v", input.size)
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
					attributes += fmt.Sprintf(" size=%v", len(input.options))
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
				output += "<button" + attributes + ">" + input.inner + "</button>"
			} else {
				if input.html != "text" {
					attributes += " type=" + input.html
				}
				if input.dataList && options != "" {
					if input.id == "" {
						warning.Println("datalist needs a unique ID")
					}
					if input.html == "text" {
						warning.Println("datalist type should be search")
					}
					attributes += " list=" + input.id
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
		output = "<tr><td>" + strings.Join(formElements, "<td>")
		return output
	}

	output = strings.Join(formElements, " ")
	if form.title != "" {
		output = fieldSet(form.title) + output + "</fieldset>"
	} else {
		//devModeCheckForm(false, "all forms should have a title")
	}
	return fmt.Sprintf("<form%v action=%v method=post>%v</form>", formID, addQuotes(form.action), output)
}*/

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
		output += ">" + option.Display
		//devModeCheckForm(!(option.Display == "" && option.Value == "" && option.Selected == false), "option must have display text")
	}
	if input.dataList {
		output = "<datalist id=" + input.id + ">" + output + "</datalist>"
	} else if input.placeholder != "" {
		output += "<option selected value disabled>" + input.placeholder
	}
	return output
}

func fieldSet(title string) string {
	return fmt.Sprintf("<fieldset><legend>%v</legend>", title)
}
