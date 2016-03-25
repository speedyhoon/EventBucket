'use strict';
function stuff(row) {
  row = row.parentElement;
  var tds = row.children, t = document.importNode(document.querySelector('template').content, true);
  t.querySelector('[name=s]').value = row.children[1].querySelector('span').textContent;
  tds[1].querySelector('span').textContent = '';
  t.querySelector('[name=f]').value = tds[1].textContent.trim();
  t.querySelector('[name=C]').value = tds[2].textContent;
  t.querySelector('[name=g]').value = findValue(t.querySelector('[name=g]'), tds[3].textContent);
  t.querySelector('[name=r]').value = findValue(t.querySelector('[name=r]'), tds[4].textContent);

  t.querySelector('[name=I]').value = tds[0].textContent;
  t.querySelector('td').textContent = tds[0].textContent;

  //  console.log(t);
  row.innerHTML = '';
  row.appendChild(t);
  //  console.log(t);
}
function findValue(element, label) {
  var options = element.options, i = options.length;
  while (i--) {
    if (label === options[i].textContent) {
      return options[i].value;
    }
  }
}
