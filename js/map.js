(function () {
  var places,
    j = new XMLHttpRequest(),
    path = document.location.pathname.split("/");
  //Add clubID to GET parameters if it exists. Otherwise the response will return all clubs, i.e. path = ['', 'page-name', 'clubID']
  j.open(
    "GET",
    "/forms.mapClubs.action" + (path[2] ? "?C=" + path[2] : ""),
    true
  );
  j.send();
  j.onreadystatechange = function stateChanger() {
    if (j.readyState === 4 && j.status === 200 && j.response) {
      places = JSON.parse(j.response);
      var script = document.createElement("script");
      script.src =
        "https://maps.googleapis.com/maps/api/js?key=AIzaSyCvMoHGEB9iW2i7VakZevwh3GhdXsL2eik&callback=initMap";
      document.body.appendChild(script);
    }
  };

  //When the user clicks the marker, an info window opens.
  window.initMap = function () {
    var map,
      t,
      u = {
        zoom: 4,
      },
      infoW = function (gh) {
        return function () {
          gh.infowindow.open(map, gh.marker);
        };
      };
    if (places.length) {
      u.center = {
        lat: places[0].x,
        lng: places[0].y,
      };
    }
    map = new google.maps.Map(document.getElementById("map"), u);

    t = places.length;
    while (t--) {
      var place = places[t];
      var url = place.u ? `<a href=${place.u}>${place.n}</a>` : place.n;
      var town = place.t ? `<br>${place.t}` : "";
      town += town && place.p ? `, ${place.p}` : "";
      place.infowindow = new google.maps.InfoWindow({
        content: `<h1>${url}</h1>` + (place.a || "") + town,
      });
      place.marker = new google.maps.Marker({
        position: {
          lat: place.x,
          lng: place.y,
        },
        map: map,
        title: place.n,
      });
      place.marker.addListener("click", infoW(place));
    }
  };
})();
