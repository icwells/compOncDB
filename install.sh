#!/bin/bash

##############################################################################
# This script will install scripts for the compOncDB package.
# 
# Required programs:	Go 1.11+
##############################################################################

MAIN="compOncDB"

installPackages () {
	echo "Installing dependencies..."
	ETREE="github.com/beevik/etree"
	DATAFRAME="github.com/icwells/go-tools/dataframe"
	FRACTION="github.com/icwells/go-tools/fraction"
	IOTOOLS="github.com/icwells/go-tools/iotools"
	STRARRAY="github.com/icwells/go-tools/strarray"
	KINGPIN="gopkg.in/alecthomas/kingpin.v2"
	MUX="github.com/gorilla/mux"
	SCHEMA="github.com/gorilla/schema"
	COOKIE="github.com/gorilla/securecookie"
	SEESSIONS="github.com/gorilla/sessions"
	DBIO="github.com/icwells/dbIO"
	SIMPLESET="github.com/icwells/simpleset"
	FUZZY="github.com/lithammer/fuzzysearch/fuzzy"
	ASPELL="github.com/trustmaster/go-aspell"
	for I in $ETREE $DATAFRAME $FRACTION $IOTOOLS $STRARRAY $KINGPIN $MUX $SCHEMA $COOKIE $SEESSIONS $DBIO $SIMPLESET $FUZZY $ASPELL; do
		go get $I
	done
}

installMain () {
	# compOncDB 
	echo "Building $MAIN..."
	go build -i -o $GOBIN/$MAIN src/*.go
	echo ""
}

echo ""
echo "Preparing compOncDB package..."
echo "GOPATH identified as $GOPATH"
echo ""

if [ $# -eq 0 ]; then
	installMain
elif [ $1 = "main" ]; then
	installMain
elif [ $1 = "all" ]; then
	installPackages
	installMain
elif [ $1 = "help" ]; then
	echo "Installs Go scripts for compOnDB"
	echo ""
	echo "main	Installs scripts to GOBIN only."
	echo "app	Installs web application only."
	echo "all	Installs all scripts and dependencies."
	echo "help	Prints help text and exits."
	echo ""
fi

echo "Finished"
echo ""
