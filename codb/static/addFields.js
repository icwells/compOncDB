// Contains functions for dynamically presenting logical search options 

var count = 0;

const tables = {
	Patient: [["ID", "INT"],
		["Sex", "TEXT"],
		["Age", "DOUBLE"],
		["Castrated", "TINYINT"],
		["taxa_id", "INT"],
		["source_id", "TEXT"],
		["Date", "TEXT"],
		["Comments", "TEXT"]],
	Diagnosis: [["ID", "INT"],
		["Masspresent", "TINYINT"],
		["Hyperplasia", "TINYINT"],
		["Necropsy", "TINYINT"],
		["Metastasis", "TINYINT"]],
	Tumor: [["ID", "INT"],
		["primary_tumor", "TINYINT"],
		["Malignant", "TINYINT"],
		["Type", "TEXT"],
		["Location", "TEXT"]],
	Taxonomy: [["taxa_id", "INT"],
		["Kingdom", "TEXT"],
		["Phylum", "TEXT"],
		["Class", "TEXT"],
		["Orders", "TEXT"],
		["Family", "TEXT"],
		["Genus", "TEXT"],
		["Species", "TEXT"],
		["Source", "TEXT"]],
	Common: [["taxa_id", "INT"],
		["Name", "TEXT"],
		["Curator", "TEXT"]],
	Life_history: [["taxa_id", "INT"],
		["female_maturity", "DOUBLE"],
		["male_maturity", "DOUBLE"],
		["Gestation", "DOUBLE"],
		["Weaning", "DOUBLE"],
		["Infancy", "DOUBLE"],
		["litter_size", "DOUBLE"],
		["litters_year", "DOUBLE"],
		["interbirth_interval", "DOUBLE"],
		["birth_weight", "DOUBLE"],
		["weaning_weight", "DOUBLE"],
		["adult_weight", "DOUBLE"],
		["growth_rate", "DOUBLE"],
		["max_longevity", "DOUBLE"],
		["metabolic_rate", "DOUBLE"]],
	Totals: [["taxa_id", "INT"],
		["Total", "INT"],
		["Avgage", "DOUBLE"],
		["Adult", "INT"],
		["Male", "INT"],
		["Female", "INT"],
		["Cancer", "INT"],
		["Cancerage", "DOUBLE"],
		["Malecancer", "INT"],
		["Femalecancer", "INT"]],
	Denominators: [["taxa_id", "INT"],
		["Noncancer", "INT"]],
	Source: [["ID", "INT"],
		["service_name", "TEXT"],
		["account_id", "INT"]],
	Accounts: [["account_id", "INT"],
		["Account", "TEXT"],
		["submitter_name", "TEXT"]],
	Unmatched: [["sourceID", "TEXT"],
		["name", "TEXT"],
		["sex", "TEXT"],
		["age", "DOUBLE"],
		["date", "TEXT"],
		["masspresent", "TINYINT"],
		["necropsy", "TINYINT"],
		["comments", "TEXT"],
		["Service", "TEXT"]]
};

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
	if (type === "TEXT") {
      ret = `      <input type="text" id="Value${n}" name="Value${n}">`;
	} else if (type === "TINYINT") {
		// Return selected for true/false/na
		ret = `      <select id="Value${n}" name="Value${n}">
        <option name="Value${n}" value="1" selected="selected">1</option>
        <option name="Value${n}" value="0">0</option>
        <option name="Value${n}" value="-1">-1</option>
      </select>`;
	} else if (type === "INT") {
		ret = `      <input type="number" name="Value${n}" min="0">`;
	} else {
		// Return number for decimal
		ret = `      <input type="number" name="Value${n}" min="0" step="0.01" >`;
	}
	// Index count after last input is formatted
	count++
	return ret;
}

function addValue(divname, id) {
	// Adds operation and value inputs
	if (id != "Empty") {
		let column = document.getElementById(id).value;
		let n = id.slice(-1);
		// let type = 

		document.getElementById(divname).innerHTML += getOperation(n, type);
	}
}

//-----------------------Columns----------------------------------------------

function getColumns(n, table) {
	// Returns select with all columns for given table
	let columns = tables[table];
	console.log(n, table, columns);
	let body = columns.map(name => {
		return `        <option name="Column${n}" value="${name}">${name}</option>`;
	});
	// Append close and prepend opening
	body.push(`      </select>`);
	body.unshift(`        <option name="Column${n}" value="Empty"></option>`);
	body.unshift(`      <select id="Column${n}" name="Column${n}" onchange="addValue('searchField', this.id);">`);
	return body.join("\n");
}

function addColumns(divname, id) {
	// Adds selector for columns from given table
	if (id != "Empty") {
		let table = document.getElementById(id).value;
		// Subset item count from id
		document.getElementById(divname).innerHTML += getColumns(id.slice(-1), table);
	}
}

//-----------------------Tables-----------------------------------------------

function getTables(n) {
	// Returns formatted select field for tables
	let tablenames = Object.keys(tables);
	let body = tablenames.map(name => {
		return `        <option name="Table${n}" value="${name}">${name}</option>`;
	});
	// Append close and prepend opening
	body.push(`      </select>`);
	body.unshift(`        <option name="Table${n}" value="Empty"></option>`);
	body.unshift(`      <select id="Table${n}" name="Table${n}" onchange="addColumns('searchField', this.id);">`);
	return body.join("\n");
}

function addField(divname) {
	// Adds new table field
	let newdiv = document.createElement("div");
	newdiv.innerHTML = getTables(count);
	document.getElementById(divname).appendChild(newdiv);
	count++;
}
