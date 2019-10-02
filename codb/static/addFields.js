// Contains functions for dynamically presenting logical search options 

let COUNT = 0;
let SEARCHES = {};

const TABLES = {
	Patient: [["ID", "INT"], ["Sex", "TEXT"], ["Age", "DOUBLE"], ["Castrated", "TINYINT"], ["taxa_id", "INT"], ["source_id", "TEXT"], ["source_name", "TEXT"], ["Date", "TEXT"], ["Comments", "TEXT"]],
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

const ROW = (i, n) => `    <div id="parameter${n}">
      <select id="${i}Table${n}" name="${i}Table${n}" onchange="SEARCHES[${i}].params[${n}].addColumns();">
        <option name="${i}Table${n}" value="Empty"></option>
      </select>
      <select id="${i}Column${n}" name="${i}Column${n}" onchange="SEARCHES[${i}].params[${n}].addValue();">
        <option name="${i}Column${n}" value="Empty"></option>
      </select>
      <select id="${i}Operator${n}" name="${i}Operator${n}">
        <option name="${i}Operator${n}" value="=" selected="selected">=</option>
        <option name="${i}Operator${n}" value="!=">!=</option>
      </select>
      <input type="text" id="${i}Value${n}" name="${i}Value${n}">
      <select id="${i}Select${n}" name="${i}Select${n}">
        <option name="${i}Select${n}" value="1">1</option>
        <option name="${i}Select${n}" value="0">0</option>
        <option name="${i}Select${n}" value="-1">-1</option>
      </select>
      <input type="button" value="Remove" onClick="SEARCHES[${i}].removeRow('${n}');">
    </div>`;

const SEARCH = (n) => `    <div id="search${n}">
	</div>
    <div id="options${n}">
      <input type="button" value="Add search parameter" onClick="SEARCHES[${n}].addRow();">
      <input type="button" value="Remove Search" onClick="removeSearch('${n}');">
	<hr>
    </div>`;

//----------------------Parameter---------------------------------------------

class Parameter {
	// Stores values for single search field
	constructor(i, n) {
		this.num = i;
		this.count = n;
		this.table = "Empty";
		this.column = "Empty";
		this.type = "TEXT";
	}

	clearSelect(id, end) {
		// Removes all but empty line from select field
		let sel = document.getElementById(id)
		for (let i = sel.options.length-1; i > end; i--) {
			// Remove existing element
			sel.remove(i);
		}
	}

	clearFields(field) {
		// Removes existing and subsequent select fields
		if (field === "Table") {
			this.clearSelect(`${this.num}Column${this.count}`, 0);
			this.clearSelect(`${this.num}Operator${this.count}`, 1);
		} else if (field === "Column") {
			this.clearSelect(`${this.num}Operator${this.count}`, 1);
		}
	}

	addOperations() {
		// Adds values to dropdown
		let ops = [">=", "<=", ">", "<"];
		ops.forEach(o => {
			let opt = document.createElement("option");
			opt.text = o;
			document.getElementById(`${this.num}Operator${this.count}`).add(opt);
		});
	}

	toggleInputs() {
		// Switches between select/text input
		let value = document.getElementById(`${this.num}Value${this.count}`);
		let select = document.getElementById(`${this.num}Select${this.count}`);
		if (this.type === "TINYINT") {
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

	getDataType() {
		// Returns data type of given column value
		let t = TABLES[this.table];
		for (let i in t) {
			if (t[i][0] === this.column) {
				this.type = t[i][1];
				break;
			} 
		}
	}

	addValue() {
		// Adds operation and value inputs
		this.clearFields("Column");
		this.column = document.getElementById(`${this.num}Column${this.count}`).value;
		if (this.column != "Empty") {
			this.getDataType();
			this.toggleInputs();
			if (this.type === "INT" || this.type === "DOUBLE") {
				// Format number input for decimal/integer
				this.addOperations();
				let elem = document.getElementById(`${this.num}Value${this.count}`);
				elem.type = "number";
				elem.min = "0";
				if (this.type === "DOUBLE") {
					// Add step
					elem.step = "0.01";
				}
			}
		}
	}

	addColumns() {
		// Adds selector for columns from given table
		this.clearFields("Table");
		this.table = document.getElementById(`${this.num}Table${this.count}`).value;
		if (this.table != "Empty") {
			TABLES[this.table].forEach(name => {
				let opt = document.createElement("option");
				opt.text = name[0];
				document.getElementById(`${this.num}Column${this.count}`).add(opt);
			});		
		}
	}

	getTables() {
		// Formats select field for tables
		let tableselect = document.getElementById(`${this.num}Table${this.count}`);
		setTableNames(tableselect);
	}
}

//--------------------------Search--------------------------------------------

class Search {
	constructor(n) {
		this.num = n;
		this.count = 0;
		this.params = {};
	}

	removeRow(n) {
		// Removes row from document and map
		let child = document.getElementById(`parameter${n}`);
		let parent = child.parentNode;
		parent.removeChild(child);
		delete this.params[n];
	}

	addRow() {
		// Adds new row to table search
		if (this.count < 10) {
			// Create new div
			let newdiv = document.createElement("div");
			newdiv.innerHTML = ROW(this.num, this.count);
			document.getElementById(`search${this.num}`).appendChild(newdiv);
			// Get new object, populate tables list, hide value select, and add to map
			let p = new Parameter(this.num, this.count);
			p.getTables();
			p.toggleInputs();
			this.params[this.count] = p;
			this.count++;
		}
	}
}

//--------------------------Functions-----------------------------------------

function setTableNames(tableselect) {
	// Assigns table names as options to given select
	let tablenames = Object.keys(TABLES);
	tablenames.forEach(name => {
		let opt = document.createElement("option");
		opt.text = name;
		tableselect.add(opt);
	});
}

function newSearch(divname) {
	// Adds new search to form
	if (COUNT < 10) {
		if (COUNT === 0){
			// Add options to single table select
			let tableselect = document.getElementById("Table");
			setTableNames(tableselect);
		}
		// Create new div
		let newdiv = document.createElement("div");
		newdiv.innerHTML = SEARCH(COUNT);
		document.getElementById("SearchBlock").appendChild(newdiv);
		// Create new search object
		let s = new Search(COUNT);
		s.addRow();
		SEARCHES[COUNT] = s;
		COUNT++;
	}
}

function removeSearch(n) {
	// Removes search from document and map
	let child = document.getElementById(`search${n}`);
	let opt = document.getElementById(`options${n}`);
	let parent = child.parentNode;
	parent.removeChild(child);
	parent.removeChild(opt);
	delete SEARCHES[n];
}
