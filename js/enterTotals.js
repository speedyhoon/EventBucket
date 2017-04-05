'use strict';
var eventID = window.location.pathname.split('/')[2],
	rangeID = window.location.pathname.split('/')[3],
	$inputs = document.querySelectorAll('table input'),
	i = $inputs.length,
	ws, intervalId; //intervalId global variable stops reconnect() snowballing into an infinite loop.
while(i--){
	$inputs[i].onchange = save;
}

function reconnect (){
	ws = new WebSocket('ws://'+window.location.host+'/k/');
	ws.onopen = function(){
		if(intervalId){
			clearInterval(intervalId);
			intervalId = undefined;
		}
	};
	//TODO
	//Update UI with save / error message.
	ws.onmessage = function(message){
		console.log(message.data);
	};
	ws.onclose = function(){
		if(!intervalId){
			intervalId = setInterval(reconnect, 3000); //try to reconnect every 3 seconds after the connection is dropped.
		}
	};
}
reconnect();

function save(event){
	var $row = event.target.parentElement.parentElement,
		name = event.target.name,
		$total = $row.querySelector('input[name=t]'),
		$centre = $row.querySelector('input[name=c]'),
		$shootOff = $row.querySelector('input[name=h]');
	//Assigning values as arrays so Golang json.Marshal can convert it to url.Values straight away & doesn't require custom validation code
	var score = {
		//Using string attribute names here otherwise ClosureCompiler changes them.
		//TODO replace there mapping characters with dynamic build variables?
		'E': [eventID],
		'R': [rangeID],
		'S': [$row.children[0].textContent],
		//Strip any decimal places with double bitwise operator & then convert to string because the backend is expecting a string array.
		't': [~~$total.value+''],
		'c': [~~$centre.value+''],
		'h': [~~$shootOff.value+'']
	};
	if(errorMessage(score, $total, $centre)){
		ws.send('\u000E' + JSON.stringify(score));
	}
}

function errorMessage(score, $total, $centre){
	var popup = document.querySelector('.^popup^'),
		totalMax = ~~$total.max,
		highestPossibleCentres = ~~$centre.getAttribute('data-max'), //string
		highestShot = $centre.getAttribute('data-top'),   //string
		highestCentres = ~~(score.t[0] / highestShot),
		errorMessage;
	$centre.max = highestCentres;
	switch(true){
		//Using string attribute names here otherwise ClosureCompiler changes them.
	case score['t'][0] < 0 || score['t'][0] > totalMax:
		errorMessage = 'Please enter a total between 0 and ' + totalMax + '.';
		break;
	case score['c'][0] < 0 || score['c'][0] > highestPossibleCentres:
		errorMessage = 'Please enter centres between 0 and ' + highestPossibleCentres + '.';
		break;
	case score['c'][0] > highestCentres:
		errorMessage = 'Score has too many centres for a total of ' + score['t'][0] + '.<br>Please decrease centres to ' + highestCentres + '<br><b>OR</b><br>increase total to ' + score['c'][0] * highestShot + '.';
	}
	if(errorMessage){
		//Display error message
		popup.innerHTML = errorMessage;
		//TODO get element exact position taking into account scrolling.
		popup.style.top = $centre.getBoundingClientRect().top + $centre.clientHeight + 'px';
		popup.removeAttribute('hidden');
		return
	}
	//No validation errors with total or centres
	popup.setAttribute('hidden', '');
	return 1;   //return true value
}