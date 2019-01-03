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
WD=$(pwd)
PRSRC="$WD/src/parseRecords/*.go"
DBSRC="$WD/src/compOncDB/*.go"
DUSRC="$WD/src/dbupload/*.go"
CDB="$WD/bin/compOncDB"
PR="$WD/bin/parseRecords"
TABLES="$WD/bin/tableColumns.txt"

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
	echo "Running white box tests..."
	go test $PRSRC
	go test $DBSRC
	go test $DUSRC
}

testParseRecords () {
	echo ""
	echo "Running black box tests on parseRecords..."
	$PR extract -s $SERVICE -c $DICT -i $INPUT -o $DOUT
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
	$CDB test --tables $TABLES -i $PATIENTS --taxonomy $TAXA --lifehistory $LIFEHIST --diagnosis $DIAG --denominators $DENOM -o "$TESTDIR/tables/"
	# Compare tables to expected
	go test $TESTDB --run TestDumpTables --args --indir="$TESTDIR/tables/"
}

testSearch () {
	# Test search output
	echo ""
	echo "Running black box tests on database search..."
	$CDB test --search --tables $TABLES -i $CASES -o "$TESTDIR/searchResults/"
	go test $TESTDB --run TestSearches --args --indir="$TESTDIR/searchResults/"
}

testUpdates () {
	# Test search output
	echo ""
	echo "Running black box tests on database update..."
	$CDB test --update --tables $TABLES -i $UPDATE -o "$TESTDIR/updateResults/"
	go test $TESTDB --run TestUpdates --args --indir="$TESTDIR/updateResults/"
}

if [ $# -eq 0 ]; then
	whiteBoxTests
	testParseRecords
	testUpload
	testSearch
	testUpdates
elif [ $1 = "install" ]; then
	# Compile binaries and call test functions
	./install.sh
	whiteBoxTests
	testParseRecords
	testUpload
	testSearch
	testUpdates
elif [ $1 = "white" ]; then
	whiteBoxTests
elif [ $1 = "parse" ]; then
	testParseRecords
elif [ $1 = "upload" ]; then
	testUpload
elif [ $1 = "search" ]; then
	testSearch
elif [ $1 = "update" ]; then
	testUpdate
elif [ $1 = "help" ]; then
	echo ""
	echo "Runs test scripts for compOncDB. Omit command line arguments to run all tests."
	echo "Usage: ./test.sh {install/white/parse/upload/search/update}"
	echo "install		Installs binaries and runs all tests"
	echo "white		Runs white box tests"
	echo "parse		Runs parseRecords black box tests"
	echo "upload		Runs compOncDB upload black box tests"
	echo "search		Runs compOncDB search black box tests"
	echo "update		Runs compOncDB update black box tests"
fi
echo ""
