'use strict';
//Search for a shooter on Scorecards & total-scores pages.
function shooterBarcode(search, shooter){
	document.getElementById('barcodeErr').setAttribute('hidden', '');
	document.getElementById('shooterErr').setAttribute('hidden', '');
	var pathName = window.location.pathname.split('/')[1],
		eventID = window.location.pathname.split('/')[2],
		rangeID = window.location.pathname.split('/')[3];
	if(shooter && shooter.value && /^\d+$/g.test(shooter.value)){
		otherFunc(shooter.value, shooter, pathName, eventID, rangeID);
		return false;
	}else if(!search || !search.value || !/^\d+\/\d+#\d+$/g.test(search.value)){
		search.select();
		document.getElementById('barcodeErr').removeAttribute('hidden');
		return false;
	}
	var barcodeEventID = search.value.split('/')[0], barcodeRangeID = search.value.split('/')[1].split('#')[0], shooterID = search.value.split('#')[1];
	if(eventID !== barcodeEventID){
		//Go to a different event if user presses OK.
		if(confirm('event is not the same. do you want to go to event with id X?')){
			window.location.href = '/' + pathName + '/' + barcodeEventID + '/' + barcodeRangeID + '#' + shooterID;
		}
		//Else do nothing.
		search.select();
		return false;
	}
	if(rangeID !== barcodeRangeID){
		//Go to a different range if user presses OK.
		if(confirm('Jump to range with ID ' + barcodeRangeID + '?')){
			window.location.href = '/' + pathName + '/' + barcodeEventID + '/' + barcodeRangeID + '#' + shooterID;
		}
		//Else do nothing.
		search.select();
		return false;
	}
	return otherFunc(shooterID, search, pathName, eventID, rangeID);
}
if(!window.location.hash){
	document.querySelector('[name=B]').setAttribute('autofocus', '');
}else if(!document.getElementById(window.location.hash.replace('#', ''))){
	document.getElementById('shooterErr').removeAttribute('hidden');
	document.querySelector('[name=B]').select();
}
function otherFunc(shooterID, search, pathName, eventID, rangeID){
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
