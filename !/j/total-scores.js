(function(){
	'use strict';
	var index,
		inputs = document.getElementsByName('score'),
		summitOnchange = function(){
			return function(){
				console.log('down');
//				updateTotalScores
			};
		};
	for(index in inputs){
		if(inputs.hasOwnProperty(index)){
			inputs[index].onkeyup = summitOnchange();
		}
	}



	var j;
	if(window.XMLHttpRequest){
		j = new XMLHttpRequest();
	}
	function ajax(){
		j.open('POST', '/queryShooterList', true);// + 'scoreSave=' + classes.eventId + '~' + Id + '~' + shots, true);
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
}());