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

function toSelect(id) {
	// Changes value input to select
	let opts = ["1", "0", "-1"];
	let old = document.getElementById(id);
	old.parentElement.removeChild(old);
	let sel = document.createElement("select");
	opts.forEach(o => {
		let opt = document.createElement("option");
		opt.text = o;
		sel.add(opt);
	});
}

function addValue(divname, id) {
	// Adds operation and value inputs
	let s = new Session(id);
	let type = getDataType(s.table, s.column);
	//console.log(s.count, s.column, s.table, type);
	if (type != null) {
		if (type === "TEXT") {
			document.getElementId(id).type = "text";
		} else if (type === "TINYINT") {
			// Convert to select for true/false/na
			toSelect(id)
		} else {
			// Format number input for decimal/integer
			addOperations(n);
			document.getElementId(id).type = "number";
			document.getElementId(id).min = "0";
			if (type === "DOUBLE") {
				// Add step
				document.getElementId(id).step = "0.01";
			}
		}
	}
}

function addColumns(divname, id) {
	// Adds selector for columns from given table
	let s = new Session(id);
	if (s.table) {
		tables[s.table].forEach(name => {
			let opt = document.createElement("option");
			opt.text = name[0];
			document.getElementById(`Column${s.count}`).add(opt);
		});
	}
}

//--------------------------NewRow--------------------------------------------

function getTables(divname, n) {
	// Formats select field for tables
	let tablenames = Object.keys(tables);
	tablenames.forEach(name => {
		let opt = document.createElement("option");
		opt.text = name;
		document.getElementById(`Table${n}`).add(opt);
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
		count++;
	}
}
