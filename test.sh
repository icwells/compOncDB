#!/bin/bash

##############################################################################
#	Manages black box tests on the comaprative oncology mysql database
#	Make sure a test database exists prior to running.
#
#	Required:	go 1.10+
#				mysql 14.14+
##############################################################################
TESTDIR=$(pwd)/test
CDB=bin/compOncDB
PR=bin/parseRecords
TESTPR="$TESTDIR/parseRecords_test.go"
TESTDB="$TESTDIR/coDB_test.go"

TAXA="$TESTDIR/taxonomies.csv"
LIFEHIST="$TESTDIR/testLifeHistories.csv"
DENOM=
DIAG="$TESTDIR/testDiagnosis.csv"
PATIENTS=

SERVICE="NWZP"
INPUT="$TESTDIR/testInput.csv"
DOUT="$TESTDIR/diagnoses.csv"
MOUT="$TESTDIR/merged.csv"

# Compile binaries
./install.sh

# Test parseRecords
$PR extract -s $SERVICE -i $INPUT -o $DOUT
#$PR merge -s $SERVICE -i $INPUT -o $MOUT
go test $TESTPR --args $TESTDIR
# Delete test files
#rm $DOUT
#rm $MOUT

# Upload test data
#$CDB test -i $PATIENTS --taxafile $TAXA --lifehistory $LIFEHIST --diagnosis $DIAG --denominators $DENOM

# Dump tables and compare to expected

# Check table sizes

# Test taxaSearch output

# Test columnSearch output
