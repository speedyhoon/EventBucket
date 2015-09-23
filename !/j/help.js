'use strict';
var c = document.getElementsByTagName('c'), index = c.length;
while(--index){
	c[index].onclick = function (text) {
		return function () {
			alert(text);
		}
	}(c[index].title);
}


















clearInterval(intervalId);
var intervalId = setInterval(function(){
	var selectors = ['a', 'form'],
		selectorsQty = selectors.length,
		index = Math.floor(Math.random() * selectorsQty) - 1,
		elements = document.querySelectorAll(selectors[index]),
		qty = elements.length,
		element = elements[Math.floor(Math.random() * qty)];
	if(element.nodeName === 'form'){
		submitForm(element);
	}else{
		element.click();
	}
}, 1000);

//var specialForms = {
//	shooterInsert: function(){return 6}
//};

function submitForm(form){
	var action = form.getAttribute('action'),
		form_backup = p.cloneNode(true);
	clearValidation(form);
//	for (var property in specialForms) {
//		if (specialForms.hasOwnProperty(property) && action === '/' + property) {
//			return specialForms[property]();
//		}
//	}
//	standardFormSubmit(form)
	form.submit();
}
//function standardFormSubmit(form){
//}
function clearValidation(form){
	var qty = form.elements.length;
	while(qty--){
		form.elements[qty].removeAttribute('max');
		form.elements[qty].removeAttribute('maxlength');
		form.elements[qty].removeAttribute('min');
		form.elements[qty].removeAttribute('minlength');
		form.elements[qty].removeAttribute('step');
		form.elements[qty].removeAttribute('required');
		form.elements[qty].removeAttribute('pattern');
	}
}

/*
hard code server check to 10.1.1.1:80
websocket listener
add small js file to EventBucket
on js load connect to test server
	send back document.location.href
server writes location to log file
js randomly selects between link or form
	js sends back link go to location if it picked link
	OR
	js removes validation & submits form & sends server a copy of the form data and original form
js asserts that every form element has maxlength
*/