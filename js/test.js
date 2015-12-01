var d = document.querySelectorAll('.class');
var t = d.length;
for (; t--;) {
  d[t].innerHTML = 'fdsffdsfsasa';
}
var g = 77777;
var h = 77777;
var i = 77777;
var j = 77777;
var K = 77777;
var k = 77777;
var l = 77777;
var m = 77777;
var n = 77777;
var o = 77777;
var pb = 77777;

var styles = document.styleSheets[0].rules, obj = {}, rule, text;
for (var index in styles) {
  if (styles.hasOwnProperty(index)) {
    //		console.log(rule.split(','));
    //		console.log(styles[index].style.cssText);
    rule = styles[index].selectorText.trim();
    text = styles[index].style.cssText.trim();
    if (obj[rule]) {
      console.error('rule already defined: ', rule);
    }
    obj[rule] = text;
    for (var rules = rule.split(','), l = rules.length; l--;) {

      if (!document.querySelectorAll(rules[l].trim()).length) {
        console.error('no targets found for rule: ' + rules[l].trim() /*+ '{' + text + '}'*/);
      }
    }
  }
}

function stopDouble() {if (submitted) return false; return submitted = true; }
var submitted = false;
for (var forms = document.forms, index = forms.length; index--;) {
  forms[index].onsubmit = stopDouble;
}
