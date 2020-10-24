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
CNRT="$WD/src/cancerrates/*.go"
CLSRC="$WD/src/clusteraccounts/*.go"
CUSRC="$WD/src/codbutils/*.go"
DBSRC="$WD/src/*.go"
DUSRC="$WD/src/dbupload/*.go"
DESRC="$WD/src/dbextract/*.go"
DIAG="$WD/src/diagnoses/*.go"
PRSRC="$WD/src/parserecords/*.go"
TSTDIR="$WD/test/*.go"

getUser () {
	# Reads mysql user name and password from command line
	read -p "Enter MySQL username: " USER
	echo -n "Enter MySQL password: "
	read -s PW
	echo ""
	ARGS="--args --user=$USER --password=$PW"
}

whiteBoxTests () {
	echo ""
	echo "Running white box tests..."
	go test $CNRT
	go test $CLSRC
	go test $DIAG
	go test $PRSRC
	go test $CUSRC
	go test $DUSRC
	go test $DESRC
}

testParseRecords () {
	echo ""
	echo "Running black box tests on parseRecords..."
	go test $TSTDIR --run TestParseRecords
}

testCancerRates () {
	# Tests cancer rate calculation with clean upload
	echo ""
	echo "Running black box tests on database upload..."
	go test $TSTDIR --run TestUpload $ARGS
	echo ""
	echo "Running black box tests on cancer rate calculation..."
	go test $TSTDIR --run TestCancerRates $ARGS
}

testDataBase () {
	# Installs and tests database functions
	testCancerRates

	echo ""
	echo "Running black box tests on database filtering..."
	go test $TSTDIR --run TestFilterPatients $ARGS

	echo ""
	echo "Running black box tests on database search..."
	go test $TSTDIR --run TestSearches $ARGS

	echo ""
	echo "Running black box tests on database update..."
	go test $TSTDIR --run TestUpdates $ARGS

	echo ""
	echo "Running black box tests on database deletion..."
	go test $TSTDIR --run TestDelete $ARGS
}

checkSource () {
	# Runs go fmt/vet on source files (vet won't run in loop)
	echo ""
	echo "Running go $1..."
	go $1 $APP
	go $1 $CNRT
	go $1 $CNRT
	go $1 $CLSRC
	go $1 $DIAG
	go $1 $PRSRC
	go $1 $DBSRC
	go $1 $DUSRC
	go $1 $DESRC
	go $1 $PRSRC
	go $1 $TSTDIR
}

helpText () {
	echo ""
	echo "Runs test scripts for compOncDB."
	echo "Usage: ./test.sh {all/whitebox/blackbox/parse/db/fmt/vet}"
	echo ""
	echo "all		Runs all tests."
	echo "whitebox		Runs white box tests."
	echo "blackbox		Runs all black box tests (parse, upload, search, and update)."
	echo "parse		Runs parseRecords black box tests."
	echo "cancerrate	Runs cancer rate calculation black box tests."
	echo "db		Runs upload, search, update, and delete black box tests."
	echo "fmt		Runs go fmt on all source files."
	echo "vet		Runs go vet on all source files."
	echo "help		Prints help text."
}

if [ $# -eq 0 ]; then
	helpText
elif [ $1 = "all" ]; then
	getUser
	whiteBoxTests
	testParseRecords
	testDataBase
elif [ $1 = "whitebox" ]; then
	whiteBoxTests
elif [ $1 = "blackbox" ]; then
	getUser
	testParseRecords
	testDataBase
elif [ $1 = "parse" ]; then
	testParseRecords
elif [ $1 = "cancerrate" ]; then
	getUser
	testCancerRates
elif [ $1 = "db" ]; then
	getUser
	testDataBase
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
