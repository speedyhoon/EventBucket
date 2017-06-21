'use strict';
var currentCell, currentRow, currentType, currentClass,
	eventID = window.location.pathname.split('/')[2],
	rangeID = window.location.pathname.split('/')[3],
	passed = true,
	classes;

function getCurrentNth(){
	return Array.prototype.indexOf.call(currentRow.getElementsByTagName('td'), currentCell);
}

function recalculateTotal(value){
	var shots = currentClass.marking.shots,
		newValue = shots[value].value,
		oldValue = ~~(shots[currentCell.textContent] && shots[currentCell.textContent].value),
		newCentre = ~~shots[value].center,
		oldCentre = ~~(shots[currentCell.textContent] && shots[currentCell.textContent].center),
		total = 0, centers = 0;
	currentCell.textContent = value;
	if(getCurrentNth() >= currentClass.sightersQty){
		total = ~~currentRow.getAttribute('data-total') - oldValue + newValue;
		centers = ~~currentRow.getAttribute('data-centers') - oldCentre + newCentre;
		currentRow.querySelector('th:last-of-type').innerHTML = total + (centers ? '<sup>' + centers + '</sup>' : '');
		currentRow.setAttribute('data-total', total);
		currentRow.setAttribute('data-centers', centers);
	}
}

function getShots(){
	var cells = currentRow.getElementsByTagName('td'),
		send = '',
		value,
		i = -1,
		max = cells.length;
	while(++i < max){
		value = cells[i].textContent;
		if(!value){
			send += '-';
		}else if(i < currentClass.sightersQty){
			send += currentClass.marking.shots[value].sighter || '-';
		}else{
			send += currentClass.marking.shots[value].shot;
		}
	}
	return send.replace(/-+$/, '');  //trim training hyphens
}

function highlightOnlyTheCell(cell){
	if(cell !== currentCell){
		if(currentCell){
			currentCell.removeAttribute('data-selected');
		}
		if(cell){
			cell.setAttribute('data-selected', 1);
		}
		currentCell = cell;
	}
}

function changeValue(value){
	if(currentCell){
		if(currentCell.textContent !== value){//prevents recalculating the score if it is the same value
			recalculateTotal(value);
			var id = currentRow.children[0].textContent;
			ws.send('\u0010' + JSON.stringify({E: [eventID], R: [rangeID], S: [id], s: [getShots()]}));
		}
		if(currentCell.nextSibling && currentCell.nextSibling.nodeName === 'TD'){
			highlightOnlyTheCell(currentCell.nextSibling);
		}
	}
}

function generateButtons(){
	if(currentClass && currentType !== currentClass.marking.buttons){
		currentType = currentClass.marking.buttons;
		var h = -1, th = document.createElement('th'), buttonLength = currentType.length, buttonOnClickEvent = function buttonClickEventer(buttonValue){
				return function buttonClicker(){
					changeValue(buttonValue);
				};
			};
		while(++h < buttonLength){
			var button = document.createElement('button');
			button.textContent = currentType[h];
			button.onclick = buttonOnClickEvent(currentType[h]);
			th.appendChild(button);
		}
		var buttonsCell = document.getElementById('bu');
		th.id = 'bu';
		th.setAttribute('colspan', buttonsCell.getAttribute('colspan'));
		buttonsCell.parentNode.replaceChild(th, buttonsCell);
	}
}

function moveHeader(){
	//If currentRow is not the first row in tbody
	if(Array.prototype.indexOf.call(currentRow.parentElement.children, currentRow)){
		//move header row before currentRow
		currentRow.parentNode.insertBefore(document.getElementById('h'), currentRow);
		//Make the header row visible
		document.getElementById('h').removeAttribute('hidden');
	}else{
		//Otherwise hide the header row if the currentRow is first in tbody
		document.getElementById('h').setAttribute('hidden','');
	}
	currentRow.parentNode.insertBefore(document.getElementById('x'), currentRow.nextSibling);//equivalent to insertAfter!
	generateButtons();
	document.getElementById('x').removeAttribute('hidden');
}
function highlightRow(row){//HIGHLIGHT THE SELECTED ROW
	if(row !== currentRow){
		row.setAttribute('data-selected', '1');
		if(currentRow){
			currentRow.removeAttribute('data-selected');
		}
		currentRow = row;
		setCurrentClass();
	}
}

function setCurrentClass(){
	passed = classes && classes[currentRow.getAttribute('data-class')];
	if(passed){
		currentClass = classes[currentRow.getAttribute('data-class')];
		highlightOnlyTheCell(currentRow.querySelector('td'));
		moveHeader();
		//modifySelectBox();
	}
}

function highlightCell(cell){//change the selected table cell (td) to the currentCell selected
	if(cell !== currentCell){
		highlightOnlyTheCell(cell);
		if(currentRow !== currentCell.parentNode){
			highlightRow(cell.parentNode);//used for clicking on a different shooters shot without clicking on a name first
		}
	}
}

function shooterNameOnclick(trElement){
	return function shooterClick(){
		highlightRow(trElement);

		//if visited attribute is present it has already been processed
		if(!currentRow.visited){
			currentRow.onclick = function trClicker(tdElement){
				return function trClick(event){
					//console.log('tdClicker');
					if(event.target.nodeName === 'TD'){
						highlightCell(tdElement);
					}
				};
			};
			//TODO click on all td elements in the row.
			//var trElement.getElementsByTagName('td')
			currentRow.visited = 1;
		}//else{console.log('already visited');}
	};
}

var shooters = document.querySelectorAll('tbody :not(#h) th:nth-child(4)'), shooterQty = shooters.length;
while(shooterQty--){		//assign onclick events to all shooters names
	shooters[shooterQty].onclick = shooterNameOnclick(shooters[shooterQty].parentNode);
}

var ws, intervalId;
function reconnect (){
	ws = new WebSocket('ws://'+window.location.host+'/k/');
	ws.onopen = function(){
		if(intervalId){
			clearInterval(intervalId);
			intervalId = 0;
		}else{
			ws.send('\u007E');   //126
		}
	};
	//Update UI with save / error message.
	ws.onmessage = function(message){
		var data = JSON.parse(message.data.substr(1));
		switch(message.data[0]){
		case'~':
			classes = data;
			if(!passed){
				setCurrentClass();
			}
			//console.log(classes);
			break;
		case'!':
			document.getElementById(data.S).parentElement.children[4].className = '';
			if(rangeID === data.R){
				document.getElementById(data.S).parentElement.children[4].innerHTML = data.T;
				//TODO status ok && html == same - GREEN	else	RED
				setTimeout(function() {
					document.getElementById(data.S).parentElement.children[4].className = '^save^';
				}, 10);
			}
		}
	};
	ws.onclose = function(){
		if(!intervalId){
			intervalId = setInterval(reconnect, 3000); //try to reconnect every 3 seconds after the connection is dropped.
		}
	};
}
reconnect();