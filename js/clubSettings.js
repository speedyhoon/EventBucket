'use strict';
function buildRow(t, tds){
	t.querySelector('[name=I]').value = tds[0].textContent;
	t.querySelector('[name=n]').value = tds[1].textContent;
	return t;
}