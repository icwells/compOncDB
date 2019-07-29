#!/bin/bash

##############################################################################
#	Manages black box tests on the comaprative oncology mysql database
#	Make sure a test database exists prior to running.
#
#	Required:	go 1.10+
#				mysql 14.14+
#
#	Usage:		./test.sh {help/...}
##############################################################################
USER=""
PW=""
WD=$(pwd)

APP="$WD/codb/*.go"
CDB="compOncDB"
CUSRC="$WD/src/codbutils/*.go"
DBSRC="$WD/src/*.go"
DUSRC="$WD/src/dbupload/*.go"
DESRC="$WD/src/dbextract/*.go"
PRSRC="$WD/src/parserecords/*.go"
TSTDIR="$WD/test/*.go"

getUser () {
	# Reads mysql user name and password from command line
	read -p "Enter MySQL username: " USER
	echo -n "Enter MySQL password: "
	read -s PW
	echo ""
}

whiteBoxTests () {
	echo ""
	echo "Running white box tests..."
	go test $PRSRC
	go test $CUSRC
	go test $DUSRC
	go test $DESRC
	#go test $APP
}

testParseRecords () {
	echo ""
	echo "Running black box tests on parseRecords..."
	go test $TSTDIR --run TestParseRecords
}

testUpload () {
	# Upload test data
	echo ""
	echo "Running black box tests on database upload..."
	# Compare tables to expected
	go test $TSTDIR --run TestUpload --args --user=$USER --password=$PW
}

testSearch () {
	# Test search output
	echo ""
	echo "Running black box tests on database search..."
	go test $TSTDIR --run TestSearches --args --user=$USER --password=$PW
}

testUpdates () {
	# Test search output
	echo ""
	echo "Running black box tests on database update..."
	go test $TSTDIR --run TestUpdates --args --user=$USER --password=$PW
}

blackBoxTests () {
	testParseRecords
	testUpload
	testSearch
	testUpdates
}

checkSource () {
	# Runs go fmt/vet on source files (vet won't run in loop)
	echo ""
	echo "Running go $1..."
	go $1 $APP
	go $1 $PRSRC
	go $1 $DBSRC
	go $1 $DUSRC
	go $1 $DESRC
	go $1 $PRSRC
	go $1 $TSTDIR
}

helpText () {
	echo ""
	echo "Runs test scripts for compOncDB. Omit command line arguments to run all tests."
	echo "Usage: ./test.sh {install/white/parse/upload/search/update}"
	echo ""
	echo "all		Runs all tests"
	echo "whitebox		Runs white box tests"
	echo "blackbox		Runs all black box tests (parse, upload, search, and update)"
	echo "parse		Runs parseRecords black box tests"
	echo "upload		Runs compOncDB upload black box tests"
	echo "search		Runs compOncDB search black box tests"
	echo "update		Runs compOncDB update black box tests"
	echo "fmt		Runs go fmt on all source files."
	echo "vet		Runs go vet on all source files."
	echo "username	MySQL username (root by default)."
}

if [ $# -eq 0 ]; then
	helpText
elif [ $1 = "all" ]; then
	getUser
	whiteBoxTests
	blackBoxTests
elif [ $1 = "whitebox" ]; then
	whiteBoxTests
elif [ $1 = "blackbox" ]; then
	getUser
	blackBoxTests
elif [ $1 = "parse" ]; then
	testParseRecords
elif [ $1 = "upload" ]; then
	getUser
	testUpload
elif [ $1 = "search" ]; then
	getUser
	testSearch
elif [ $1 = "update" ]; then
	getUser
	testUpdates
elif [ $1 = "fmt" ]; then
	checkSource $1
elif [ $1 = "vet" ]; then
	checkSource $1
elif [ $1 = "help" ]; then
	helpText
else
	helpText
fi
echo ""
