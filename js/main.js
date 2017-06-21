'use strict';
//Prevent users from double/triple clicking on submit buttons that could cause a form to submit more than once with exactly the same data. If forms used a CSRF token this JavaScript would be redundant
var submitted = false, n = document.forms.length;
function stopDoubleSubmit(){
	if(submitted)return false;
	return submitted = true;
}

while(n--){
	if(!document.forms[n].onsubmit){
		document.forms[n].onsubmit = stopDoubleSubmit;
	}
}

function tableSort(th){
	if(th.parentElement.parentElement.nodeName !== 'THEAD'){
		return;
	}
	//If th.textContent == id compare using base36 else use textContent
	var tbody = th.parentElement.parentElement.parentElement.querySelector('tbody'),
		column = Array.prototype.indexOf.call(th.parentElement.children, th),
		direction = th.className === '^asc^' ? -1 : 1,
		rows = Array.from(tbody.children);
	var sortBy = function(cell){
		if(!cell){
			return '';
		}
		var input = cell.textContent;
		switch(th.textContent){
			case 'ID':
				return ~~input;
			//Used on the Shooters page
			case 'Id':
				return parseInt(input, 36);   //Id = base 36 string (0-9a-z) e.g. a2e = 13046
		}
		return input;
	};
	rows.sort(function(a, b){
		a = sortBy(a.children[column]);
		b = sortBy(b.children[column]);
		return a > b ? direction : a < b ? -1 * direction : 0;
	});
	tbody.innerHTML = '';
	for(var i = 0, max = rows.length; i < max; i++){
		tbody.appendChild(rows[i]);
	}
	var ths = th.parentElement.querySelectorAll('.^asc^,.^desc^'), qty = ths.length;
	while(qty--){
		ths[qty].removeAttribute('class');
	}
	th.className = direction > 0 ? '^asc^' : '^desc^';
}

//Form help dialog popup
var dialog = document.createElement('dialog'), label;
dialog.onclick = function(){
	dialog.close();
	if(label){
		label.click();
		label = 0;
	}
};
document.body.appendChild(dialog);

//Add table sort to TH elements within a THEAD
document.onclick = function(event){
	var target = event.target;
	switch(target.nodeName){
	//Help message pop-up
	case 'ABBR':
		dialog.textContent = target.title;
		//Assuming the <abbr> element will always be the immediate child of a <label> element
		label = target.parentElement;
		dialog.showModal();
		break;
	//Table sort
	case 'TH':
		tableSort(target);
	}
};

//Start event listeners that add a class to form fields based on their valid/invalid values that changes the background colour.
var inputs = document.querySelectorAll('input'), i = inputs.length, flagClass = function(evt){
		evt.srcElement.classList.toggle('^dirty^', true);
	};
	while(i--){
		inputs[i].addEventListener('blur', flagClass);
		inputs[i].addEventListener('invalid', flagClass);
		inputs[i].addEventListener('valid', flagClass);
}