[![Build Status](https://travis-ci.com/icwells/compOncDB.svg?branch=master)](https://travis-ci.com/icwells/compOncDB)
# Go-based MySQL database CRUD program  
This program is meant specifically for managing the comparative oncology database for the Maley lab at Arizona State University.  
It is currently being developed, so many features are not yet available.   

Copyright 2019 by Shawn Rupp

1. [Description](#Description)
2. [Installation](#Installation)  
3. [Usage](#Usage)  
4. [Commands](#Commands)  
5. [Extracting Data](#Extract)

## Description  
compOncDB is a program written to manage veterinary pathology data and identify cancer records using a MySQL database. 
The parse command can be used to extract diagnosis information from an input file. 
The program provides basic CRUD (create, read, update, delete) functionality, as well as specific analysis functions. 
These include calculating cancer rates per species, calculating summary statistics, and searching for specific records.  

## Installation  

### Dependencies  
[Go version 1.11 or higher](https://golang.org/doc/install)  
MySQL 14.14 or higher  
GNU Aspell 

### Installing Go and Setting Paths  
Go requires a GOPATH environment variable to set to install packages, an compOncDB requires the GOBIN variable to be set as well.  
Follow the directions [here](https://github.com/golang/go/wiki/SettingGOPATH) to set your GOPATH. Before you close your .bashrc or 
similar file, add the following lines after you deifne you GOPATH:  

	export GOBIN=$GOPATH/bin  
	export PATH=$PATH:$GOBIN  

### Installing GNU Aspell
parseRecords uses GNU Aspell to assist in resolving spelling errors in account and submitter names.  

	sudo apt-get install aspell libaspell-dev  

### Download  
Download the repository into correct Go src directory (required for package imports):  

	cd $GOPATH/src
	mkdir -p github.com/icwells/
	cd github.com/icwells/
	git clone https://github.com/icwells/compOncDB.git  

### Compiling scripts:
Any missing Go packages will be downloaded and installed when running install.sh.  

	cd compOncDB/  
	./install.sh  

### Config File  
The config file is located in the utils directory. It contains basic connection information which will be used for all connections. This includes the host (leave blank for local host), 
database name, test database name, and the path to tableColumns.txt (also located in the bin by default). The host is the only field that may need changing, depending on whether you are using 
a local or a remote connection.  

Be sure to change the name of "example_config.txt" to "config.txt" (to prevent git from overwriting it).

### Testing the Installation  
Run the following in a terminal:

	./test.sh all

You will be prompted the enter your MySQL user name and password at the beginning. All of the output from the test scripts should begin with "ok".  

## Usage  
Once compiled, the compOncDB program can be used by giving it a base command and the appropriate flags.  
The program will prompt for a mysql password for the given username.  

Make sure the "comparativeOncology" database has been created in MySQL before running.  

### Overview  

	compOncDB command {flags} {args...}  

	--help {command}	Show help {for given command}.  
	version			Prints version info and exits.  
	backup			Backs up database to local machine.  
	new			Initializes new tables in new database.
	parse			Parse and organize records for upload to the comparative oncology database.  
	upload			Upload data to the database.  
	update			Update or delete existing records from the database.  
	extract			Extract data from the database.
	search			Search database for matches to queries.
	cancerrates		Calculate neoplasia prevalence for species.

### Commands  

#### Backup  
	compOncDB backup

Backs up database to local machine. Must use root password. Output is written to current directory.  

#### New  
	compOncDB new {-u username}  

	-u, --user="root"	MySQL username (default is root).  

Initializes new tables in new database. The database itself must be initialized manually.  
Make sure tableColumns.txt is in the bin/ directory.  

#### Parse  
Uses a regular expression based approach to identify diagnosis information and merge the input and diagnosis 
information with taxonomy information (output of the [Kestrel](https://github.com/icwells/Kestrel) tool) 
to make a csv file ready to upload to the MySQL database.  

	compOncDB parse -s service_name -t taxa_file -i infile -o outfile

	--help					Show help.  
	-i, --infile=INFILE			Path to input file (required).  
	-o, --outfile=OUTFILE			Path to output file (required).  
	-s, --service=SERVICE			Database/service name (required).  
	-t, --taxa="nil"			Path to kestrel output (used with the merge command).  
	-d, --debug				Adds cancer and code column (if present) for hand checking.  

Input files for parsing should have columns with the following names (in no particular order):

	ID				Unique patient ID.
	CommonName		Common species name of patient.
	ScientificName	Binomial scientific species name.
	Age				Patient age in months.
	Sex				Patient sex (male/female).
	Castrated		Whether patient was neutered/spayed.
	Location		Tissue tumor occured in.
	Type			Type of tumor (e.g. carcinoma).
	PrimaryTumor	Whether tumor was the original tumor (if multiple found).
	Metastasis		Whether tumor metastasized.
	Malignant		Whether tumor was malignant.
	Necropsy		Whether a necropsy was performed.
	Date			Date of diagnosis.
	Year			Year of diangosis.
	Comments		Additional diangosis info.
	Account			Account number or code.
	Client			Name of record submitter.

#### Upload  
	compOncDB upload {-u username} --{type_from_list_below} -i infile

	-u, --user="root"		MySQL username (default is root).  
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
	compOncDB update {-u username} {infile}

	-u, --user="root"	MySQL username (default is root).  
	--count			Recount species totals and update the Totals table.  
	--delete		Delete records from given table if column = value. 
	--table="nil"		Perform operations on this table only.   
	-c, --column="nil"	Column to be updated with given value if --eval column == value.
	-v, --value="nil"	Value to write to column if --eval column == value (only supply one evaluation statement).
	-e, --eval="nil"	Searches tables for matches (table is automatically determined) ('column operator value'; valid operators: != = <= >= > <; wrap statement in quotation marks).  

	-i infile		Path to input file (see below for formatting).  

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
	compOncDB extract {-u username} {--flags...} {-o outfile}

	-u, --user="root"	MySQL username (default is root).  
	--table="nil"		Perform operations on this table only.  
	-d, --dump="nil"	Name of table to dump (writes all data from table to output file).  
	--lhsummary			Summarizes life history table.  
	-r, --reference_taxonomy	Returns merged common and taxonomy tables.  
	--summarize		Compiles basic summary statistics of the database.  

	-i infile		Path to input file (see below for formatting).  
	-o outfile		Name of output file (writes to stdout if not given).  

Extract data from the database.  


#### Search
	compOncDB search {-u username} {--flags...} {-o outfile}

	-u, --user="root"	MySQL username (default is root).  
	-e, --eval="nil"	 Searches tables for matches (table is automatically determined) ('column operator value'; valid operators: != = <= >= > < ^; wrap statement in quotation marks and seperate multiple statements with commas; '^' will return match if the column contains the value).  
	--infant		Include infant records in results (excluded by default).  
  	-n, --names=0		Column of input file containing scientific/common species names to search.    
	--table="nil"		Perform operations on this table only.  
	--taxonomies		Searches for taxonomy matches given column of common/scientific names in a file.  

	-i infile		Path to input file (see below for formatting).  
	-o outfile		Name of output file (writes to stdout if not given).  

For searching most tables, the only valid operators for the eval flag are = (or ==), !=, or ^. For searching the Totals or Life_history tables, valid operations also include less than (or equal to) (</<=) and greater than (or equal to) (>/>=). Options given with -e should wrapped in single or double quotes to avoid errors.  

If an input file is specified, it should follow a similar format to the update file, with the exception that every column will be used as a search criterion. Each line of the file represents a seperate search, while each column of the row will be added the same search and will be cumulative. The program will search for matches where the value in the file equals the value in the given column (empty cells will be skipped). Therefore:

	taxa_id	Age	Masspresent
		20	
	2		0
	3	10	1

will search for any match where Age == 20, records where taxa_id == 2 and Masspresent == 0, and records where taxa_id == 3, Age == 10, and Masspresent == 1.  

Taxonomy information can be extracted for target species in a given input file by specifying the "--taxonomies" flag. 
This will search for matches in the "-n" column of an input file (the first column by default). The species names can be either common or scientific names.  

#### Cancer Rates
	compOncDB cancerrates {-u username} {--flags...} {-o outfile}

	-u, --user="root"	MySQL username (default is root).  
	-e, --eval="nil"	 Searches tables for matches (table is automatically determined) ('column operator value'; valid operators: != = <= >= > < ^; wrap statement in quotation marks and seperate multiple statements with commas; '^' will return match if the column contains the value).  
	--approved		Calculate neoplasia prevalence using only records from approved zoos.  
	--infant		Include infant records in results (excluded by default).  
	--lifehistory	Append life history values to cancer rate data.  
	--location=""		Include tumor location summary for each species for given location.  
	-m, --min=50		Minimum number of entries required for calculations (default = 50).  
	--necropsy		Extract only necropsy records (extracts all matches by default).  

	-i infile		Path to input file (see below for formatting).  
	-o outfile		Name of output file (writes to stdout if not given).  

Returns the cancer rates by species for records matching given search criteria. The "--min" flag specifies the minimum number of species required to report cancer rates.  


