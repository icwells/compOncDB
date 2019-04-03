#!/bin/bash

##############################################################################
# This script will install scripts for the compOncDB package.
# 
# Required programs:	Go 1.11+
##############################################################################

AP="github.com/trustmaster/go-aspell"
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

installPackage () {
	# Installs go package if it is not present in src directory
	echo "Installing $1..."
	go get -u $1
	echo ""
}

installDependencies () {
# Get dependencies
	for I in $AP $DR $FZ $IO $KP $PR $SA ; do
		if [ ! -e "$PDIR/$1.a" ]; then
			installPackage $I
		fi
	done
}

installMain () {
	# Install parseRecords
	echo "Building $PARSE..."
	go build -i -o bin/$PARSE src/$PARSE/*.go
	echo ""

	for I in $DBU $DBE; do
		# Install submodules
		installPackage $I
	done

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
	installPackage $DBI
	installMain
elif [ $1 = "db" ]; then
	installPackage $DBI
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
