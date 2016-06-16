'use strict';
var ws, intervalId, eventID = window.location.pathname.split('/')[2];
function reconnect (){
	ws = new WebSocket('ws://'+window.location.host+'/k/');
	ws.onopen = function(){
		if(intervalId){
			clearInterval(intervalId);
			intervalId = 0;
		}
		console.log("open");
		ws.send('}');
	};
	//TODO refresh page when receiving any message
	//Update UI with save / error message.
	ws.onmessage = function(message){
		console.log(message.data);
	};
	ws.onclose = function(){
		if(!intervalId){
			intervalId = setInterval(reconnect, 15000); //try to reconnect every 15 seconds after the connection is dropped.
		}
	};
}
reconnect();