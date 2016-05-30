'use strict';
var editRows = document.querySelectorAll('#editRow tbody tr td:last-of-type'), i = editRows.length;
while(i--){
	editRows[i].onclick = editRow;
}
function editRow(editCell){
	var row = editCell.target.parentElement,
		tds = row.children,
		t = document.importNode(document.querySelector('template').content, true),
		form = t.querySelector('form'),
		outerFormFields = t.querySelectorAll('[form=editRow]'),
		i = outerFormFields.length;
	form.id = '_' + tds[0].textContent;
	while(i--){
		outerFormFields[i].setAttribute('form', form.id);
	}
	t.querySelector('td').textContent = tds[0].textContent;
	if(buildRow){
		t = buildRow(t, tds, form.id);
		row.innerHTML = '';
		row.appendChild(t);
	}
}
function findValues(element, labels){
	var i = labels.length;
	while(i--){
		findValue(element, labels[i]);
	}
}
function findValue(element, label){
	var i = element.options.length;
	while(i--){
		//Select the option if its text = the label. Trim isn't required because the server outputs html without any wrapping whitespace.
		if(label === element.options[i].textContent){
			//Assign the selected option and exit the loop
			return element.options[i].setAttribute('selected', '');
		}
	}
}
function checked(cell, element){
	if(cell.className === '^tick^'){
		element.setAttribute('checked', '');
	}
}