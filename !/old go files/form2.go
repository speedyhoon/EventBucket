package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Form is exported
type Form2 struct {
	action string
	title  string
	fields []Field
	help   string
	table  bool
	id     string
}

type Field struct {
	//	min, max, step    float64
	min, max          *int
	step              float64
	isValid, required bool
	minLen, maxLen    int
	value             interface{}
	kind              interface{}
	name, errMsg      string
	err               error

	//used html instead of type because type is a language keyword
	html, label, help, pattern, placeholder, autoComplete string //AutoComplete values can be: "on" or "off"
	checked, multiSelect                                  bool

	size    int
	options []Option

	varType              string //the type of variable to return
	maxLength, minLength int    //the length of variable to return
	error                string
	snippet              interface{}
	autofocus            bool
	action               string //Way to switch the forms action to a different purpose

	accessKey, inner, id string
	dataList             bool
	class                string
}

/*

var on = 1
var formAttrs = map[string]Field{
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
*/

func makeForm(form Form2) string {
	if conn == nil {
		return "<p class=error>Unable to connect to the EventBucket database.</p>"
	}
	var formElements, labelClasses []string
	var formID, attributes, element, options string
	for i, field := range form.fields {
		labelClasses = []string{}
		allowedAttrs := formAttrs[field.html]
		if field.snippet != nil {
			element = fmt.Sprintf("%v", field.snippet)
		}
		if allowedAttrs.html != "" {
			attributes += " type=" + field.html
		}
		if allowedAttrs.name != "" && field.name != "" {
			attributes += " name=" + field.name
		}
		if field.accessKey != "" {
			attributes += " accesskey=" + field.accessKey
		}
		if allowedAttrs.autoComplete != "" && field.autoComplete == "on" || field.autoComplete == "off" {
			attributes += " autocomplete=" + field.autoComplete
		}
		if allowedAttrs.autofocus && field.autofocus {
			attributes += " autofocus"
		}
		if form.table && form.id != "" && i > 0 {
			attributes += " form=" + form.id
		}
		if allowedAttrs.checked && field.checked {
			attributes += " checked"
		}
		if allowedAttrs.action != "" && field.action != "" {
			attributes += " formaction=" + field.action
		}
		if field.class != "" {
			attributes += " class=" + field.class
		}

		/*var elementClass []string
		if field.class != "" {
			elementClass = append(elementClass, field.class)
		}
		if field.help != "" && field.class != "il" {
			elementClass = append(elementClass, "il")
		}
		if len(elementClass) >= 1 {
			attributes += " class=" + strings.Join(elementClass, " ")
		}*/

		if allowedAttrs.dataList && field.dataList && field.id != "" {
			attributes += " list=" + field.id
			if len(field.options) > 0 {
				options = "<datalist id=" + field.id + ">" + drawOnlyOptions(field.options) + "</datalist>"
			}
		}
		if allowedAttrs.max != nil && field.max != nil {
			attributes += fmt.Sprintf(" max=%v", *field.max)
		}
		if allowedAttrs.maxLength > 0 && field.maxLength > 0 {
			attributes += fmt.Sprintf(" maxlength=%v", field.maxLength)
		}
		if allowedAttrs.min != nil && field.min != nil {
			attributes += fmt.Sprintf(" min=%v", *field.min)
		}
		if allowedAttrs.minLength > 0 && ((field.required && field.minLength > 1) || (!field.required && field.minLength > 0)) {
			attributes += fmt.Sprintf(" minlength=%v", field.minLength)
		}
		if allowedAttrs.multiSelect && field.multiSelect {
			attributes += " multiple"
			if field.html == "select" && len(field.options) > 0 {
				attributes += fmt.Sprintf(" size=%v", min(len(field.options), 12))
			}
		}
		if field.html == "select" {
			options = drawOnlyOptions(field.options)
		}
		if allowedAttrs.pattern != "" && field.pattern != "" {
			attributes += " pattern=" + addQuotes(field.pattern)
		}
		if allowedAttrs.placeholder != "" && field.placeholder != "" {
			attributes += " placeholder=" + addQuotes(field.placeholder)
		}
		//		if allowedAttrs.required && field.required {
		//			attributes += " required"
		//		}
		if allowedAttrs.size > 0 && field.size > 0 {
			if field.html == "select" && len(field.options) > 0 {
				attributes += fmt.Sprintf(" size=%v", min(len(field.options), 12))
			} else if field.html != "select" {
				attributes += fmt.Sprintf(" size=%v", field.size)
			}
		}
		if allowedAttrs.step > 0 && field.step > 0 {
			//Trim leading zeros & FormatFloat removes any trailing decimal places. e.g. 1.000000 = "1", 0.001000 = ".001".
			attributes += " step=" + strings.TrimLeft(strconv.FormatFloat(field.step, 'f', -1, 32), "0")
		}
		if allowedAttrs.value != nil && field.value != nil {
			value := fmt.Sprintf("%v", field.value)
			if value != "" {
				attributes += " value=" + addQuotes(value)
			}
		}

		/*var hasHelp string
		if field.help != "" {
			hasHelp = "<br>"
		}*/
		if allowedAttrs.html != "" || field.html == "text" {
			element += fmt.Sprintf("<input%v>%v", attributes, options)
		} else if field.html == "select" {
			element += fmt.Sprintf("<select%v>%v</select>", attributes, options)
		} else if field.html == "submit" {
			element += fmt.Sprintf("<button%v>%v</button>", attributes, field.inner)
		} else if field.html != "" {
			element += fmt.Sprintf("<%v%v>%v</%v>", field.html, attributes, field.inner, field.html)
		}
		//		var labelClasses []string
		var elementHelp string
		if field.help != "" {
			//			labelClasses = []string{"help"}
			//			elementHelp = "<img src=/p/help2 alt=? title='" + field.help + "'>"
			//			elementHelp = "<c title='" + field.help + "'>?<div>" + field.help + "</div></c>" //U+2754: WHITE QUESTION MARK ORNAMENT
			elementHelp = "<c title='" + field.help + "'>&#10068</c>" //U+2754: WHITE QUESTION MARK ORNAMENT
		}
		if field.label != "" {

			//			if field.label != "" && !form.table {
			//				var errorClass string
			//				if field.error != "" {
			//					errorClass = " class=error"
			//					field.error = " " + field.error
			//				}
			//				output += "<label" + errorClass + ">" + field.label + ": " + element + field.error + "</label>"

			if field.error != "" {
				labelClasses = append(labelClasses, "error")
				element = "<label class=" + addQuotes(strings.Join(labelClasses, " ")) + ">" + field.label + element + "<p>" + field.error + "</p></label>"
				//				element = "<label class=error>" + field.label + elementHelp + element + "<p>" + field.error + "</p></label>"
			} else {
				//				class := addQuotes(strings.Join(labelClasses, " "))
				//				if class != "" {
				//					class = " class=" + class
				//				}
				//				element = "<label" + class + ">" + field.label + element + "</label>"
				element = "<label>" + field.label + elementHelp + element + "</label>"
			}
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

//
//func min(a, b int) int {
//	if a < b {
//		return a
//	}
//	return b
//}

//func drawOnlyOptions(options []Option) string {
//	var output, optionValue string
//	for _, option := range options {
//		optionValue = fmt.Sprintf("%v", option.Value)
//		if optionValue != "" || option.Display != "" { //Exclude options with only one populated
//			output += "<option"
//			if option.Selected {
//				output += " selected"
//			}
//			if optionValue != option.Display {
//				output += " value=" + addQuotes(optionValue)
//			} else {
//				//				trace.Printf("option.Value shouldn't be included when it has the save value as option.Display %v\n%v", options, optionValue)
//			}
//			output += ">" + option.Display
//		} else {
//			//			trace.Println("options contains an empty option")
//			fmt.Println(options)
//		}
//	}
//	return output
//}

/*func drawOptions(field Inputs) string {
	//devModeCheckForm(len(field.options) > 0, "select should have at least one option to select from for element='"+field.name+"' type='"+field.html+"'")
	if field.required {
		//devModeCheckForm(len(field.options) > 0, "select shouldn't be required with no available options to select")
	}
	output := ""
	var optionValue string
	for _, option := range field.options {
		output += "<option"
		if option.Selected {
			output += " selected"
			//devModeCheckForm(field.html != "datalist", "datalist shouldn't have any selected values! change it to a value attribute")
			//devModeCheckForm(!(field.placeholder != "" && field.html != "datalist"), "shouldn't set a placeholder when options are already selected")
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
	if field.dataList {
		output = "<datalist id=" + field.id + ">" + output + "</datalist>"
	} else if field.placeholder != "" {
		output += "<option selected value disabled>" + field.placeholder
	}
	return output
}

func fieldSet(title string) string {
	return fmt.Sprintf("<fieldset><legend>%v</legend>", title)
}
*/

//func isValidInt(strNum string, field Inputs) (int, error) {
//	num, err := strconv.Atoi(strNum)
//	if err != nil {
//		return num, err
//	}
//	if num >= int(field.min) && num <= int(field.max) && (field.step == 0 || field.step != 0 && num%int(field.step) == 0) || !field.required && num == 0 {
//		return num, nil
//	}
//	return num, errors.New("field integer doesn't pass validation")
//}

func isValidStr(str string, field Field) (string, error) {
	length := len(str)
	if length >= field.minLen && length <= field.maxLen || !field.required && str == "" {
		return str, nil
	}
	return str, errors.New("field string doesn't pass validation")
}

func isValid(fields []Field, r *http.Request) ([]Field, bool) {
	r.ParseForm()
	if len(r.Form) < 1 {
		return fields, false
	}
	//Process the post request as normal if len(r.Form) > len(fields)
	//	var fieldValue []string
	var ok bool
	//	var err error
	valid := true
	for _, field := range fields {
		_, ok = r.Form[field.name]
		if !ok {
			valid = false
			continue //Skip to the next loop iteration.
		}

		/*switch field.kind.(type) {
		case bool:
			fmt.Printf("boolean %t\n", field.kind)
			//		case int:
			//			fields[i].value, err = isValidInt(strings.TrimSpace(fieldValue[0]), field)
			//			if err == nil {
			//				fields[i].isValid = true
			//			} else {
			//				valid = false
			//				fields[i].errMsg = "integer supplied was wrong"
			//				fmt.Println(field.name + "integer supplied was wrong")
			//			}
		case string:
			fields[i].value, fields[i].err = isValidStr(strings.TrimSpace(fieldValue[0]), field)
			fields[i].isValid = fields[i].err == nil
			fields[i].errMsg = fmt.Sprintf("%v", fields[i].err)
			if !fields[i].isValid {
				valid = false
				fields[i].errMsg = "string doesn't match"
				fmt.Println(field.name + "string doesn't match")
			}
		case []string:
			fmt.Println("within string slice")
			for key, thingy := range fieldValue {
				fmt.Println(key, thingy)
			}
			fields[i].value = fieldValue
		default:
			fmt.Printf("unexpected type %T", field.kind) // %T prints whatever type t is
		}*/
	}
	//	if valid{
	//
	//	}
	return fields, valid
}

//func addQuotes(in string) string {
//	return in
//}
//
//func fieldSet(in string) string {
//	return in
//}
//
//// Option is exported
//type Option struct {
//	Value    interface{} `json:"v,omitempty"`
//	Display  string      `json:"d,omitempty"`
//	Selected bool        `json:"s,omitempty"`
//}
