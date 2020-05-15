#!/bin/bash

##############################################################################
# This script will install scripts for the compOncDB package.
# 
# Required programs:	Go 1.11+
##############################################################################

APP="codbApplication"
MAIN="compOncDB"

installMain () {
	# compOncDB 
	echo "Building $MAIN..."
	go build -i -o $GOBIN/$MAIN src/*.go
	echo ""
}

installApp () {
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
	installApp
elif [ $1 = "main" ]; then
	installMain
elif [ $1 = "app" ]; then
	installApp
elif [ $1 = "help" ]; then
	echo "Installs Go scripts for compOnDB"
	echo ""
	echo "main	Installs scripts to GOBIN only."
	echo "app	Installs web application only."
	echo "help	Prints help text and exits."
	echo ""
fi

echo "Finished"
echo ""
