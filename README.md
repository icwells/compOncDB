[![Build Status](https://travis-ci.com/icwells/compOncDB.svg?branch=master)](https://travis-ci.com/icwells/compOncDB)
# Go-based MySQL database CRUD program  
This program is meant specifically for managing the comparative oncology database for the Maley lab at Arizona State University.  
It is currently being developed, so many features are not yet available.   

Copyright 2019 by Shawn Rupp

1. [Description](#Description)
2. [Installation](#Installation)  
3. [Usage](#Usage)  
4. [Commands](#Commands)  
5. [parseRecords](#parseRecords)

## Description  
compOncDB is a program written to manage veterinary pathology data and identify cancer records using  a MySQL database. The accompanying parseRecords tool can be used to extract diagnosis information from an input file. It can then be used to merge the source data, diagnosis output, and taxonomy information into a format ready for upload to the MySQL database. The program provides basic CRUD (create, read, update, delete) functionality, as well as specific analysis functions. These include calculating cancer rates per species, calculating summary statistics, and searching for specific records.  

## Installation  

### Dependencies  
Go version 1.11 or higher  
MySQL 14.14 or higher  

### Download  
Download the repository:  

	git clone https://github.com/icwells/compOncDB.git  

### Compiling scripts:
Any missing Go packages will be downloaded and installed when running install.sh.  

	cd compOncDB/  
	./install.sh  

### Config File  
The config file is located in the bin directory by default (although a different location can be specified with the --config flag). It contains 
basic connection information which will be used for all connections. This includes the host (leave blank for local host), database name, test database name, 
and the path to tableColumns.txt (also located in the bin by default). The host is the only field that may need changing, depending on whether you are using 
a local or a remote connection.  

Be sure to change the name of "example_config.txt" to "config.txt" (to prevent git from overwritting it).

### Testing the Installation  
Replace "mysql_username" with your username and run the following in a terminal:

	./test.sh all mysql_username

You will be prompted the enter your MySQL password several times. All of the output from the test scripts should begin with "ok".  

## Usage  
Once compiled, the compOncDB program can be used by giving it a base command and the appropriate flags.  
The program will prompt for a mysql password for the given username.  

Make sure the "comparativeOncology" database has been created in MySQL before running.  

### Overview  

	./compOncDB command {flags} {args...}  

	--help {command}	Show help {for given command}.  
	version			Prints version info and exits.  
	backup			Backs up database to local machine.  
	new				Initializes new tables in new database.
	upload			Upload data to the database.  
	update			Update or delete existing records from the database.  
	extract			Extract data from the database and perform optional analyses.  
	search			Searches database for matches to given term.  
	test			Tests database functionality using testDataBase instead of comaprative oncology.  

### Commands  

#### Backup  
	./compOncDB backup

Backs up database to local machine. Must use root password. Output is written to current directory.  

#### New  
	./compOncDB new {-u username}  

	-u, --user="root"	MySQL username (default is root).  
	--config="config.txt"  Path to config.txt (Default is in bin directory).  

Initializes new tables in new database. The database itself must be initialized manually.  
Make sure tableColumns.txt is in the bin/ directory.  

#### Upload  
	./compOncDB upload {-u username} --{type_from_list_below} -i infile

	-u, --user="root"	MySQL username (default is root).  
	--config="config.txt"  Path to config.txt (Default is in bin directory).  
	--taxa				Load taxonomy tables from Kestrel output to update taxonomy table.  
	--common			Additionally extract common names from Kestrel output to update common name tables.  
	--lh				Upload life history info from merged life history table to the database.   
	--den				Uploads file to denominator table for databases where only cancer records were extracted.  
	--patient			Upload patient info from input table to database.  


	-i infile			Path to appropriate input file (Required).  

Uploads data from input files to appropriate tables. Only one flag may be given to indicate the type of 
input data, and therefore which tables must be updated. The only exception is for the --common flag which 
may be given in addition to the --taxa flag. This indicates that the Kestrel search terms are comman 
names and should be uploaded to the common names table. The input for --accounts, --diagnosis, and --patient
are all the same file which must in the format of uploadTemplate.csv.  

#### Update  
	./compOncDB update {-u username} {infile}

	-u, --user="root"	MySQL username (default is root).  
	--config="config.txt"  Path to config.txt (Default is in bin directory).  
	--count				Recount species totals and update the Totals table.  
	--delete			Delete records from given table if column = value. 
	--table="nil"		Perform operations on this table only.   
	-c, --column="nil"	Column to be updated with given value if --eval column == value.
	-v, --value="nil"	Value to write to column if --eval column == value (only supply one evaluation statement).
	-e, --eval="nil"	 Searches tables for matches (table is automatically determined) ('column operator value'; valid operators: 
!= = <= >= > <; wrap statement in quotation marks).  

	-i infile			Path to input file (see below for formatting).  

Update or delete existing records from the database. Command line updates (given with the -c, -v, and -e flags) will only perform a single 
update operation; however, multiple updates can be run at once using an input file.  

If an input file is given (with the -i flag), it must be a csv/tsv/txt file where columns in the header are also columns in the database. 
The first column will be used to identify records to update (they will usually be ID numbers, but don't have to be). Values in this column 
will not be changed in the database. The values in the remaining columns will be used to update records with a matching identifier. Each value 
will be used to update the column with the same name as apears in the header. The tables will be automatically determined.  

For example:  

	ID	Age	Masspresent
	1	20	0

indicates that the record with an ID of "1" should have its Age (in the Patient table) changed to "20" and its Masspresent code (in the Diagnosis table) 
changed to "0". Since ID is a unique identifier this will only change one record, but if a taxonomic level were given, for example, all taxonomies with 
a matching taxonomic level would be updated.  

#### Extract  
	./compOncDB extract {-u username} {--flags...} {-o outfile}

	-u, --user="root"	MySQL username (default is root).  
	--config="config.txt"  Path to config.txt (Default is in bin directory).  
	-d, --dump="nil"	Name of table to dump (writes all data from table to output file).  
	--summarize			Compiles basic summary statistics of the database.  
	--cancerRate		Calculates cancer rates for species with greater than min entries.  
	-m, --min=50		Minimum number of entries required for calculations (default = 50).  
	--necropsy			Extract only necropsy records (extracts all matches by default).  

	-o outfile			Name of output file (writes to stdout if not given).  

Extract data from the database and perform optional analyses.  

#### Search  
	./compOncDB search {-u username} {--column/level --value/species ...} {-o outfile}

	-u, --user="root"		MySQL username (default is root).   
	--config="config.txt"  Path to config.txt (Default is in bin directory).  
	-l, --level="Species"	Taxonomic level of taxon (or entries in taxon file)(default = Species).  
	-t, --taxa="nil"		Name of taxonomic unit to extract data for or path to file with single column of units.  
	--common				Indicates that common species name was given for taxa.  
	--count					Returns count of target records instead of printing entire records.  
	-e, --eval="nil"	 Searches tables for matches (table is automatically determined) ('column operator value'; valid operators: 
!= = <= >= > <; wrap statement in quotation marks and seperate multiple statements with commas).  
	--table="nil"		Perform operations on this table only.  
	--infant			Include infant records in results (excluded by default).  
	--taxonomies		Searches for taxonomy matches given column of common/scientific names in a file.  
	-n, --names=0		Column of input file containing scientific/common species names to search.  

	-o outfile				Name of output file (writes to stdout if not given).  

Searches database for matches to given criteria. The taxonomy search is given special consideration since it is the most common search type.  
For most tables, the only valid operator for the eval flag is = (or ==). For searching the Totals or Life_history tables, valid operations also 
include less than (or equal to) (</<=) and greater than (or equal to) (>/>=). Options given with -e should wrapped in single or double quotes to avoid errors.  

#### Test  
	./compOncDB test {paths_to_input_files}/ {--search --eval operation -o path_to_output file}  

	-u, --user="root"				MySQL username (default is root).  
	--config="config.txt"  Path to config.txt (Default is in bin directory).  
	-i, --infile="nil"				Path to input file (if using).  
	-o, --outfile="nil"				Name of output file (writes to stdout if not given).  
	--taxonomy=TAXONOMY				Path to taxonomy file.  
	--diagnosis=DIAGNOSIS			Path to extracted diganoses file.  
	--lifehistory=LIFEHISTORY		Path to life history data.  
	--denominators=DENOMINATORS		Path to file conataining non-cancer totals.  
	--search						Search for matches using above commands.  

Runs tests using "testDataBase" (make sure it has been created in MySQL first).  

## parseRecords  
The parseRecords utility has been provided to help with parsing input files before uploading them to the database.  
It will use a regular expression based approach to identify diagnosis information and merge the input and diagnosis 
information with taxonomy information (output of the [Kestrel](https://github.com/icwells/Kestrel) tool) 
to make a csv file ready to upload to the MySQL database.  
Previous versions performed this in two steps, but as of v0.3 they have been merged into one to steamline the process.  


### Usage  

	./parseRecords command -s service_name -i infile -o outfile

	--help							Show help.  
	-i, --infile=INFILE				Path to input file (required).  
	-o, --outfile=OUTFILE			Path to output file (required).  
	-s, --service=SERVICE			Database/service name (required).  
	-t, --taxa="nil"				Path to kestrel output (used with the merge command).  
