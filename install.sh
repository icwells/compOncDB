#!/bin/bash

##############################################################################
# This script will install scripts for the compOncDB package.
# 
# Required programs:	Go 1.11+
##############################################################################

DBI="github.com/icwells/dbIO"
DBU="dbupload"
DBE="dbextract"
DR="github.com/go-sql-driver/mysql"
IO="github.com/icwells/go-tools/iotools"
KP="gopkg.in/alecthomas/kingpin.v2"
MAIN="compOncDB"
PARSE="parseRecords"
PR="github.com/Songmu/prompter"
SA="github.com/icwells/go-tools/strarray"

# Get install location
SYS=$(ls $GOPATH/pkg | head -1)
PDIR=$GOPATH/pkg/$SYS

echo ""
echo "Preparing compOncDB package..."
echo "GOPATH identified as $GOPATH"
echo ""

# Get dependencies
for I in $DR $IO $KP $PR $SA $DBI; do
	if [ ! -e "$PDIR/$I.a" ]; then
		echo "Installing $I..."
		go get -u $I
		echo ""
	fi
done

for I in $DBU $DBE; do
	echo "Installing $I..."
	cp -r "src/$I" $GOPATH/src/
	go install $I
	echo ""
done

# Intall parseRecords
echo "Building $PARSE..."
go build -o bin/$PARSE src/$PARSE/*.go
echo ""

# compOncDB 
echo "Building main..."
go build -o bin/$MAIN src/$MAIN/*.go

echo "Finished"
echo ""
