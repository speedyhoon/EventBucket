(function(){
	'use strict';
	var browserVersion = window.navigator.appVersion.match(/Chrome\/(\d+)\./), shooterEntryValues = {};
	if(!browserVersion){
		var note = document.createElement('a');
		note.id = 'update';
		note.href= '//google.com/chrome/browser/features.html';
		note.target = '_blank';
		note.textContent = 'EventBucket works best with Google Chrome';
		document.body.insertBefore(note, document.body.childNodes[0]);
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

	var textboxes = document.getElementById('ShooterEntry'),
		inputChange = function(inputElement){
			return function(){
				searchShooter(inputElement);
			};
		};


	if(textboxes && textboxes.length){
//	if(textboxes && max){
		var i, max = textboxes.length;
		for(i=0; i < max; i++){
//		var i = textboxes.length;
//		while(i--){
			if(textboxes[i].type === 'text'){
				shooterEntryValues[textboxes[i].name] = textboxes[i].value;
				textboxes[i].onkeyup = inputChange(textboxes[i]);
			}
		}
	}

	var j;
	if(window.XMLHttpRequest){
		j = new XMLHttpRequest();
	}
	function ajax(){
//		var shots = getAjax();
//		shots = encodeURI(shots).replace(/#/gi, '%23');	//hashes are converted after encodeURI to stop % being converted twice
		j.open('POST', '/queryShooterList', true);// + 'scoreSave=' + classes.eventId + '~' + Id + '~' + shots, true);
//		j.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
		j.send(JSON.stringify(shooterEntryValues));
		j.onreadystatechange = function(){
			if(j.status === 200){
//				if(!j.response.length){
				document.getElementById('sid').innerHTML = !j.response.length ? '<option value>No shooters found...</option>' : j.response;
//				}else{
//					document.getElementById('sid').innerHTML = j.response;
//				}
			}
		};
	}



	function changeRequiredAttrs(formAction){
		var inputs = document.getElementById('ShooterEntry').getElementsByTagName('input');
		if(max > 0){
			for(i = 0, max = inputs.length; i < max; i++){
				if(inputs[i].type === 'text'){
					if(formAction){
						inputs[i].removeAttribute('required');
						document.getElementById('sid').setAttribute('required', '');
					}else{
						inputs[i].setAttribute('required', '');
						document.getElementById('sid').removeAttribute('required');
					}
				}
			}
		}
	}

	var shooterButtons = document.getElementById('ShooterEntry').getElementsByTagName('button'),
		clickIt = function(formAction){
			return function(){
				changeRequiredAttrs(formAction);
			};
		};
	if(max > 0){
		for(i = 0, max = shooterButtons.length; i < max; i++){
			shooterButtons[i].onclick = clickIt(shooterButtons[i].getAttribute('formaction'));
		}
	}

}());