(function(){
	'use strict';
	var currentCell = null, currentRow = null, j, classes = {
		eventId:1,
		currentType:null,
		0:{
			sighters:2,
			validShots:'012345V6X',
			validScore:'012345555',
			validCenta:'000000111',
			validSighters:')!@#$%v^x',
			buttons:'012345VX'
		},
		1:{
			sighters:2,
			validShots:'012345V6X',
			validScore:'012345666',
			validCenta:'000000001',
			validSighters:')!@#$%v^x',
			buttons:'0123456X'
		},
		2:{
			sighters:2,
			validShots:'012345V6X',
			validScore:'012345555',
			validCenta:'000000111',
			validSighters:')!@#$%v^x',
			buttons:'012345VX'
		}
	};
	if(window.XMLHttpRequest){
		j = new XMLHttpRequest();
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
		return ~~thisCell.getAttribute('data-' + attribute);
	}

	function recalculateTotal(value){
		var type = getShootersClass(),
			index = classes[type].validShots.indexOf(value),
			newer = ~~classes[type].validScore[index],
			newerC = ~~classes[type].validCenta[index],
			total,
			centres;
		currentCell.innerHTML = value;
		if(getCurrentNth() >= getNoOfSighters()){
			total = getValue(currentRow, 'total') - getValue(currentCell, 'value') + newer;
			centres = getValue(currentRow, 'centres') - getValue(currentCell, 'center') + newerC;
			currentRow.lastChild.innerHTML = total + '<sup>' + centres + '</sup>';
			currentRow.setAttribute('data-total', total);
			currentRow.setAttribute('data-centres', centres);
		}
		currentCell.setAttribute('data-value', newer);
		currentCell.setAttribute('data-center', newerC);
	}

	function getAjax(){
		var cells = currentRow.getElementsByTagName('td'), shootersClass = classes[getShootersClass()], send = '', sighters = getNoOfSighters(), value = '', i, max = cells.length, index;
		for(i = 0; i < max; i++){
			value = cells[i].innerHTML;
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
		return send;
	}

	function ajax(id){
		var shots = getAjax();
		shots = encodeURI(shots).replace(/#/gi, '%23');	//hashes are converted after encodeURI to stop % being converted twice
		j.open('POST', '/updateShotScores?eventid='+eventId+'&rangeid='+rangeId+'&shooterid='+id+'&shots='+shots, true);// + 'scoreSave=' + classes.eventId + '~' + Id + '~' + shots, true);
		j.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
		j.send();
		j.onreadystatechange = function(){
			if(j.status === 200){
				document.getElementById(id).getElementsByClassName('t')[0].innerHTML = j.response;//total+'<sup>'+centres+'</sup>';
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
			if(currentCell.innerHTML !== value){//prevents recalculating the score if it is the same value
				recalculateTotal(value);
				ajax(currentRow.getAttribute('id'));
			}
			if(currentCell.nextSibling && currentCell.nextSibling.nodeName === 'TD'){
				highlightOnlyTheCell(currentCell.nextSibling);
			}
		}
	}

	function buttonsAddEvents(){
		var n = 0,
			buttons = document.getElementsByTagName('button'),
			max = buttons.length, buttonOnClickEvent = function(buttonValue){
				return function(){
					changeValue(buttonValue);
				};
			};
		for(n; n < max; n++){
			buttons[n].onclick = buttonOnClickEvent(buttons[n].innerHTML);
		}
	}

	function generateButtons(){
		var type = getShootersClass(), buttonLength, buttonHtml, n;
		if(classes.currentType !== type){
			classes.currentType = type;
			if(type && classes[type].buttons){
				buttonLength = classes[type].buttons.length;
				buttonHtml = '';
				for(n = 0; n < buttonLength; n++){
					buttonHtml += '<button>' + classes[type].buttons[n] + '</button>';
				}
				document.getElementById('bu').innerHTML = buttonHtml;
				buttonsAddEvents();
			}
		}
	}

	function moveHeader(){
		var tbody = currentRow.parentNode;
		tbody.insertBefore(document.getElementById('h'), currentRow);
		tbody.insertBefore(document.getElementById('x'), currentRow.nextSibling);//equivilent to insertAfter!
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
				recalculateTotal(tds[i].innerHTML);
			}
			currentCell = selectedCell;
			highlightOnlyTheCell(tds[(sighters + selected)]);
			ajax(currentRow.getAttribute('id'));
		}
	}*/

/*	function modifySelectBox(){//alter the options in the sighters select box for different classes
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
			highlightOnlyTheCell(row.getElementsByTagName('td')[0]);
			moveHeader();
//			modifySelectBox();
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

	function listenToTds(){
		var n = 0, tdElements = currentRow.getElementsByTagName('td'), tdQuantity = tdElements.length, tdOnClickEvent = function(tdElement){
				return function(){
					highlightCell(tdElement);
				};
			};
		for(n; n < tdQuantity; n++){
			tdElements[n].onclick = tdOnClickEvent(tdElements[n]);
		}
		currentRow.setAttribute('data-visited', 1);
	}

	function runOnLoad(){
		var n = 0, shooters = document.getElementsByClassName('name'), shooterQuantity = shooters.length, shooterNameOnclick = function(trElement){
				return function(){
					highlightRow(trElement);
					if(!currentRow.getAttribute('data-visited')){  //if the attribute is present it has already been processed
						listenToTds();
					}
				};
			};
		if(shooters){//assign onclick events to all shooters names
			for(n; n < shooterQuantity; n++){
				shooters[n].onclick = shooterNameOnclick(shooters[n].parentNode);
			}
		}
	}

	runOnLoad();//178,161,166,168,238,260,241
	//6328,5685,7920
}());