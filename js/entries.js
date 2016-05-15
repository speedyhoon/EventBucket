'use strict';
var editRows = document.querySelectorAll('#editRow tbody tr td:last-of-type'), qty = editRows.length;
while(qty--){
	editRows[qty].onclick = stuff;
}
function stuff(editCell){
	var row = editCell.target.parentElement;
	var tds = row.children, t = document.importNode(document.querySelector('template').content, true), form = t.querySelector('form');
	form.id = '_' + tds[0].textContent;
	var outerFormFields = t.querySelectorAll('[form=editRow]'), index = outerFormFields.length;
	while(index--){
		outerFormFields[index].setAttribute('form', form.id);
	}
	t.querySelector('td').textContent = tds[0].textContent;

	t = custom(t, tds);
	row.innerHTML = '';
	row.appendChild(t);
}
function findValues(element, labels){
	var i = labels.length, j = element.options.length;
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
	if(cell.className == '^tick^'){
		element.setAttribute('checked', '');
	}
}

function custom(t, tds){
	t.querySelector('[name=s]').value = tds[1].querySelector('span').textContent;
	tds[1].querySelector('span').textContent = '';
	t.querySelector('[name=f]').value = tds[1].textContent.trim();
	t.querySelector('[name=C]').value = tds[2].textContent;
	findValues(t.querySelector('[name=g]'), tds[3].textContent);
	findValue(t.querySelector('[name=r]'), tds[4].textContent);
	t.querySelector('[name=I]').value = tds[0].textContent;
	t.querySelector('td').textContent = tds[0].textContent;
	checked(tds[5], t.querySelector('[name=k]'));
	checked(tds[6], t.querySelector('[name=x]'));
	return t;
}