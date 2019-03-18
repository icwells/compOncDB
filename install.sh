#!/bin/bash

##############################################################################
# This script will install scripts for the compOncDB package.
# 
# Required programs:	Go 1.11+
##############################################################################

DBI="github.com/icwells/dbIO"
DBU="github.com/icwells/compOncDB/src/dbupload"
DBE="github.com/icwells/compOncDB/src/dbextract"
DR="github.com/go-sql-driver/mysql"
FZ="github.com/lithammer/fuzzysearch/fuzzy"
IO="github.com/icwells/go-tools/iotools"
KP="gopkg.in/alecthomas/kingpin.v2"
MAIN="compOncDB"
PARSE="parseRecords"
PR="github.com/Songmu/prompter"
SA="github.com/icwells/go-tools/strarray"

# Get install location
SYS=$(ls $GOPATH/pkg | head -1)
PDIR=$GOPATH/pkg/$SYS

installDependencies () {
# Get dependencies
	for I in $DR $FZ $IO $KP $PR $SA ; do
		if [ ! -e "$PDIR/$I.a" ]; then
			echo "Installing $I..."
			go get -u $I
			echo ""
		fi
	done
}

installDBIO () {
	for I in $DBI $DBU $DBE; do
		echo "Installing $I..."
		go get -u $I
		echo ""
	done
}

installMain () {
	# Install parseRecords
	echo "Building $PARSE..."
	go build -i -o bin/$PARSE src/$PARSE/*.go
	echo ""

	# compOncDB 
	echo "Building main..."
	go build -i -o bin/$MAIN src/$MAIN/*.go
}

echo ""
echo "Preparing compOncDB package..."
echo "GOPATH identified as $GOPATH"
echo ""

if [ $# -eq 0 ]; then
	installMain
elif [ $1 = "all" ]; then
	installDependencies
	installDBIO
	installMain
elif [ $1 = "db" ]; then
	installDBIO
	installMain
elif [ $1 = "help" ]; then
	echo "Installs Go scripts for compOnDB"
	echo ""
	echo "all	Installs scripts and all depenencies."
	echo "db	Installs scripts and dbIO only."
	echo "help	Prints help text and exits."
	echo ""
fi

echo "Finished"
echo ""
