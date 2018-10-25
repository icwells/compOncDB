#!/bin/bash

##############################################################################
# This script will install scripts for the compOncDB package.
# 
# Required programs:	Go 1.7+
##############################################################################

DBI="dbIO"
DR="github.com/go-sql-driver/mysql"
IO="github.com/icwells/go-tools/iotools"
KP="gopkg.in/alecthomas/kingpin.v2"
LAN="golang.org/x/text/language"
MAIN="compOncDB"
PARSE="parseRecords"
PR="github.com/Songmu/prompter"
SA="github.com/icwells/go-tools/strarray"
TS="golang.org/x/text/search"

# Get install location
SYS=$(ls $GOPATH/pkg | head -1)
PDIR=$GOPATH/pkg/$SYS

echo ""
echo "Preparing compOncDB package..."
echo "GOPATH identified as $GOPATH"
echo ""

# Get dependencies
for I in $DR $IO $KP $LAN $PR $SA $TS; do
	echo "Installing $I..."
	go get -u $I
	echo ""
done

# Intall parseRecords
echo "Building $PARSE..."
go build -o bin/$PARSE src/$PARSE/*.go
echo ""

# Install dbIO
for I in $DBI; do
	echo "Installing $I..."
	cp -R src/$I/ $GOPATH/src/
	go install $I
	echo ""
done

# compOncDB 
echo "Building main..."
go build -o bin/$MAIN src/$MAIN/*.go

echo "Finished"
echo ""
