'use strict';
var pathName = window.location.pathname.split('/')[1],
	eventID = window.location.pathname.split('/')[2],
	rangeID = window.location.pathname.split('/')[3],
	inputs = document.querySelectorAll('table input'),
	i = inputs.length,
	ws, intervalId; //intervalId global variable stops reconnect() snowballing into an infinite loop.
while(i--){
	inputs[i].onchange = save;
}

//Set shooter barcode form onsubmit because it's not allowed with the current Content Security Policy.
document.querySelector('#sb').onsubmit='return shooterBarcode(B)';

function reconnect (){
	ws = new WebSocket('ws://'+window.location.host+'/w/');
	ws.onopen = function(){
		if(intervalId){
			clearInterval(intervalId);
			intervalId = undefined;
		}
	};
	//TODO
	//Update UI with save / error message.
	//ws.onmessage = function(message){};
	ws.onclose = function(){
		if(!intervalId){
			intervalId = setInterval(reconnect, 3000); //try to reconnect every 3 seconds after the connection is dropped.
		}
	};
}
reconnect();

//If shooter ID is provided, try to find the shooter. Otherwise autofocus on barcode search textbox
if(!window.location.hash){
	document.querySelector('[name=B]').setAttribute('autofocus', '');
}else{
	goToShooter(window.location.hash.replace('#', ''), document.querySelector('[name=B]'));
}

//Search for a shooter on Scorecards & total-scores pages.
function shooterBarcode(search){
	document.getElementById('searchErr').setAttribute('hidden', '');
	document.getElementById('barcodeErr').setAttribute('hidden', '');
	document.getElementById('shooterErr').setAttribute('hidden', '');
	if(!search||!search.value){
		document.getElementById('searchErr').removeAttribute('hidden');
		return false;
	}
	if(/^\d+$/g.test(search.value)){
		goToShooter(search.value, search);
		return false;
	}
	//If barcode doesn't match display error message
	if(!/^\d+\/\d+#\d+$/g.test(search.value)){
		document.getElementById('barcodeErr').removeAttribute('hidden');
		search.select();
		return false;
	}
	var barcodeEventID = search.value.split('/')[0],
		barcodeRangeID = search.value.split('/')[1].split('#')[0],
		shooterID = search.value.split('#')[1];
	if(eventID !== barcodeEventID){
		//Go to a different event if user presses OK.
		if(confirm('This barcode is for a different event. Do you want to go to event with id '+barcodeEventID+'?')){
			window.location.href = '/' + pathName + '/' + barcodeEventID + '/' + barcodeRangeID + '#' + shooterID;
		}
		//Else do nothing.
		search.select();
		return false;
	}
	if(rangeID !== barcodeRangeID){
		//Go to a different range if user presses OK.
		if(confirm('This barcode is for a different range. Do you want to go to range with id ' + barcodeRangeID + '?')){
			window.location.href = '/' + pathName + '/' + barcodeEventID + '/' + barcodeRangeID + '#' + shooterID;
		}
		//Else do nothing.
		search.select();
		return false;
	}
	return goToShooter(shooterID, search);
}

function goToShooter(shooterID, search){
	//If the shooter textbox exists in the DOM, set focus to their text box.
	var d = document.getElementById(shooterID);
	if(d){
		search.value = '';
		window.location.hash = '#' + shooterID;
		d.select();
		return false;
	}
	if(pathName.indexOf('-all') >= 0){
		//Display shooter not is this event error message.
		document.getElementById('shooterErr').removeAttribute('hidden');
		return false;
	}
	//If the shooter doesn't exist go to the scorecards-all OR total-scores-all page
	window.location.href = '/' + pathName + '-all/' + eventID + '/' + rangeID + '#' + shooterID;
	return false;
}

function save(event){
	var row = event.target.parentElement.parentElement,
		name = event.target.name,
		otherInput = name==='t'?'c':'t';
	//Assigning values as arrays so json.Marshal can convert it to url.Values straight away & doesn't require custom validation code
	var score = {
		E: [eventID],
		R: [rangeID],
		S: [row.children[0].textContent]
	};
	//Use double bitwise operator to convert strings to an integer
	score[name] = [~~event.target.value+''];
	score[otherInput]= [~~row.querySelector('input[name='+otherInput+']').value+''];
	if(validateScore(score)){
		ws.send('\u000E' + JSON.stringify(score));
	}else{
		//Display validation error
		console.warn('Score is invalid: total=',score.t[0], 'centers=', score.c[0]);
	}
}

function validateScore(score){
	//TODO dynamicly change 5 to 6
	if(score.c[0] > ~~(score.t[0]/5)){
		return false;
	}
	return true;
}