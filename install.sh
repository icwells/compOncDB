#!/bin/bash

##############################################################################
# This script will install scripts for the compOncDB package.
# 
# Required programs:	Go 1.7+
##############################################################################

DBS="dbsearch"
DR="github.com/go-sql-driver/mysql"

# Get install location
SYS=$(ls $GOPATH/pkg | head -1)
PDIR=$GOPATH/pkg/$SYS

echo ""
echo "Preparing compOncDB package..."
echo "GOPATH identified as $GOPATH"
echo ""

# Get mysql driver
if [ ! -e "$GOPATH/src/$DR.a" ]; then
	echo "Installing $DR..."
	go get -u $DR
fi

# Install dbsearch
#if [ ! -e "$PDIR/$DBS.a" ]; then
cp -R $DBS $GOPATH/src/gopkg.in/
go install $DBS
#fi
