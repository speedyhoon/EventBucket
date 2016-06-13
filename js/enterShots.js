(function(){
	'use strict';
	var currentCell, currentRow, currentType, currentClass,
		eventID = window.location.pathname.split('/')[2],
		rangeID = window.location.pathname.split('/')[3],
		passed = true,
		classes;
	if(window.XMLHttpRequest){
		var j = new XMLHttpRequest();
	}

	function getCurrentNth(){
		return Array.prototype.indexOf.call(currentRow.getElementsByTagName('td'), currentCell);
	}

	function getNoOfSighters(){
		return currentClass.sighters;
	}

	function getValue(thisCell, attribute){
		//Putting plus in front of a string converts it to a number (integer or float)
		return+thisCell.getAttribute('data-' + attribute);
	}

	function recalculateTotal(value){
		var newValue = currentClass.marking.Shots[value].value,
			newCentre = ~~currentClass.marking.Shots[value].center,
			total = 0, centers = 0;
		currentCell.textContent = value;
		//console.log(currentRow.total, currentRow.centers, currentCell.value, currentCell.center);
		if(getCurrentNth() >= getNoOfSighters()){
			total = getValue(currentRow, 'total') - getValue(currentCell, 'value') + newValue;
			//total = ~~currentRow.total - ~~currentCell.value + newValue;
			centers = getValue(currentRow, 'centers') - getValue(currentCell, 'center') + newCentre;
			//centers = ~~currentRow.centers - ~~currentCell.center + newCentre;
			currentRow.lastChild.innerHTML = total + (centers ? '<sup>' + centers + '</sup>' : '');
			currentRow.setAttribute('data-total', total);
			//currentRow.total = total;
			currentRow.setAttribute('data-centers', centers);
			//currentRow.centers = centers;
		}


		currentCell.value = newValue;
		currentCell.center = newCentre;
		console.log('bugger', currentCell.value, currentCell.center);
	}

	function getShots(){
		var cells = currentRow.getElementsByTagName('td'), send = '', sighters = getNoOfSighters(), value = '', i = -1, max = cells.length, index;
		while(++i < max){
			value = cells[i].textContent;
			if(!value){
				send += '-';
			}else if(i < sighters){
				index = currentClass.validShots.indexOf(value);
				if(index && currentClass.validSighters[index]){
					//send+=encodeURIComponent(currentClass.validSighters[index]);
					send += currentClass.validSighters[index];
				}else{
					send += '-';
				}
			}else{
				//send+=encodeURIComponent(value);
				send += value;
			}
		}
		return send.replace(/-+$/, '');  //trim training hyphens
	}

	function ajax(id){
		/*TODO reimplement for websockets
		j.onreadystatechange = function stateChanger(){
			if(j.readyState == 4){
				//TODO status ok && html == same - GREEN    else    RED
				currentRow.querySelector('.t').innerHTML = j.status === 200 ? j.response : 'failed';
			}
		};*/
		//Form 16 - eventUpdateShotScore
		console.log({E: eventID, R: rangeID, S: id, s: getShots()});
		ws.send('\u0010' + JSON.stringify({E: [eventID], R: [rangeID], S: [id], s: [getShots()]}));
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
				ajax(currentRow.children[0].textContent);
			}
			if(currentCell.nextSibling && currentCell.nextSibling.nodeName === 'TD'){
				highlightOnlyTheCell(currentCell.nextSibling);
			}
		}
	}

	function generateButtons(){
		if(currentClass && currentType !== currentClass.marking.Buttons){
			currentType = currentClass.marking.Buttons;
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
			//console.log('visited?', currentRow.visited);
			if(!currentRow.visited){
				currentRow.onclick = function tdClicker(tdElement){
					return function tdClick(event){
						//console.log('tdClicker');
						if(event.target.nodeName === 'TD'){
							highlightCell(tdElement);
						}
					};
				};
				currentRow.visited = 1;
			}
		};
	};

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
		//TODO
		//Update UI with save / error message.
		ws.onmessage = function(message){
			switch(message.data[0]){
			case'~':
				classes = JSON.parse(message.data.substr(1));
				if(!passed){
					setCurrentClass();
				}
				console.log(classes);
				break;
			case'!':
				classes = JSON.parse(message.data.substr(1));
				if(!passed){
					setCurrentClass();
				}
				console.log(classes);
				break;
			default:
				console.log('no handler for:', message.data);
			}
		};
		ws.onclose = function(){
			if(!intervalId){
				intervalId = setInterval(reconnect, 3000); //try to reconnect every 3 seconds after the connection is dropped.
			}
		};
	}
	reconnect();
}());




/*
classes= {
			currentType:null,
			0:{
				sighters:2,
				validShots:'012345V6X',
				validScore:'012345555',
				validCenta:'000000111',
				validSighters:'abcdefvvx',
				buttons:'012345VX'
			},
			1:{
				sighters:3,
				validShots:'012345V6X',
				validScore:'012345666',
				validCenta:'000000001',
				validSighters:'abcdefggx',
				buttons:'0123456X'
			},
			2:{
				sighters:2,
				validShots:'012345V6X',
				validScore:'012345555',
				validCenta:'000000111',
				validSighters:'abcdefvvx',
				buttons:'012345VX'
			}
		}

function changeSighters(){//when changing the value of the select box "sighters"
	var selected = document.getElementById('selectSighters').selectedIndex, tds, iteration = 0, sighters, selectedCell, i;
	if(currentRow && selected > 0){
		tds = currentRow.getElementsByTagName('td');
		sighters = getNoOfSighters();
		selectedCell = currentCell;
		for(i = sighters - selected; i < sighters; i++){
			currentCell = tds[sighters + iteration++];//increment iteration AFTER this line
			recalculateTotal(tds[i].textContent);
		}
		currentCell = selectedCell;
		highlightOnlyTheCell(tds[(sighters + selected)]);
		ajax(currentRow.getAttribute('id'));
	}
}
function modifySelectBox(){//alter the options in the sighters select box for different classes
	var additional = '', k, selectBox = document.createElement('select'), selectSighters;
	for(k = getNoOfSighters(); k >= 2; k--){
		additional += '<option>Keep S' + k + ' &gt;</option>';
	}
	selectBox.id = 'selectSighters';
	selectBox.innerHTML = '<option>Drop All</option>' + additional + '<option>Keep All</option>';
	if(currentRow.getAttribute('data-sighters')){
		selectBox.selectedIndex = currentRow.getAttribute('data-sighters');
	}
	selectBox.onchange = function(){
		return function(){
			changeSighters();
		};
	};
	selectSighters = document.getElementById('selectSighters');
	selectSighters.parentNode.replaceChild(selectBox, selectSighters);
}

//	function getShootersClass(){
//		return currentRow.getAttribute('data-class');
//	}
*/