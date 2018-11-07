#!/bin/bash

##############################################################################
#	Manages black box tests on the comaprative oncology mysql database
#	Make sure a test database exists prior to running.
#
#	Required:	go 1.10+
#				mysql 14.14+
##############################################################################

CDB=bin/compOncDB
PR=bin/parseRecords

SERVICE="NWZP"
TAXA=test/taxonomies.csv
INPUT=test/testInput.csv


# Compile binaries
./install.sh

# Test extractDiagnosis
$PR extract -s $SERVICE -i $INPUT -o  

# Test merger

# Upload test data
$CDB new --test
$CDB upload --test --common --taxa -i $TAXA
$CDB upload --test --taxa -i $TAXA
$CDB upload --test --lh -i 
$CDB upload --test --den -i 
$CDB upload --test --patient -i 

# Check table sizes

# Dump tables and compare to expected

# Test taxaSearch output

# Test columnSearch output
