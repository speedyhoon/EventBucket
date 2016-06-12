(function(){
	'use strict';
	var currentCell, currentRow,
		eventID = window.location.pathname.split('/')[2],
		rangeID = window.location.pathname.split('/')[3],
		classes = {
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
		};
	if(window.XMLHttpRequest){
		var j = new XMLHttpRequest();
	}

	function getShootersClass(){
		return currentRow.getAttribute('data-class');
	}

	function getCurrentNth(){
		return Array.prototype.indexOf.call(currentRow.getElementsByTagName('td'), currentCell);
	}

	function getNoOfSighters(){
		return classes[getShootersClass()].sighters;
	}

	function getValue(thisCell, attribute){
		//Putting plus in front of a string converts it to a number (integer or float)
		return+thisCell.getAttribute('data-' + attribute);
	}

	function recalculateTotal(value){
		var type = getShootersClass(), index = classes[type].validShots.indexOf(value), newer = ~~classes[type].validScore[index], newerC = ~~classes[type].validCenta[index], total, centers;
		currentCell.textContent = value;
		//console.log(currentRow.total, currentRow.centers, currentCell.value, currentCell.center);
		if(getCurrentNth() >= getNoOfSighters()){
			total = getValue(currentRow, 'total') - getValue(currentCell, 'value') + newer;
			centers = getValue(currentRow, 'centers') - getValue(currentCell, 'center') + newerC;
			currentRow.lastChild.innerHTML = total + (centers ? '<sup>' + centers + '</sup>' : '');
//currentRow.setAttribute('data-total', total);
			currentRow.total = total;
//currentRow.setAttribute('data-centers', centers);
			currentRow.centers = centers;
		}
//currentCell.setAttribute('data-value', newer);
		currentCell.value = newer;
		currentCell.center = newerC;
	}

	function getShots(){
		var cells = currentRow.getElementsByTagName('td'), shootersClass = classes[getShootersClass()], send = '', sighters = getNoOfSighters(), value = '', i = -1, max = cells.length, index;
		while(++i < max){
			value = cells[i].textContent;
			if(!value){
				send += '-';
			}else if(i < sighters){
				index = shootersClass.validShots.indexOf(value);
				if(index && shootersClass.validSighters[index]){
					//send+=encodeURIComponent(shootersClass.validSighters[index]);
					send += shootersClass.validSighters[index];
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
		var type = getShootersClass();
		if(type && classes.currentType !== type && classes[type].buttons){
			classes.currentType = type;
			var h = -1, td = document.createElement('td'), buttonLength = classes[type].buttons.length, buttonOnClickEvent = function buttonClickEventer(buttonValue){
					return function buttonClicker(){
						changeValue(buttonValue);
					};
				};
			while(++h < buttonLength){
				var button = document.createElement('button');
				button.textContent = classes[type].buttons[h];
				button.onclick = buttonOnClickEvent(classes[type].buttons[h]);
				td.appendChild(button);
			}
			var buttonsCell = document.getElementById('bu');
			td.id = 'bu';
			td.setAttribute('colspan', buttonsCell.getAttribute('colspan'));
			buttonsCell.parentNode.replaceChild(td, buttonsCell);
		}
	}

	function moveHeader(){
		currentRow.parentNode.insertBefore(document.getElementById('h'), currentRow);
		currentRow.parentNode.insertBefore(document.getElementById('x'), currentRow.nextSibling);//equivilent to insertAfter!
		generateButtons();
	}

	/*function changeSighters(){//when changing the value of the select box "sighters"
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
	}*/
	function highlightRow(row){//HIGHLIGHT THE SELECTED ROW
		if(row !== currentRow){
			row.setAttribute('data-selected', '1');
			if(currentRow){
				currentRow.removeAttribute('data-selected');
			}
			currentRow = row;
			highlightOnlyTheCell(row.querySelector('td'));
			moveHeader();
			//modifySelectBox();
		}
	}

	function highlightCell(cell){//change the selected table cell (td) to the currentCell selected
		//console.log('highlightCell');
		if(cell !== currentCell){
			highlightOnlyTheCell(cell);
			if(currentRow !== currentCell.parentNode){
				highlightRow(cell.parentNode);//used for clicking on a different shooters shot without clicking on a name first
			}
		}
	}

	var shooters = document.querySelectorAll('tbody th:nth-child(4)'), shooterQty = shooters.length, shooterNameOnclick = function(trElement){
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
				}else{console.log('row already visited');}
			};
		};
	while(shooterQty--){		//assign onclick events to all shooters names
		shooters[shooterQty].onclick = shooterNameOnclick(shooters[shooterQty].parentNode);
	}

	var ws, intervalId;
	function reconnect (){
		ws = new WebSocket('ws://'+window.location.host+'/k/');
		ws.onopen = function(){
			if(intervalId){
				clearInterval(intervalId);
				intervalId = undefined;
			}
		};
		//TODO
		//Update UI with save / error message.
		ws.onmessage = function(message){
			console.log(message)
		};
		ws.onclose = function(){
			if(!intervalId){
				intervalId = setInterval(reconnect, 3000); //try to reconnect every 3 seconds after the connection is dropped.
			}
		};
	}
	reconnect();
}());