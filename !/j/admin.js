(function(){
	'use strict';
	var browserVersion = window.navigator.appVersion.match(/Chrome\/(\d+)\./), shooterEntryValues = {};
	if(!browserVersion){
		var note = document.createElement('a');
		note.class = 'update';
		note.href= '//google.com/chrome/browser/features.html';
		note.target = '_blank';
		note.textContent = 'EventBucket works best with Google Chrome';
		document.body.insertBefore(note, document.body.childNodes[0]);
	}

	var j;
	if(window.XMLHttpRequest){
		j = new XMLHttpRequest();
	}

	var select = document.querySelector('select#sid');
	if(select){
		select.onchange=function(select){
			return function(){
				shooterSelected(select.value);
			};
		}(select);
	}

	function searchShooter(inputElement){
		if(shooterEntryValues[inputElement.name] !== inputElement.value){
			shooterEntryValues[inputElement.name] = inputElement.value;

//			console.log("it is", shooterEntryValues, shooterEntryValues.length);
			if(shooterEntryValues.first.length + shooterEntryValues.surname.length + shooterEntryValues.club.length >= 2){
//			for(var i= 0; i < 3; i++){
//				console.log("looking");
//				if(shooterEntryValues[i].length >= 2){
//					console.log("searched");
				ajax();
//					return;
//				}
			}
		}
	}

	function shooterSelected(shooterId){
		j.open('POST', '/queryShooterGrade', true);
		j.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
		j.send('shooterid='+shooterId);
		j.onreadystatechange = function(){
			if(j.status === 200){
				document.getElementById('palette').innerHTML = !j.response.length ? 'No grades found.' : j.response;
			}
		};
	}

	var textboxes = document.querySelectorAll('#ShooterEntry input[type=search]'),
		inputChange = function(inputElement){
			return function(){
				searchShooter(inputElement);
			};
		};

	if(textboxes && textboxes.length){
		var i=-1;
		while(++i < textboxes.length){
			shooterEntryValues[textboxes[i].name] = textboxes[i].value;
			textboxes[i].onkeyup = inputChange(textboxes[i]);
		}
	}

	function ajax(){
//		var shots = getAjax();
//		shots = encodeURI(shots).replace(/#/gi, '%23');	//hashes are converted after encodeURI to stop % being converted twice
		j.open('POST', '/queryShooterList', true);// + 'scoreSave=' + classes.eventId + '~' + Id + '~' + shots, true);
		j.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
//		j.send(JSON.stringify(shooterEntryValues));
		j.send('first='+shooterEntryValues.first+'&surname='+shooterEntryValues.surname+'&club='+shooterEntryValues.club);
		j.onreadystatechange = function(){
			if(j.status === 200){
//				if(!j.response.length){
				document.getElementById('sid').innerHTML = !j.response.length ? '<option value>0 shooters found.</option>' : j.response;
				if(j.response.length){
					shooterSelected(document.getElementById('sid').value)
				}
//				}else{
//					document.getElementById('sid').innerHTML = j.response;
//				}
			}
		};
	}

	function changeRequiredAttrs(formAction){
		var inputs = document.querySelectorAll('#ShooterEntry [type=search]');
		if(formAction){
			document.getElementById('sid').setAttribute('required', '');
		}else{
			document.getElementById('sid').removeAttribute('required');
		}
		if(input && inputs.length){
			var i = inputs.length;
			while(--i){
				if(formAction){
					inputs[i].removeAttribute('required');
				}else{
					inputs[i].setAttribute('required', '');
				}
			}
		}
	}

	var shooterButtons = document.querySelector('#ShooterEntry button'),
		clickIt = function(formAction){
			return function(){
				changeRequiredAttrs(formAction);
			};
		};
	if(shooterButtons && shooterButtons.length){
		var i = shooterButtons.length;
		while(--i){
			shooterButtons[i].onclick = clickIt(shooterButtons[i].getAttribute('formaction'));
		}
	}

	var addExistingShooter = document.getElementById('addExistingShooter');
	if (addExistingShooter){
		addExistingShooter.onclick = function(button){
			return function(){
				button.form.first.required = false;
				button.form.surname.required = false;
				button.form.club.required = false;
				button.form.sid.required = true;
			};
		}(addExistingShooter);
	}

	var addNewShooter = document.getElementById('addNewShooter');
	if (addNewShooter){
		addNewShooter.onclick = function(button){
			return function(){
				button.form.first.required = true;
				button.form.surname.required = true;
				button.form.club.required = true;
				button.form.sid.required = false;
			};
		}(addNewShooter);
	}

//Disable submit buttons on form submit to prevent double form submits
	var disableButtons = function (form) {
			var b = form.querySelectorAll("button"), g = b.length;
			while(g--){  //post decrement because we need 0
				b[g].disabled = true;
			}
		},
		f = document.forms.length;
	while(f--){  //post decrement because we need 0
		document.forms[f].onsubmit = disableButtons;
	}
}());

//display errors when these elements are found: , input[type=submit], input[type=reset]