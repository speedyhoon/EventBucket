var xhr = new XMLHttpRequest, f = '', s = '', c = document.querySelector('[action="/form.action.eventShooterNew"] [name=C]').value;
document.querySelector('form[action="/form.action.eventShooterNew"]').oninput = function search(event){
	var form = event.currentTarget;
	//Ignore form inputs into other fields
	if(f === form.f.value && s === form.s.value && c === form.C.value)return;
	f = form.f.value;
	s = form.s.value;
	c = form.C.value;
	xhr.onreadystatechange = function(){
		if(xhr.readyState == 4){
			form.S.removeAttribute('class');
			if(xhr.status == 200 && xhr.responseText.length){
				form.S.innerHTML = xhr.responseText;
				form.S.parentElement.removeAttribute('hidden');
			}else{
				form.S.parentElement.setAttribute('hidden', '');
			}
		}
	};
	//TODO fix loading class
	form.S.setAttribute('class', '^loading^');
	xhr.open('GET', '/form.action.shooterSearch' + (f ? '?f=' + f : '') + (s ? '?s=' + s : '') + (c ? '?c=' + c : ''), true);
	xhr.send();
};