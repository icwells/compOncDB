// Contains functions for dynamically presenting logical search options 

var count = 0;
const addrow = "addRow";
const tableindeces = ["Patient", "Diagnosis", "Tumor"];

const tables = {
	Patient: [["ID", "INT"], ["Sex", "TEXT"], ["Age", "DOUBLE"], ["Castrated", "TINYINT"], ["taxa_id", "INT"], ["source_id", "TEXT"], ["Date", "TEXT"], ["Comments", "TEXT"]],
	Diagnosis: [["ID", "INT"], ["Masspresent", "TINYINT"], ["Hyperplasia", "TINYINT"], ["Necropsy", "TINYINT"], ["Metastasis", "TINYINT"]],
	Tumor: [["ID", "INT"], ["primary_tumor", "TINYINT"], ["Malignant", "TINYINT"], ["Type", "TEXT"], ["Location", "TEXT"]],
	Taxonomy: [["taxa_id", "INT"], ["Kingdom", "TEXT"], ["Phylum", "TEXT"], ["Class", "TEXT"], ["Orders", "TEXT"], ["Family", "TEXT"], ["Genus", "TEXT"], ["Species", "TEXT"], ["Source", "TEXT"]],
	Common: [["taxa_id", "INT"], ["Name", "TEXT"], ["Curator", "TEXT"]],
	Life_history: [["taxa_id", "INT"], ["female_maturity", "DOUBLE"], ["male_maturity", "DOUBLE"], ["Gestation", "DOUBLE"], ["Weaning", "DOUBLE"], ["Infancy", "DOUBLE"], ["litter_size", "DOUBLE"], ["litters_year", "DOUBLE"], ["interbirth_interval", "DOUBLE"], ["birth_weight", "DOUBLE"], ["weaning_weight", "DOUBLE"], ["adult_weight", "DOUBLE"], ["growth_rate", "DOUBLE"], ["max_longevity", "DOUBLE"], ["metabolic_rate", "DOUBLE"]],
	Totals: [["taxa_id", "INT"], ["Total", "INT"], ["Avgage", "DOUBLE"], ["Adult", "INT"], ["Male", "INT"], ["Female", "INT"], ["Cancer", "INT"], ["Cancerage", "DOUBLE"], ["Malecancer", "INT"], ["Femalecancer", "INT"]],
	Denominators: [["taxa_id", "INT"], ["Noncancer", "INT"]],
	Source: [["ID", "INT"], ["service_name", "TEXT"], ["account_id", "INT"]],
	Accounts: [["account_id", "INT"], ["Account", "TEXT"], ["submitter_name", "TEXT"]],
	Unmatched: [["sourceID", "TEXT"], ["name", "TEXT"], ["sex", "TEXT"], ["age", "DOUBLE"], ["date", "TEXT"], ["masspresent", "TINYINT"], ["necropsy", "TINYINT"], ["comments", "TEXT"], ["Service", "TEXT"]]
};

//-----------------------Session----------------------------------------------

class Session {
	// Stores values for single search field
	constructor(n) {
		this.count = n;
		this.table = localStorage.getItem(`Table${this.count}`);
		this.column = localStorage.getItem(`Column${this.count}`);
		this.operator = localStorage.getItem(`Operator${this.count}`);
		this.value = localStorage.getItem(`Value${this.count}`);
		console.log(this.count, this.table, this.column);
	}

	setSelected(id, option) {
		document.getElementById(id).selectedIndex = 1;
	}

	setField(id) {
		// Returns value at given field and stores locally
		let val = document.getElementById(id).value;
		localStorage.setItem(id, val);
		this.setSelected(id, val);
		return val
	}

	setTable(id) {
		// Gets selected table from form
		this.table = this.setField(id);
	}

	setColumn(id) {
		// Gets selected column from form
		this.column = this.setField(id);
	}

	setValue(id) {
		// Gets selected operator and value from form
		this.operator = this.setField(`Operator${this.count}`);
		this.value = this.setField(id);
	}
}

