'use strict';
var editShooters = document.querySelectorAll('#editRow tbody tr td:last-of-type'), qty = editShooters.length;
while(qty--){
	editShooters[qty].onclick = stuff;
}
function stuff(editCell){
	var row = editCell.target.parentElement;
	var tds = row.children, t = document.importNode(document.querySelector('template').content, true), shooterID = '_' + tds[0].textContent;
	t.querySelector('form').id = shooterID;
	t.querySelector('td').textContent = tds[0].textContent;
	t.querySelector('[name=n]').value = tds[1].textContent;
	var R = t.querySelector('[name=R]');
	if(tds[2].textContent){
		findValues(R, tds[2].textContent.replace(/, $/, '').split(', '));
		t.querySelector('form').setAttribute('action','/21')
	}else{
		R.parentElement.removeChild(R);
	}

	var l = t.querySelector('[name=k]');
	if(!tds[2].textContent){
		if(tds[3].className == '^tick^'){
			l.setAttribute('checked', '');
		}
	}else{
		l.parentElement.removeChild(l);
	}
	t.querySelector('[name=I]').value = tds[0].textContent;
	var outerFormFields = t.querySelectorAll('[form=editRow]'), index = outerFormFields.length;
	while(index--){
		outerFormFields[index].setAttribute('form', shooterID);
	}
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
