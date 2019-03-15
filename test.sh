#!/bin/bash

##############################################################################
#	Manages black box tests on the comaprative oncology mysql database
#	Make sure a test database exists prior to running.
#
#	Required:	go 1.10+
#				mysql 14.14+
#
#	Usage:		./test.sh {install/white/parse/upload/search/update}
##############################################################################
USER="root"
WD=$(pwd)
PRSRC="$WD/src/parseRecords/*.go"
DBSRC="$WD/src/compOncDB/*.go"
DUSRC="$WD/src/dbupload/*.go"
DESRC="$WD/src/dbextract/*.go"
CDB="$WD/bin/compOncDB"
PR="$WD/bin/parseRecords"
CONFIG="$WD/bin/config.txt"

TESTDIR=$WD/test
TESTPR="$TESTDIR/parseRecords_test.go"
TESTDB="$TESTDIR/coDB_test.go"

TAXA="$TESTDIR/taxonomies.csv"
LIFEHIST="$TESTDIR/testLifeHistories.csv"
DENOM="$TESTDIR/testDenominators.csv"
DIAG="$TESTDIR/testDiagnosis.csv"
PATIENTS="$TESTDIR/testUpload.csv"

SERVICE="NWZP"
INPUT="$TESTDIR/testInput.csv"
DOUT="$TESTDIR/diagnoses.csv"
MOUT="$TESTDIR/merged.csv"
CASES="$TESTDIR/searchCases.csv"
UPDATE="$TESTDIR/testUpdate.csv"

whiteBoxTests () {
	echo ""
	echo "Running white box tests..."
	go test $PRSRC
	go test $DBSRC
	go test $DUSRC
	go test $DESRC
}

testParseRecords () {
	echo ""
	echo "Running black box tests on parseRecords..."
	$PR extract -s $SERVICE -i $INPUT -o $DOUT
	go test $TESTPR --run TestExtractDiagnosis  --args --expected=$DIAG --actual=$DOUT
	$PR merge -s $SERVICE -i $INPUT -t $TAXA -d $DIAG -o $MOUT
	go test $TESTPR --run TestMergeRecords --args --expected=$PATIENTS --actual=$MOUT
	# Delete test files
	rm $DOUT
	rm $MOUT
}

testUpload () {
	# Upload test data
	echo ""
	echo "Running black box tests on database upload..."
	$CDB test --config $CONFIG -u $USER -i $PATIENTS --taxonomy $TAXA --lifehistory $LIFEHIST --diagnosis $DIAG --denominators $DENOM -o "$TESTDIR/tables/"
	# Compare tables to expected
	go test $TESTDB --run TestDumpTables --args --indir="$TESTDIR/tables/"
}

testSearch () {
	# Test search output
	echo ""
	echo "Running black box tests on database search..."
	$CDB test --search --config $CONFIG -u $USER -i $CASES -o "$TESTDIR/searchResults/"
	go test $TESTDB --run TestSearches --args --indir="$TESTDIR/searchResults/"
}

testUpdates () {
	# Test search output
	echo ""
	echo "Running black box tests on database update..."
	$CDB test --update --config $CONFIG -u $USER -i $UPDATE -o "$TESTDIR/updateResults/"
	go test $TESTDB --run TestUpdates --args --indir="$TESTDIR/updateResults/"
}

testAll () {
	whiteBoxTests
	testParseRecords
	testUpload
	testSearch
	testUpdates
}

checkSource () {
	# Runs go fmt/vet on source files
	echo ""
	echo "Running go $1..."
	go $1 $PRSRC
	go $1 $DBSRC
	go $1 $DUSRC
	go $1 $DESRC
}

if [ $# -eq 2 ]; then
	# Set user to input value
	USER=$2
fi

if [ $# -eq 0 ]; then
	testAll
elif [ $1 = "all" ]; then
	testAll
elif [ $1 = "install" ]; then
	# Compile binaries and call test functions
	./install.sh
	testAll
elif [ $1 = "white" ]; then
	whiteBoxTests
elif [ $1 = "parse" ]; then
	testParseRecords
elif [ $1 = "upload" ]; then
	testUpload
elif [ $1 = "search" ]; then
	testSearch
elif [ $1 = "update" ]; then
	testUpdates
elif [ $1 = "fmt" ]; then
	checkSource $1
elif [ $1 = "vet" ]; then
	checkSource $1
elif [ $1 = "help" ]; then
	echo ""
	echo "Runs test scripts for compOncDB. Omit command line arguments to run all tests."
	echo "Usage: ./test.sh {install/white/parse/upload/search/update} username"
	echo ""
	echo "all		Runs all tests"
	echo "install		Installs binaries and runs all tests"
	echo "white		Runs white box tests"
	echo "parse		Runs parseRecords black box tests"
	echo "upload		Runs compOncDB upload black box tests"
	echo "search		Runs compOncDB search black box tests"
	echo "update		Runs compOncDB update black box tests"
	echo "fmt		Runs go fmt on all source files."
	echo "vet		Runs go vet on all source files."
	echo "username	MySQL username (root by default)."
fi
echo ""
