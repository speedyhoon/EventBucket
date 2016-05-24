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

	t = entries(t, tds);
	row.innerHTML = '';
	row.appendChild(t);
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

function entries(t, tds){
	var surname = tds[1].querySelector('span');
	t.querySelector('[name=s]').value = surname.textContent;
	surname.textContent = '';
	t.querySelector('[name=f]').value = tds[1].textContent.trim();
	t.querySelector('[name=C]').value = tds[2].textContent;
	findValue(t.querySelector('[name=g]'), tds[3].textContent);
	findValue(t.querySelector('[name=r]'), tds[4].textContent);
	t.querySelector('[name=I]').value = tds[0].textContent;
	checked(tds[5], t.querySelector('[name=k]'));
	checked(tds[6], t.querySelector('[name=x]'));
	return t;
}
var button = document.querySelectorAll('[action="/8"] button'), i=button.length;
while(i--){
	button[i].onclick = function shooterEntry(){
		var isNew = !this.getAttribute('formAction');
		this.form.f.required = isNew;
		this.form.s.required = isNew;
		this.form.S.required = !isNew;
		this.form.C.required = isNew;
	};
}