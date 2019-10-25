// Contains expand options function 

function expandCollapse(divname) {
	var x = document.getElementById(divname);
	if (x.style.display === "none") {
		x.style.display = "block";
	} else {
		x.style.display = "none";
	}
}
