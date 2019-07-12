// Contains functions for dynamically presenting logical search options 

var count = 0;

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

const template = () => `      <select id="Table${count}" name="Table${count}" onchange="addColumns('searchField', this.id);">
        <option name="Table${count}" value="Empty"></option>
      </select>
      <select id="Column${count}" name="Column${count}" onchange="addValue('searchField', this.id);">
        <option name="Column${count}" value="Empty"></option>
      </select>
      <select id="Operator${count}" name="Operator${count}">
        <option name="Operator${count}" value="=" selected="selected">=</option>
        <option name="Operator${count}" value="!=">!=</option>
      </select>
      <input type="text" id="Value${count}" name="Value${count}">
      <select id="Select${count}" name="Select${count}">
        <option name="Select${count}" value="1">1</option>
        <option name="Select${count}" value="0">0</option>
        <option name="Select${count}" value="-1">-1</option>
      </select>
`;

//----------------------------------------------------------------------------

class Session {
	// Stores values for single search field
	constructor(id) {
		// Isolate count and input field type 
		this.count = id.slice(-1);
		this.field = id.replace(this.count, "");
		this.table = document.getElementById(`Table${this.count}`).value;
		this.column = document.getElementById(`Column${this.count}`).value;
		this.clearFields();
	}

	clearSelect(id, end) {
		// Removes all but empty line from select field
		let sel = document.getElementById(id)
		for (let i = sel.options.length-1; i > end; i--) {
			// Remove existing element
			sel.remove(i);
		}
	}

	clearFields() {
		// Removes existing and subsequent select fields
		if (this.field === "Table") {
			this.clearSelect(`Column${this.count}`, 0);
			this.clearSelect(`Operator${this.count}`, 1);
		} else if (this.field === "Column") {
			this.clearSelect(`Operator${this.count}`, 1);
		}
	}
}

//-----------------------Values-----------------------------------------------

function addOperations(n) {
	// Adds values to dropdown
	let ops = [">=", "<=", ">", "<"];
	ops.forEach(o => {
		let opt = document.createElement("option");
		opt.text = o;
		document.getElementById(`Operator${n}`).add(opt);
	});
}

function toggleInputs(type, n) {
	// Switches between select/text input
	let value = document.getElementById(`Value${n}`);
	let select = document.getElementById(`Select${n}`);
	if (type === "TINYINT") {
		// Reveal select block
		select.style.display = "inline";
		value.style.display = "none";
	} else {
		// Reveal text input
		select.style.display = "none";
		value.style.display = "inline";
		value.type = "text";
	}
}

function getDataType(table, column) {
	// Returns data type of given column value
	let t = tables[table];
	for (let i in t) {
		if (t[i][0] === column) {
			return t[i][1];
		} 
	}
	return null;
}

function addValue(divname, id) {
	// Adds operation and value inputs
	let s = new Session(id);
	let type = getDataType(s.table, s.column);
	if (type != null) {
		let val = `Value${s.count}`;
		let sel = `Select${s.count}`;
		toggleInputs(type, s.count);
		if (type === "INT" || type === "DOUBLE") {
			// Format number input for decimal/integer
			addOperations(s.count);
			document.getElementById(val).type = "number";
			document.getElementById(val).min = "0";
			if (type === "DOUBLE") {
				// Add step
				document.getElementById(val).step = "0.01";
			}
		}
	}
}

function addColumns(divname, id) {
	// Adds selector for columns from given table
	let s = new Session(id);
	if (document.getElementById(id).value != "Empty") {
		if (s.table) {
			tables[s.table].forEach(name => {
				let opt = document.createElement("option");
				opt.text = name[0];
				document.getElementById(`Column${s.count}`).add(opt);
			});
		}		
	}
}

//--------------------------NewRow--------------------------------------------

function getTables(divname, n) {
	// Formats select field for tables
	let tableselect = document.getElementById(`Table${n}`);
	let tablenames = Object.keys(tables);
	tablenames.forEach(name => {
		let opt = document.createElement("option");
		opt.text = name;
		tableselect.add(opt);
	});
}

function addRow(divname) {
	// Adds new row to table search
	if (count < 10) {
		// Limit to ten search fields
		if (count === 0) {
			// Insert template into existing div
			document.getElementById(divname).innerHTML = template();
		} else {
			// Create new div
			let newdiv = document.createElement("div");
			newdiv.innerHTML = template();
			document.getElementById(divname).appendChild(newdiv);
		}
		// Populate tables list
		getTables(divname, count);
		// Hide value select
		toggleInputs("TEXT", count);
		count++;
	}
}