//-----------------------Values-----------------------------------------------

function getOperation(n, type) {
	// Returns formatted operation block for given data type
	let ops = ["=", "!=", ">=", "<=", ">", "<"];
	let end = 2;
	let ret = [`      <select id="Operator${n}" name="Operator${n}">`];
	if (type === "DOUBLE" || type === "INT") {
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
	let action = `onchange="addField('searchField', this.id);"`;
	if (type === "TEXT") {
      ret = `      <input type="text" id="Value${n}" name="Value${n}" ${action}>`;
	} else if (type === "TINYINT") {
		// Return selected for true/false/na
		ret = `      <select id="Value${n}" name="Value${n} ${action}">
        <option name="Value${n}" value="1" selected="selected">1</option>
        <option name="Value${n}" value="0">0</option>
        <option name="Value${n}" value="-1">-1</option>
      </select>`;
	} else if (type === "INT") {
		ret = `      <input type="number" name="Value${n}" min="0" ${action}>`;
	} else {
		// Return number for decimal
		ret = `      <input type="number" name="Value${n}" min="0" step="0.01" ${action}>`;
	}
	ret += "<br/>"
	// Index count after last input is formatted
	count++
	return ret;
}

function getDataType(table, column) {
	// Returns data type of given column value
	for (let i in tables[table]) {
		//console.log(i, column, table);
		if (i[0] === column) {
			return i[1];
		} 
	}
	return null;
}

function addValue(divname, s) {
	// Adds operation and value inputs
	console.log(s.count, s.column, s.table);
	let type = getDataType(s.table, s.column);
	if (type != null) {
		document.getElementById(divname).innerHTML += getOperation(self.count, type);
		document.getElementById(divname).innerHTML += getValue(self.count, type);
		// Enable add field button
		document.getElementById(addrow).disabled = null;
	}
}

//-----------------------Tables-----------------------------------------------

function addColumns(divname, s) {
	// Adds selector for columns from given table
	if (document.getElementById(`Column${s.count}`)) {
		// Replace existing elment
		document.getElementById(`Column${s.count}`).outerHTML = "";
	}
	let body = tables[s.table].map(name => {
		return `        <option name="Column${s.count}" value="${name[0]}">${name[0]}</option>`;
	});
	// Append close and prepend opening
	body.push(`      </select>`);
	body.unshift(`        <option name="Column${s.count}" value="Empty"></option>`);
	body.unshift(`      <select id="Column${s.count}" name="Column${s.count}" onchange="addValue('searchField', this.id);">`);
	document.getElementById(divname).innerHTML += body.join("\n");
}

function getTables(divname, n) {
	// Formats select field for tables
	let tablenames = Object.keys(tables);
	let body = tablenames.map(name => {
		return `        <option name="Table${n}" value="${name}">${name}</option>`;
	});
	// Append close and prepend opening
	body.push(`      </select>`);
	body.unshift(`        <option name="Table${n}" value="Empty"></option>`);
	body.unshift(`      <select id="Table${n}" name="Table${n}" onchange="addField('searchField', this.id);">`);
	// Add select field to div and index count
	document.getElementById(divname).innerHTML = body.join("\n");
	// Disable button
	document.getElementById(addrow).disabled = true;
	count++;
}

function addField(divname, id = null) {
	// Adds new table field
	if (id == null && count < 10) {
		// Limit to ten search fields
		getTables(divname, count);
	} else if (id != "Empty") {
		// Isolate count and input field type
		let n = id.slice(-1);
		let field = id.replace(n, "");
		let session = new Session(n);
		if (field === "Table") {
			session.setTable(id)
			addColumns(divname, session);
		} else if (field === "Column") {
			session.setColumn(id)
			console.log(session.table, session.column);
			addValue(divname, session);
		} else {
			// Record selected values
			session.setValue(id)
		}
	}
}
