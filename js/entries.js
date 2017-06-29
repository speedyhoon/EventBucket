var button = document.querySelectorAll('[action="/8"] button'), i=button.length;
while(i--){
	button[i].onclick = function shooterEntry(){
		var isNew = !this.getAttribute('formAction');
		this.form.f.required = isNew;
		this.form.s.required = isNew;
		this.form.S.required = !isNew;
		this.form.C.required = isNew;
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