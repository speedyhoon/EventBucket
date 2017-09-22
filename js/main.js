//Prevent users from double/triple clicking on submit buttons that could cause a form to submit more than once with exactly the same data. If forms used a CSRF token this JavaScript would be redundant
var submitted = 0, i = document.forms.length;
function stopDoubleSubmit(){
	if(submitted)return 0;
	return submitted = 1;
}

while(i--){
	if(!document.forms[i].onsubmit){
		document.forms[i].onsubmit = stopDoubleSubmit;
	}
}

var $dialog
	,$label;

document.onclick = function(event){
	switch(event.target.nodeName){
	case 'ABBR':
		//Help message pop-up
		if(!$dialog){
			//Form help dialog popup
			$dialog = document.createElement('dialog');
			$dialog.onclick = function(){
				$dialog.open = 0;
				//If label is set, click on the parent label element
				if($label){
					$label.click();
					$label = 0;
				}
			};
			document.body.appendChild($dialog);
		}
		$dialog.textContent = event.target.title;
		//Assuming the <abbr> element will always be the immediate child of a <label> element
		$label = event.target.parentElement;
		$dialog.open = 1;
		break;
	case 'TD':
		tableSort(event.target);
	}
};

//Start event listeners that add a class to form fields based on their valid/invalid values that changes the background colour.
var inputs = document.querySelectorAll('input,select')
	,flagClass = function(event){
		event.srcElement.classList.toggle('^dirty^', !event.srcElement.validity.valid);
	};
i = inputs.length;
while(i--){
	inputs[i].addEventListener('blur', flagClass);
	inputs[i].addEventListener('invalid', flagClass);
	inputs[i].addEventListener('input', flagClass);
}

//Add table sort to <td> elements within a <thead>
function tableSort($th){
	if($th.parentElement.parentElement.nodeName !== 'THEAD'){
		return;
	}
	//If th.textContent == id compare using base36 else use textContent
	var $tbody = $th.parentElement.parentElement.parentElement.querySelector('tbody')
		,column = Array.prototype.indexOf.call($th.parentElement.children, $th)
		,direction = $th.className === '^asc^' ? 1 : -1
		,$rows = Array.from($tbody.children);
	var sortBy = function($cell){
		if(!$cell){
			return '';
		}
		switch($th.textContent){
		case 'ID': //Numeric integer identifier
			return ~~$cell.textContent;
		case 'Id': //Base 36 [0-9a-z] identifier string e.g. a2e = 13046. Used on the Shooters page
			return parseInt($cell.textContent, 36);
		}
		return $cell.textContent;
	};
	$rows.sort(function(a, b){
		a = sortBy(a.children[column]);
		b = sortBy(b.children[column]);
		return a > b ? direction : a < b ? -1 * direction : 0;
	});
	$tbody.innerHTML = '';
	for(var i = 0, max = $rows.length; i < max; i++){
		$tbody.appendChild($rows[i]);
	}
	var $ths = $th.parentElement.querySelectorAll('.^asc^,.^desc^'), qty = $ths.length;
	while(qty--){
		$ths[qty].removeAttribute('class');
	}
	$th.className = direction < 0 ? '^asc^' : '^desc^';
}

var headings = document.querySelectorAll('thead');
i = headings.length;
while(i--){
	if(!headings[i].querySelector('.^asc^,.^desc^')){
		headings[i].querySelector('td').classList.add('^asc^');
	}