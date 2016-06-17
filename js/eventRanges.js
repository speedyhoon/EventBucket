'use strict';
function buildRow(t, tds, formID, row){
	var R = t.querySelector('[name=R]');
	if(tds[2].textContent){
		findValues(R, tds[2].textContent.replace(/, $/, '').split(', '));
		t.querySelector('form').setAttribute('action','/21')
	}else{
		R.parentElement.removeChild(R);
	}
	t.querySelector('[name=n]').value = tds[1].textContent;
	var l = t.querySelector('[name=k]');
	if(!tds[2].textContent){
		if(tds[3].className === '^tick^'){
			l.setAttribute('checked', '');
		}
	}else{
		l.parentElement.removeChild(l);
	}
	t.querySelector('[name=I]').value = tds[0].textContent;
	t.querySelector('[name=o]').value = Array.prototype.indexOf.call(row.parentNode.children, row);


	var aTags = t.querySelectorAll('span'), i = aTags.length;
	while(i--){
		aTags[i].onclick = moveRange;
	}
	return t;
}
function moveRange(event){
	var currentRow = event.target.parentNode.parentNode,
		tbody = currentRow.parentNode,
		qty = tbody.children.length- 1,
		index = Array.prototype.indexOf.call(tbody.children, currentRow);

	console.log(index , qty);
	//TODO change moving rows to use numbers instead. that might use less code.
	if(event.target.classList.contains('^asc^')){
		//Move Down
		//if row is not the last in the table;
		if(index !== qty){
			//Insert row after next sibling;
			currentRow.parentNode.insertBefore(currentRow.nextSibling, currentRow);
		}else{
			//Otherwise insert currentRow before the row first.
			currentRow.parentNode.insertBefore(currentRow, tbody.children[0]);
		}
	}else{
		//Move Up
		if(index){
			currentRow.parentNode.insertBefore(currentRow, currentRow.previousSibling);
		}else{
			currentRow.parentNode.insertBefore(currentRow, tbody.children[qty]);
			currentRow.parentNode.insertBefore(tbody.children[qty], currentRow);
		}
	}
	currentRow.querySelector('[name=o]').value = Array.prototype.indexOf.call(tbody.children, currentRow);
}