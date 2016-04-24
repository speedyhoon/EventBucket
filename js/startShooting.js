(function(){
	'use strict';
	var currentCell, currentRow, classes = {
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
		return +thisCell.getAttribute('data-' + attribute);
	}

	function recalculateTotal(value){
		var type = getShootersClass(), index = classes[type].validShots.indexOf(value), newer = ~~classes[type].validScore[index], newerC = ~~classes[type].validCenta[index], total, centers;
		currentCell.textContent = value;
		if(getCurrentNth() >= getNoOfSighters()){
			total = getValue(currentRow, 'total') - getValue(currentCell, 'value') + newer;
			centers = getValue(currentRow, 'centers') - getValue(currentCell, 'center') + newerC;
			currentRow.lastChild.innerHTML = total + (centers ? '<sup>' + centers + '</sup>' : '');
			currentRow.setAttribute('data-total', total);
			currentRow.setAttribute('data-centers', centers);
		}
		currentCell.setAttribute('data-value', newer);
		currentCell.setAttribute('data-center', newerC);
	}

	function getShots(){
		var cells = currentRow.getElementsByTagName('td'), shootersClass = classes[getShootersClass()], send = '', sighters = getNoOfSighters(), value = '', i=-1, max = cells.length, index;
		while(++i<max){
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
		var table = document.querySelector('table'), eventID = table.getAttribute('data-eventID'), rangeID = table.getAttribute('data-rangeID');
		//shots = encodeURI(getShots()).replace(/#/gi, '%23');	//hashes are converted after encodeURI to stop % being converted twice
		j.open('POST', '/16?E=' + eventID + '&R=' + rangeID + '&S=' + id + '&s=' + getShots(), true);
		//j.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
		j.send();
		j.onreadystatechange = function stateChanger(){
			if(j.readyState == 4){
				//TODO status ok && html == same - GREEN    else    RED
				currentRow.querySelector('.t').innerHTML = j.status === 200 ? j.response : 'failed';
			}
		};
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
				ajax(currentRow.getAttribute('id'));
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
			var h=-1, td = document.createElement('td'), buttonLength = classes[type].buttons.length, buttonOnClickEvent = function buttonClickEventer(buttonValue){
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
			td.id='bu';
			td.setAttribute('colspan', buttonsCell.getAttribute('colspan'));
			buttonsCell.parentNode.replaceChild(td, buttonsCell);
		}
	}

	function moveHeader(){
		currentRow.parentNode.insertBefore(document.getElementById('h'), currentRow);
		currentRow.parentNode.insertBefore(document.getElementById('x'), currentRow.nextSibling);//equivilent to insertAfter!
		generateButtons();
	}

	/*	function changeSighters(){//when changing the value of the select box "sighters"
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
			if(!currentRow.getAttribute('data-visited')){
				currentRow.onclick = function tdClicker(tdElement){
					return function tdClick(event){
						if(event.target.nodeName === 'TD'){
							highlightCell(tdElement);
						}
					};
				};
				currentRow.setAttribute('data-visited', 1);
			}
		};
	};
	while(shooterQty--){		//assign onclick events to all shooters names
		shooters[shooterQty].onclick = shooterNameOnclick(shooters[shooterQty].parentNode);
	}

	//178,161,166,168,238,260,241
	//6328,5685,7920
}());