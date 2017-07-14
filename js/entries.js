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
function buildRow(t, tds){
	var surname = tds[1].querySelector('span');
	t.querySelector('[name=s]').value = surname.textContent;
	surname.textContent = '';
	t.querySelector('[name=f]').value = tds[1].textContent.trim();
	t.querySelector('[name=C]').value = tds[2].textContent;
	findValue(t.querySelector('[name=g]'), tds[3].textContent);
	findValue(t.querySelector('[name=r]'), tds[4].textContent);
	t.querySelector('[name=S]').value = tds[0].textContent;
	checked(tds[5], t.querySelector('[name=k]'));
	checked(tds[6], t.querySelector('[name=x]'));
	return t;
}