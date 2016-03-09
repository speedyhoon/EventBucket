'use strict';
var f = '', s = '', c = '';
function getter(form) {
  //Ignore form inputs into other fields
  if (f === form.f.value && s === form.s.value && c === form.C.value) return;
  f = form.f.value;
  s = form.s.value;
  c = form.C.value;
  var xhttp = new XMLHttpRequest;
  xhttp.onreadystatechange = function() {
    if (xhttp.readyState == 4) {
      if (xhttp.status == 200) {
        form.S.innerHTML = xhttp.responseText;
        form.S.removeAttribute('class');
      }else {
        form.S.setAttribute('class', 'not200');
      }
    }
  };
  form.S.setAttribute('class', 'loading');
  xhttp.open('GET', '/9?f=' + f + '&s=' + s + '&C=' + c, true);
  xhttp.send();
}
