var count = 1;

var searchfield = `      <input type="text" id="Column${count}" name="Column${count}">
      <select id="Operator${count}" name="Operator${count}">
        <option name="Operator${count}" value="=" selected="selected">=</option>
        <option name="Operator${count}" value=">=">>=</option>
        <option name="Operator${count}" value="<="><=</option>
        <option name="Operator${count}" value=">">></option>
        <option name="Operator${count}" value="<"><</option>
      </select>
      <input type="text" id="Value${count}" name="Value${count}">`;

function addField(divname) {
	// Adds new column, operator, and value fields
	var newdiv = document.createElement("div");
	newdiv.innerHTML = searchfield;
	document.getElementById(divname).appendChild(newdiv);
	count++;
}
