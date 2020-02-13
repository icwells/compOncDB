// Contains expand options and formSubmit function 

function expandCollapse(divname) {
	// Toggles div between hidden and visible
	var x = document.getElementById(divname);
	if (x.style.display === "none" || x.style.display === "") {
		x.style.display = "block";
	} else {
		x.style.display = "none";
	}
}

function displayDiv(divname) {
	// Opens div
	var x = document.getElementById(divname);
	x.style.display = "block";
}

function hideDiv(divname) {
	// Closes div
	var x = document.getElementById(divname);
	x.style.display = "none";
}
