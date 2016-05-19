'use strict';
//Prevent users from double/triple clicking on submit buttons that could cause a form to submit more than once with exactly the same data. If forms used a CSRF token this JavaScript would be redundant
var submitted = false;
function stopDoubleSubmit(){
	if(submitted)return false;
	return submitted = true;
}

for(var n = document.forms.length; n--;){
	if(!document.forms[n].onsubmit){
		document.forms[n].onsubmit = stopDoubleSubmit;
	}
}
function tableSort(th){
	th = th.target;
	/*If th.textContent == id compare using base36 else use textContent */
	var tbody = th.parentElement.parentElement.parentElement.querySelector('tbody'), column = Array.prototype.indexOf.call(th.parentElement.children, th), direction = th.className === '^asc^' ? -1 : 1, rows = Array.from(tbody.children);
	var sortBy = function(input){
		switch(th.textContent){
			case 'ID':
				return +input;
			case 'Id':
				return parseInt(input, 36);   //Id = base 36 string (0-9a-z) e.g. a2e = 13046
		}
		return input;
	};
	rows.sort(function(a, b){
		a = sortBy(a.children[column].textContent);
		b = sortBy(b.children[column].textContent);
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
var dialog = document.createElement('dialog');
dialog.onclick = dialog.close;
document.body.appendChild(dialog);
//Add table sort to TH elements within a THEAD
document.onclick = function(event){
	//Help message popup
	if(event.target.nodeName === 'ABBR'){
		dialog.textContent = event.target.title;
		return dialog.showModal();
	}
	//Table sort
	if(event.target.nodeName === 'TH' && event.target.parentElement.parentElement.nodeName === 'THEAD'){
		tableSort(event);
	}
};

function initInputs(){
	var inputs = document.getElementsByTagName('input'), i = inputs.length, flagClass = function(evt){
		evt.srcElement.classList.toggle('^dirty^', true);
	};
	while(i--){
		inputs[i].addEventListener('blur', flagClass);
		inputs[i].addEventListener('invalid', flagClass);
		inputs[i].addEventListener('valid', flagClass);
	}
}
initInputs();


/*
var t = document.querySelectorAll('[required]'), r=t.length;
while(r--){
  t[r].removeAttribute('required');
}
* */