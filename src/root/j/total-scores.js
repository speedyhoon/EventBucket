(function(){
"use strict";
	var j, index, inputs = document.getElementsByName("score");
	if(window.XMLHttpRequest){
		j = new XMLHttpRequest();
	}
	function down(form){
		var i, formLength = form.length, output = form[0].name + "=" + form[0].value;
		for(i=1; i < formLength; i++){
			output += "&"+ form[i].name + "=" + form[i].value;
		}

		j.open('POST', form.action, true);
		j.setRequestHeader('Content-type', form.encoding);
		j.send(output);
		j.onreadystatechange = function(){
			if(j.status === 200){
				console.log("saved!", j.response);
			}
		};
	}

	var summit_onchange = function(){
		return function(){
			down(this.form);
		};
	};
	for(index in inputs){
		if(inputs.hasOwnProperty(index)){
			inputs[index].onkeyup = summit_onchange();
		}
	}
}());