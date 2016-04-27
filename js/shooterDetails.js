'use strict';
function stuff(row){
	row = row.parentElement;
	var tds = row.children, t = document.importNode(document.querySelector('template').content, true), shooterID = '_' + tds[0].textContent;
	t.querySelector('form').id = shooterID;
	t.querySelector('[name=s]').value = row.children[1].querySelector('span').textContent;
	tds[1].querySelector('span').textContent = '';
	t.querySelector('[name=f]').value = tds[1].textContent.trim();
	t.querySelector('[name=C]').value = tds[2].textContent;
	t.querySelector('[name=g]').value = findValue(t.querySelector('[name=g]'), tds[3].textContent);
	t.querySelector('[name=r]').value = findValue(t.querySelector('[name=r]'), tds[4].textContent);
	t.querySelector('[name=I]').value = tds[0].textContent;
	t.querySelector('td').textContent = tds[0].textContent;
	var outerFormFields = t.querySelectorAll('[form=editShooter]'), index = outerFormFields.length;
	while(index--){
		outerFormFields[index].setAttribute('form', shooterID);
	}
	row.innerHTML = '';
	row.appendChild(t);
}
function findValue(element, label){
	var options = element.options, i = options.length;
	label = label.trim();
	while(i--){
		if(label === options[i].textContent.trim()){
			return options[i].value;
		}
	}
}
