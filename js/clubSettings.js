'use strict';
var editShooters = document.querySelectorAll('#editMounds tbody tr td:last-of-type'), qty = editShooters.length;
while(qty--){
	editShooters[qty].onclick = stuff;
}
function stuff(editCell){
	var row = editCell.target.parentElement;
	var tds = row.children, t = document.importNode(document.querySelector('template').content, true), ID = row.id;
	t.querySelector('form').id = '_'+ID;
	t.querySelector('[name=n]').value = row.children[0].textContent;
	t.querySelector('[name=I]').value = ID;
	var outerFormFields = t.querySelectorAll('[form=edit]'), index = outerFormFields.length;
	while(index--){
		outerFormFields[index].setAttribute('form', '_'+ID);
	}
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
