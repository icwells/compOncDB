// Defines search struct for compOncDB web app

package main

import (
	""
)

type SearchForm struct {
	Column		string
	Operator	string
	Value		string
	Taxon		bool
	Table		string
	Dump		bool
	Summary		bool
	Cancerrate	bool
	Min			int
	Necropsy	bool
	Common		bool
	Count		bool
	Infant		bool
	search		bool
}
