'use strict';
document.onclick = function(event) {
  if (event.target.nodeName === 'TH') {
    tableSort(event);
  }
};

//function tableSort(th, sortAs){
function tableSort(th) {
  th = th.target;
  /*If th.textContent == id compare using base36 else use textContent */
  var tbody = th.parentElement.parentElement.parentElement.querySelector('tbody'),
	  column = Array.prototype.indexOf.call(th.parentElement.children, th),
	  direction = th.className === '^asc^' ? -1 : 1,
	  rows = Array.from(tbody.children);

  //Array.prototype.sort.call(tbody.children, function (a, b) {
  rows.sort(function(a, b) {
    a = a.children[column].textContent;
    b = b.children[column].textContent;
    return a > b ? direction : a < b ? -1 * direction : 0;
  });

  tbody.innerHTML = '';
  for (var i = 0, max = rows.length; i < max; i++) {
    tbody.appendChild(rows[i]);
  }
  var ths = th.parentElement.querySelectorAll('.^asc^,.^desc^'), qty = ths.length;
  while (qty--) {
    ths[qty].removeAttribute('class');
  }
  th.className = direction > 0 ? '^asc^' : '^desc^';
}
