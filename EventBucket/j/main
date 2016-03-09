'use strict';

//Prevent users from double/triple clicking on submit buttons that could cause a form to submit more than once with exactly the same data. If forms used a CSRF token this JavaScript would be redundant
function stopDoubleSubmit() {if (submitted) return false; return submitted = true; }
var submitted = false;
for (var forms = document.forms, index = forms.length; index--;) {
  forms[index].onsubmit = stopDoubleSubmit;
}
