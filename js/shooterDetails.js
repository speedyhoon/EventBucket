function buildRow(t, tds){
	var surname = tds[1].querySelector('span');
	t.querySelector('[name=s]').value = surname.textContent;
	surname.textContent = '';
	t.querySelector('[name=f]').value = tds[1].textContent.trim();
	t.querySelector('[name=C]').value = tds[2].textContent;
	findValues(t.querySelector('[name=g]'), tds[3].textContent.replace(/, $/, '').split(', '));
	findValue(t.querySelector('[name=r]'), tds[4].textContent);
	t.querySelector('[name=I]').value = tds[0].textContent;
	t.querySelector('td').textContent = tds[0].textContent;
	checked(tds[5], t.querySelector('[name=x]'));
	return t;
}