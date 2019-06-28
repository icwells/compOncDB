#!/bin/bash

##############################################################################
# This script will install scripts for the compOncDB package.
# 
# Required programs:	Go 1.11+
##############################################################################

AP="github.com/trustmaster/go-aspell"
APP="codbApplication"
DBI="github.com/icwells/dbIO"
DR="github.com/go-sql-driver/mysql"
FZ="github.com/lithammer/fuzzysearch/fuzzy"
GM="github.com/gorilla/mux"
GC="github.com/gorilla/securecookie"
GS="github.com/gorilla/sessions"
GSC="github.com/gorilla/schema"
IO="github.com/icwells/go-tools/iotools"
KP="gopkg.in/alecthomas/kingpin.v2"
MAIN="compOncDB"
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
	for I in $AP $DBI $DR $FZ $GM $GC $GS $GSC $IO $KP $PR $SA ; do
		if [ ! -e "$PDIR/$1.a" ]; then
			installPackage $I
		fi
	done
}

installMain () {
	# compOncDB 
	echo "Building $MAIN..."
	go build -i -o $GOBIN/$MAIN src/*.go
	echo ""

	# Application
	echo "Building $APP..."
	cd codb/
	go build -i -o $APP *.go
}

echo ""
echo "Preparing compOncDB package..."
echo "GOPATH identified as $GOPATH"
echo ""

if [ $# -eq 0 ]; then
	installMain
elif [ $1 = "all" ]; then
	installDependencies
	installMain
elif [ $1 = "test" ]; then
	installDependencies
elif [ $1 = "db" ]; then
	installPackage $DBI
	installMain
elif [ $1 = "help" ]; then
	echo "Installs Go scripts for compOnDB"
	echo ""
	echo "all	Installs scripts and all depenencies."
	echo "test	Installs depenencies only (for white box testing)."
	echo "db	Installs scripts and dbIO only."
	echo "help	Prints help text and exits."
	echo ""
fi

echo "Finished"
echo ""
