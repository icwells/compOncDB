#!/bin/bash

##############################################################################
#	Manages black box tests on the comaprative oncology mysql database
#	Make sure a test database exists prior to running.
#
#	Required:	go 1.10+
#				mysql 14.14+
##############################################################################
DBSRC=src/compOncDB/*.go
PRSRC=src/parseRecords/*.go
CDB=bin/compOncDB
PR=bin/parseRecords

TESTDIR=$(pwd)/test
TESTPR="$TESTDIR/parseRecords_test.go"
TESTDB="$TESTDIR/coDB_test.go"

TAXA="$TESTDIR/taxonomies.csv"
LIFEHIST="$TESTDIR/testLifeHistories.csv"
DENOM=
DIAG="$TESTDIR/testDiagnosis.csv"
PATIENTS="$TESTDIR/testUpload.csv"

SERVICE="NWZP"
INPUT="$TESTDIR/testInput.csv"
DOUT="$TESTDIR/diagnoses.csv"
MOUT="$TESTDIR/merged.csv"

# Compile binaries
./install.sh

# White box tests
echo "Running white box tests..."
go test $DBSRC
go test $PRSRC

# Test parseRecords
echo ""
echo "Running black box tests on parseRecords..."
$PR extract -s $SERVICE -i $INPUT -o $DOUT
go test $TESTPR --run TestExtractDiagnosis  --args --expected=$DIAG --actual=$DOUT
$PR merge -s $SERVICE -i $INPUT -t $TAXA -d $DIAG -o $MOUT
go test $TESTPR --run TestMergeRecords --args --expected=$PATIENTS --actual=$MOUT
# Delete test files
rm $DOUT
#rm $MOUT

# Upload test data
#$CDB test -i $PATIENTS --taxafile $TAXA --lifehistory $LIFEHIST --diagnosis $DIAG --denominators $DENOM

# Dump tables and compare to expected

# Check table sizes

# Test taxaSearch output

# Test columnSearch output
