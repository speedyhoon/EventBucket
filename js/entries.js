var button = document.querySelectorAll('[action="/form.eventShooterNew.action"] button'), i=button.length;
while(i--){
	button[i].onclick = function shooterEntry(event){
		var isNew = !event.target.getAttribute('formAction');
		event.target.form.f.required = isNew;
		event.target.form.s.required = isNew;
		event.target.form.S.required = !isNew;
		event.target.form.C.required = isNew;
	};
}
window['buildRow'] = function(tr, tds){
	var surname = tds[1].querySelector('span');
	tr.querySelector('[name=s]').value = surname.textContent;
	surname.textContent = '';
	tr.querySelector('[name=f]').value = tds[1].textContent.trim();
	tr.querySelector('[name=C]').value = tds[2].textContent;
	findValue(tr.querySelector('[name=g]'), tds[3].textContent);
	findValue(tr.querySelector('[name=r]'), tds[4].textContent);
	tr.querySelector('[name=S]').value = tds[0].textContent;
	checked(tds[5], tr.querySelector('[name=k]'));
	checked(tds[6], tr.querySelector('[name=x]'));
	return tr;
};