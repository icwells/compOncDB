#!/bin/bash

##############################################################################
#	Manages black box tests on the comaprative oncology mysql database
#	Make sure a test database exists prior to running.
#
#	Required:	go 1.10+
#				mysql 14.14+
##############################################################################
WD=$(pwd)
DBSRC="$WD/src/compOncDB/*.go"
PRSRC="$WD/src/parseRecords/*.go"
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

whiteBoxTests () {
	echo "Running white box tests..."
	go test $DBSRC
	go test $PRSRC
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

testDataBase () {
	# Upload test data
	echo ""
	echo "Running black box tests on database upload..."
	$CDB test --tables $TABLES -i $PATIENTS --taxonomy $TAXA --lifehistory $LIFEHIST --diagnosis $DIAG --denominators $DENOM -o "$TESTDIR/tables/"
	# Compare tables to expected
	go test $TESTDB --run TestDumpTables --args --indir="$TESTDIR/tables/"
	# Test search output
	$CDB test --search --tables $TABLES -i $CASES -o "$TESTDIR/searchResults"
	go test $TESTDB --run TestSearches --args --indir="$TESTDIR/searchResults"
}

# Compile binaries and call test functions
./install.sh

#whiteBoxTests
#testParseRecords
testDataBase
echo ""
