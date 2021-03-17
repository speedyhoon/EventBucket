var button = document.querySelectorAll(
    '[action="/form.eventShooterNew.action"] button'
  ),
  i = button.length;
while (i--) {
  button[i].onclick = function shooterEntry(event) {
    var isNew = !event.target.getAttribute("formAction");
    event.target.form.schemaFirstName.required = isNew;
    event.target.form.schemaSurname.required = isNew;
    event.target.form.schemaShooter.required = !isNew;
    event.target.form.schemaClub.required = isNew;
  };
}
window.buildRow = function (tr, tds) {
  var surname = tds[1].querySelector("span");
  tr.querySelector("[name=schemaSurname]").value = surname.textContent;
  surname.textContent = "";
  tr.querySelector("[name=schemaFirstName]").value = tds[1].textContent.trim();
  tr.querySelector("[name=schemaClub]").value = tds[2].textContent;
  findValue(tr.querySelector("[name=schemaGrade]"), tds[3].textContent);
  findValue(tr.querySelector("[name=schemaAgeGroup]"), tds[4].textContent);
  tr.querySelector("[name=schemaShooter]").value = tds[0].textContent;
  checked(tds[5], tr.querySelector("[name=schemaLocked]"));
  checked(tds[6], tr.querySelector("[name=schemaSex]"));
  return tr;
};
