'use strict';
var xhttp = new XMLHttpRequest;
xhttp.onreadystatechange = function() {
  if (xhttp.readyState == 4 && xhttp.status == 200) {
    console.log(xhttp.responseText);
  }
};
xhttp.open('GET', '/9?f=Fred&s=Smith&C=Treeville', true);
xhttp.send();
