var count = 0;

const searchfield = {
	tables: ["Patient", "Diagnosis", "Tumor", "Taxonomy", "Common", "Totals", "Denominators", "Life_history", "Source", "Accounts", "Unmatched"],
	

};

/*function getOperation(n, type) {
	// Returns formatted operation block for given data type
	let ops = ["=", "!=", ">=", "<=", ">", "<"];
	let end = 2;
	let ret = [`      <select id="Operator${n}" name="Operator${n}">`];
	if (type === "DOUBLE") {
		// Add valid comparisons for floating point numbers
		end = ops.length;
	}
	for (let i = 0; i < end; i++) {
		// Add each formatted option
		if (i === 0) {
			// Select "="
			ret.push(`        <option name="Operator${n}" value="${ops[i]}" selected="selected">${ops[i]}</option>`);
		} else {
			ret.push(`        <option name="Operator${n}" value="${ops[i]}">${ops[i]}</option>`);
		}
	}
	ret.push(`      </select>`);
	return ret.join("\n");
}
	
function getValue(n, type) {
	// Returns appropriate input value field
	let ret = "";
	if (type === "TEXT" || type === "INT") {
      ret = `      <input type="text" id="Value${n}" name="Value${n}">`;
	} else if (type === "TINYINT") {
		// Return selected for true/false/na
		ret = `      <select id="Value${n}" name="Value${n}">
        <option name="Value${n}" value="1" selected="selected">1</option>
        <option name="Value${n}" value="0">0</option>
        <option name="Value${n}" value="-1">-1</option>
      </select>`;
	} else {
		// Return number for decimal
		ret = `      <input type="number" name="Value${n}" min="0" step="0.01" >`;
	}
	// Index count after last input is formatted
	count++
	return ret;
}*/

function getColumns(n, table) {
	// Returns select with all columns for given table
	return getTables(n);
}

function addColumns(divname, id) {
	// Adds selector for columns from given table
	if (id != "Empty") {
		let table = document.getElementById(id);
		document.getElementById(divname).innerHTML += getColumns(id.slice(-1), table);
	}
}

function getTables(n) {
	// Returns formatted select field for tables
	let body = searchfield.tables.map(name => {
		return `        <option name="Table${n}" value="${name}">${name}</option>`;
	});
	// Append close and prepend opening
	body.push(`      </select>`);
	body.unshift(`        <option name="Table${n}" value="Empty"></option>`);
	body.unshift(`      <select id="Table${n}" name="Table${n}" onClick="addColumns('searchField', this.id);">`);
	return body.join("\n");
}

function addField(divname) {
	// Adds new table field
	let newdiv = document.createElement("div");
	newdiv.innerHTML = getTables(count);
	document.getElementById(divname).appendChild(newdiv);
	count++;
}
