if (document.styleSheets[0]) {
  var styles = document.styleSheets[0].rules, obj = {}, rule, text, index;
  for (index in styles) {
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
}

function stopDoubleSubmit() {if (submitted) return false; return submitted = true; }
var submitted = false;
for (var forms = document.forms, index = forms.length; index--;) {
  forms[index].onsubmit = stopDoubleSubmit;
}
