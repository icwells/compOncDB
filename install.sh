#!/bin/bash

##############################################################################
# This script will install scripts for the compOncDB package.
# 
# Required programs:	Go 1.7+
##############################################################################

DBI="dbIO"
DR="github.com/go-sql-driver/mysql"
IO="github.com/icwells/go-tools/iotools"
PR="github.com/Songmu/prompter"
SA="github.com/icwells/go-tools/strarray"

# Get install location
SYS=$(ls $GOPATH/pkg | head -1)
PDIR=$GOPATH/pkg/$SYS

echo ""
echo "Preparing compOncDB package..."
echo "GOPATH identified as $GOPATH"
echo ""

# Get mysql driver
for I in $DR $IO $PR $SA; do
	if [ ! -e "$GOPATH/src/$I.a" ]; then
		echo "Installing $I..."
		go get -u $I
	fi
done

# Install dbIO
#if [ ! -e "$PDIR/$DBI.a" ]; then
cp -R $DBI $GOPATH/src/gopkg.in/
go install $DBI
#fi
