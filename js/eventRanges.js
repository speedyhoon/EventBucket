'use strict';
function buildRow(t, tds){
	var R = t.querySelector('[name=R]');
	if(tds[2].textContent){
		findValues(R, tds[2].textContent.replace(/, $/, '').split(', '));
		t.querySelector('form').setAttribute('action','/21')
	}else{
		R.parentElement.removeChild(R);
	}
	t.querySelector('[name=n]').value = tds[1].textContent;
	var l = t.querySelector('[name=k]');
	if(!tds[2].textContent){
		if(tds[3].className === '^tick^'){
			l.setAttribute('checked', '');
		}
	}else{
		l.parentElement.removeChild(l);
	}
	t.querySelector('[name=I]').value = tds[0].textContent;
	return t;
}