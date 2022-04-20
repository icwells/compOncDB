[![Build Status](https://travis-ci.com/icwells/componcdb.svg?branch=master)](https://travis-ci.com/icwells/componcdb)
# Go-based MySQL database CRUD program  
This program is meant specifically for managing the comparative oncology database for the Maley lab at Arizona State University.  
It is currently being developed, so many features are not yet available.   

Copyright 2022 by Shawn Rupp

1. [Description](#Description)
2. [Installation](#Installation)  
3. [Usage](#Usage)  
4. [Commands](#Commands)  
5. [Search](#Search)  
6. [Cancer Rates](#Cancer Rates)  
7. [run.sh](#run.sh)  
8. [test.sh](#test.sh)

## Description  
componcdb is a program written to manage veterinary pathology data and identify cancer records using a MySQL database. 
The parse command can be used to extract diagnosis information from an input file. 
The program provides basic CRUD (create, read, update, delete) functionality, as well as specific analysis functions. 
These include calculating cancer rates per species, calculating summary statistics, and searching for specific records.  

## Installation  

### Dependencies  
[Go version 1.11 or higher](https://golang.org/doc/install)  
MySQL 14.14 or higher  
GNU Aspell 

### Installing Go and Setting Paths  
Go requires a GOPATH environment variable to set to install packages, an componcdb requires the GOBIN variable to be set as well.  
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
	git clone https://github.com/icwells/componcdb.git  

### Compiling scripts:
Any missing Go packages will be downloaded and installed when running install.sh.  

	cd componcdb/  
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

### Run Server Process  
Run the following in a terminal to launch the server program:

	./run.sh start/stop

The start command will kill an existing process and start a new one. The stop command will kill an existing process.  

## Usage  
Once compiled, the componcdb program can be used by giving it a base command and the appropriate flags.  
The program will prompt for a mysql password for the given username.  

Make sure the "comparativeOncology" database has been created in MySQL before running.  

### Overview  

	componcdb command {flags} {args...}  

	--help {command}	Show help {for given command}.  
	version			Prints version info and exits.  
	backup			Backs up database to local machine.  
	new			Initializes new tables in new database.  
	parse			Parse and organize records for upload to the comparative oncology database.  
	verify			Compares parse output with NLP model predictions. Provide parse records output and new output file with -i and -o.  
	upload			Upload data to the database.  
	update			Update or delete existing records from the database.  
	extract			Extract data from the database.  
	search			Search database for matches to queries.  
	leader			Calculate neoplasia prevalence leaderboards.  
	cancerrates		Calculate neoplasia prevalence for species.
	newuser			Adds new user to database. Must performed on the server using root password.  


### Commands  

#### Backup  
	componcdb backup -u username -o outfile  

	-o, --outfile		Path to output directory (Prints to current directory by default).  
	-u, --user		MySQL username.  

Backs up database to local machine. Must use root password. Output is written to current directory.  

#### New  
	componcdb new -u username  

	-u, --user	MySQL username.  

Initializes new tables in new database. The database itself must be initialized manually.  
Make sure tableColumns.txt is in the bin/ directory.  

#### Parse  
Uses a regular expression based approach to identify diagnosis information and merge the input and diagnosis 
information with taxonomy information (output of the [Kestrel](https://github.com/icwells/Kestrel) tool) 
to make a csv file ready to upload to the MySQL database.  

	componcdb parse -s service_name -t taxa_file -i infile -o outfile

	-d, --debug			Adds cancer and code column (if present) for hand checking.  
	-i, --infile			Path to input file (required).  
	-o, --outfile			Path to output file (required).  
	-s, --service			Database/service name (required).  
	-t, --taxa			Path to kestrel output (used with the merge command).  

Input files for parsing should have columns with the following names (in no particular order):

	ID			Unique patient ID.
	CommonName		Common species name of patient.
	ScientificName		Binomial scientific species name.
	Age			Patient age in months.
	Sex			Patient sex (male/female).
	Castrated		Whether patient was neutered/spayed.
	Location		Tissue tumor occured in.
	Type			Type of tumor (e.g. carcinoma).
	PrimaryTumor		Whether tumor was the original tumor (if multiple found).
	Metastasis		Whether tumor metastasized.
	Malignant		Whether tumor was malignant.
	Necropsy		Whether a necropsy was performed.
	Date			Date of diagnosis.
	Year			Year of diangosis.
	Comments		Additional diangosis info.
	Account			Account number or code.
	Client			Name of record submitter.

#### Verify  
	componcdb verify {--merge} -i infile -o outfile  

	-i, --infile	Path to input file (required).  
	--diagnosis	Verifies type and location diagnoses only.  
	--merge		Merges currated verification results with parse output. Give path to nlp output with -i and path to parse output with -o (it will be overwritten).  
	-o, --outfile	Path to output file (required).  
	--neoplasia	Verifies masspresent diagnosis only.

Calls NLP pipeline on parse output to flag diagnosis data that may not be accurate. If the --neoplasia flag is given, mass present and hyperplasia will be examined. If the --diangosis flag is given, only tumor type annd location will be examined. Otherwise, all four columns will be examined. Records will be printed to file if any inconsistency is found. These records can be manually currated by changing the column value (Masspresent, Hyperplasia, Type, or Location) in place. You can then rerun with the --merge flag, which will write the corrected values into the parse output file. You may then proceed with uploading the file.  

##### Train NLP Model  
Prior to verifying parse output, you must train the natural language processing model. To do this, first change into the nlpModel directory and pull training data from the website:

	cd Scripts/nlpModel/
	go run main -u {username} -o outfile

Then format the trianing data for use with the model:

	python nlpModel.py -i path_to_training_data

Lastly, call the script to trian the neoplasia and diagnosis models. Each step will take around 30 minutes.  

	python nlpModel.py
	python nlpModel.py --diagnosis


#### Upload  
	componcdb upload -u username --{type_from_list_below} -i infile

	--common		Additionally extract common names from Kestrel output to update common name tables.  
	--den			Uploads file to denominator table for databases where only cancer records were extracted.  
	-i infile		Path to appropriate input file (Required).  
	--lh			Upload life history info from merged life history table to the database.   
	--patient		Upload patient info from input table to database.
	--taxa			Load taxonomy tables from Kestrel output to update taxonomy table.    
	-u, --user		MySQL username.  

Uploads data from input files to appropriate tables. Only one flag may be given to indicate the type of 
input data, and therefore which tables must be updated. The only exception is for the --common flag which 
may be given in addition to the --taxa flag. This indicates that the Kestrel search terms are comman 
names and should be uploaded to the common names table. The input for --accounts, --diagnosis, and --patient
are all the same file which must in the format of uploadTemplate.csv.  

#### Update  
	componcdb update -u username {--flags...} {infile}

	-c, --column	Column to be updated with given value if --eval column == value.
	--clean		Remove extraneous records from the database.  
	--delete	Delete records from given table if column = value. 
	-e, --eval	Searches tables for matches (table is automatically determined) ('column operator value'; valid operators: != = <= >= > <; wrap statement in quotation marks).
	-i infile	Path to input file (see below for formatting).  
	-o	outdir	Backs up database to given directory before performing update.  
	--table		Perform operations on this table only.   
	-u, --user	MySQL username.  
	-v, --value	Value to write to column if --eval column == value (only supply one evaluation statement).

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
	componcdb extract -u username {--flags...} {-o outfile}

	--alltaxa			Summarizes life history table for all species (performs summary for species with records in patient table by default).  
	-d, --dump			Name of table to dump (writes all data from table to output file).  
	--dump_db			Extracts entire database into a gzipped tarball of csv files (specify output directory with -o).  
	-i infile			Path to input file (see below for formatting).  
	--lhsummary			Summarizes life history table.  
	-o outfile			Name of output file (writes to stdout if not given).  
	-r, --reference_taxonomy	Returns merged common and taxonomy tables.  
	--summarize			Compiles basic summary statistics of the database.  
	-u, --user			MySQL username.  

Extract data from the database.  


#### Search  
	componcdb search -u username {--flags...} {-o outfile}

	-e, --eval	 Searches tables for matches (table is automatically determined) ('column operator value'; valid operators: != = <= >= > < ^; wrap statement in quotation marks and seperate multiple statements with commas; '^' will return match if the column contains the value).  
	--infant	Include infant records in results (excluded by default).  
	-o outfile	Name of output file (writes to stdout if not given).  
	--taxonomies	Searches for taxonomy matches given column of common/scientific names in a file.  
	-u, --user	MySQL username.  

For searching most tables, the only valid operators for the eval flag are = (or ==), !=, or ^. For searching the Totals or Life_history tables, valid operations also include less than (or equal to) (</<=) and greater than (or equal to) (>/>=). Options given with -e should wrapped in single or double quotes to avoid errors.   

Taxonomy information can be extracted for target species in a given input file by specifying the "--taxonomies" flag. 
This will search for matches in the "-n" column of an input file (the first column by default). The species names can be either common or scientific names.  

#### Cancer Rates  
	componcdb cancerrates -u username {--flags...} {-o outfile}

	-e, --eval	 	Searches tables for matches (table is automatically determined) ('column operator value'; valid operators: != = <= >= > < ^; wrap statement in quotation marks and seperate multiple statements with commas; '^' will return match if the column contains the value).  
	--infant		Include infant records in results (excluded by default).  
	--keepall		Keep records without specified tissue when calculating by tissue.  
	--lifehistory		Append life history values to cancer rate data.  
	--location		Include tumor location summary for each species for given location.  
	-m, --min		Minimum number of entries required for calculations (default = 1).  
	--necropsy		2: Extract only necropsy records, 1: extract all records by default, 0: extract non-necropsy records.  
	--noavgage		Will not return average age columns in output file (AverageAge(months), AverageAgeNeoplasia(months)).  
	--nosexcol		Will not return male/female specific columns in output file (Male, MaleNeoplasia, MaleMalignant, Female, FemaleNeoplasia, FemaleMalignant).  
	--notaxacol		Will not return Kingdom-Genus columns in output file.  
	-o outfile		Name of output file (writes to stdout if not given).  
	--pathology		Additionally extract pathology records for target species.  
	--tissue		Include tumor tissue type summary for each species (supercedes location analysis).  
	-u, --user		MySQL username.  
	--wild			Return results for wild records only (returns non-wild only by default).  
	-z, --source		Zoo/institute records to calculate prevalence with; all: use all records, approved (default): used zoos approved for publication, aza: use only AZA member zoos, zoo: use only zoos.  

Returns the cancer rates by species for records matching given search criteria. The "--min" flag specifies the minimum number of species required to report cancer rates.  

### New User  
	componcdb newuser -u root {--admin} --username  

	--admin		Grant all privileges to user. Also allows remote MySQL access.  
	-u, --user	MySQL username. Must be "root" to create new users in MySQL.  
	--username	MySQL username for new user. Password will be be set to this name until it is updated.  

Creates new MySQL user. Must performed on the server using root password.  

## Bash Scripts  

### run.sh  
Runs hosting server for the comparative oncology database. Since the nginx server redirects to localHost, this is used for local testing and hosting on the server.  

	start	Kills running processes and starts new server on port 8080.  
	stop	Kills process running on port 8080.  
	help	Prints help text and exits.  

### test.sh  
Runs test scripts and functions for compOncDB.  
Usage: ./test.sh {all/fmt/vet/...}  

	all		Runs all tests.
	whitebox	Runs white box tests on all files in the src directory (except search directory, which requires mysql credentials).  
	blackbox	Runs all black box tests (parse, upload, search, cancerrate, and update).  
	parse		Runs parseRecords black box tests.  
	cancerrate	Runs cancer rate calculation black box tests (also runs upload test to ensure original test data is present).  
	necropsy	Runs necropsy filtering black box tests.  
	search		Runs white box tests on database search.  
	db		Runs upload, search, update, and delete black box tests.  
	fmt		Runs go fmt on all files in src and codb directories.  
	vet		Runs go vet on all files in src and codb directories.  
	help		Prints help text.  
