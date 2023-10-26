//Set shooter barcode form onsubmit because it's not allowed to be set directly in HTML with the current Content Security Policy.
document.querySelector('#sb').onsubmit = function(event){
	shooterBarcode(event.target.B);
	return false;
};

//If shooter ID is provided, try to find the shooter. Otherwise, autofocus on barcode search text box.
if(!window.location.hash){
	document.querySelector('[name=B]').setAttribute('autofocus', '');
}else{
	goToShooter(window.location.hash.replace('#', ''), document.querySelector('[name=B]'));
}

//Search for a shooter on Scorecards & total-scores pages.
function shooterBarcode($search){
	document.getElementById('searchErr').setAttribute('hidden', '');
	document.getElementById('barcodeErr').setAttribute('hidden', '');
	document.getElementById('shooterErr').setAttribute('hidden', '');
	if(!$search || !$search.value){
		document.getElementById('searchErr').removeAttribute('hidden');
		return;
	}
	if(/^\d+$/g.test($search.value)){
		goToShooter($search.value, $search);
		return;
	}
	//If barcode doesn't match display error message
	if(!/^\d+\/\d+#\d+$/g.test($search.value)){
		document.getElementById('barcodeErr').removeAttribute('hidden');
		$search.select();
		return;
	}
	var barcodeEventID = $search.value.split('/')[0]
		,barcodeRangeID = $search.value.split('/')[1].split('#')[0]
		,shooterID = $search.value.split('#')[1]
		,pathName = window.location.pathname.split('/')[1];

	if(eventID !== barcodeEventID){
		//Go to a different event if user presses OK.
		if(confirm(`This barcode is for a different event. Do you want to go to event with id ${barcodeEventID}?`)){
			window.location.href = `/${pathName}/${barcodeEventID}/${barcodeRangeID}#${shooterID}`;
		}
		//Else do nothing.
		$search.select();
		return;
	}
	if(rangeID !== barcodeRangeID){
		//Go to a different range if user presses OK.
		if(confirm(`This barcode is for a different range. Do you want to go to range with id ${barcodeRangeID}?`)){
			window.location.href = `/${pathName}/${barcodeEventID}/${barcodeRangeID}#${shooterID}`;
		}
		//Else do nothing.
		$search.select();
		return;
	}
	goToShooter(shooterID, $search);
}

function goToShooter(shooterID, search){
	//If the shooter text box exists in the DOM, set focus to their text box.
	var d = document.getElementById(shooterID)
		,pathName = window.location.pathname.split('/')[1];
	if(d){
		search.value = '';
		if(d.tagName === 'input'){
			d.select();
		}else{
			d.click();
		}
		window.location.hash = `#${shooterID}`;
	}else if(pathName.indexOf('-all') >= 0){
		//Display shooter not is this event error message.
		document.getElementById('shooterErr').removeAttribute('hidden');
	}else{
		//If the shooter doesn't exist go to the scorecards-all OR total-scores-all page
		window.location.href = `/${pathName}-all/${eventID}/${rangeID}#${shooterID}`;
	}
}