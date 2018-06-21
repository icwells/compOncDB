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
MAIN="coDB"
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

# Get mysql driver
for I in $DR $IO $KP $LAN $PR $SA $TS; do
	if [ ! -e "$PDIR/$I.a" ]; then
		echo "Installing $I..."
		go get -u $I
		echo ""
	fi
done

# Install dbIO
#if [ ! -e "$PDIR/$DBI.a" ]; then
cp -R src/$DBI/ $GOPATH/src/
go install $DBI
echo ""
#fi

# lineageSimulator 
echo "Building main..."
go build src/$MAIN/*.go
mv $MAIN bin/

echo "Finished"
echo ""
