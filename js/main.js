'use strict';
//Prevent users from double/triple clicking on submit buttons that could cause a form to submit more than once with exactly the same data. If forms used a CSRF token this JavaScript would be redundant
function stopDoubleSubmit(){
	if(submitted)return false;
	return submitted = true;
}
var submitted = false;
for(var forms = document.forms, n = forms.length; n--;){
	forms[n].onsubmit = stopDoubleSubmit;
}
function tableSort(th){
	th = th.target;
	/*If th.textContent == id compare using base36 else use textContent */
	var tbody = th.parentElement.parentElement.parentElement.querySelector('tbody'), column = Array.prototype.indexOf.call(th.parentElement.children, th), direction = th.className === '^asc^' ? -1 : 1, rows = Array.from(tbody.children);
	//Array.prototype.sort.call(tbody.children, function (a, b) {
	rows.sort(function(a, b){
		a = a.children[column].textContent;
		b = b.children[column].textContent;
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
/*TODO doesn't submit the highlighted values. Actually submits the previous selected fields. Also bypasses form validation
//Add a shortcut to select boxes that allow their parent form to submit when Ctrl + Enter is pressed.
var selectBoxes = document.querySelectorAll('select'), qty = selectBoxes.length, formShortcut = function(event){
		if(event.keyCode === 13 && event.ctrlKey && event.target.form){
		event.target.form.submit();
		return false;
	}
	};
while(qty--){
	selectBoxes[qty].onkeydown = formShortcut;
}*/



var form;
function initForm(){
//	form = document.getElementById("usrForm");
	form = document.forms[0];
	form.addEventListener("submit", function(evt){
		if(form.checkValidity() === false){
			evt.preventDefault();
			alert("Form is invalid - submission prevented!");
			return false;
		}else{
			// To prevent data from being sent, we've prevented submission
			// here, but normally this code block would not exist.
			evt.preventDefault();
			alert("Form is valid - submission prevented to protect privacy.");
			return false;
		}
	});
}
/*function initConfirmEmail(){
	var elem = document.getElementById("frmEmailC");
	elem.addEventListener("blur", verifyEmail);
	function verifyEmail(input){
		input = input.srcElement;
		sampleCompleted("Forms-orderConfirm");
		var primaryEmail = document.getElementById('frmEmailA').value
		if(input.value != primaryEmail){
			// the provided value doesn't match the primary email address
			input.setCustomValidity('The two email addresses must match.');
			console.log("E-mail addresses do not match", primaryEmail, input.value);
		}else{
			// input is valid -- reset the error message
			input.setCustomValidity('');
		}
	}
}*/
function initInputs(){
	var inputs = document.getElementsByTagName("input");
	var inputs_len = inputs.length;
	var addDirtyClass = function(evt){
		sampleCompleted("Forms-order-dirty");
		evt.srcElement.classList.toggle("dirty", true);
	};
	for(var i = 0; i < inputs_len; i++){
		var input = inputs[i];
		input.addEventListener("blur", addDirtyClass);
		input.addEventListener("invalid", addDirtyClass);
		input.addEventListener("valid", addDirtyClass);
	}
}
/*function initNoSubmit(){
	form.addEventListener("submit", function(evt){
		evt.preventDefault();
		alert("Submission of this form is prevented.");
	});
}*/
//initForm();
initInputs();
//initConfirmEmail();

var isCompleted = {};
function sampleCompleted(sampleName){
	if (!isCompleted.hasOwnProperty(sampleName)) {
		isCompleted[sampleName] = true;
	}
}