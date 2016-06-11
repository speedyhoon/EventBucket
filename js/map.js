'use strict';

var map, places, j = new XMLHttpRequest(), path = document.location.pathname.split('/');
//Add clubID to GET parameters if it exists. Otherwise the response will return all clubs, i.e. path = ['', 'page-name', 'clubID']
j.open('GET', '/18'+(path[2] ? '?C='+path[2] :''), true);
j.send();
j.onreadystatechange = function stateChanger(){
	if(j.readyState == 4 && j.status === 200 && j.response){
		places = JSON.parse(j.response);
		var script = document.createElement('script');
		script.src = 'https://maps.googleapis.com/maps/api/js?key=AIzaSyCvMoHGEB9iW2i7VakZevwh3GhdXsL2eik&callback=initMap';
		document.body.appendChild(script);
	}
};

function initMap() {
	//When the user clicks the marker, an info window opens.
	var map, t, u = {zoom: 4}, infoW = function (gh) { return function() {
		gh.infowindow.open(map, gh.marker);
	}};
	if(places.length){
		u.center = {lat: places[0].x, lng: places[0].y};
	}
	map = new google.maps.Map(document.getElementById('map'), u);

	for(t in places){
		var url = places[t].u ? '<a href='+places[t].u+'>'+places[t].n+'</a>' : places[t].n;
		var town = places[t].t ? '<br>'+places[t].t : '';
		town += town && places[t].p ? ', '+places[t].p : '';
		places[t].infowindow = new google.maps.InfoWindow({
			content: '<h1>'+url+'</h1>'+(places[t].a||'')+town
		});
		places[t].marker = new google.maps.Marker({
			position: {lat: places[t].x, lng: places[t].y},
			map: map,
			title: places[t].n
		});
		places[t].marker.addListener('click', infoW(places[t]));
	}
}