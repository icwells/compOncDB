# Go-based MySQL database CRUD program  
This program is meant specifically for managing the comparative oncology database for the Maley lab at Arizona State University.  
It is currently being developed, so many features are not yet available.   

Copyright 2018 by Shawn Rupp

## Dependencies  
Go version 1.11 or higher  
MySQL 14.14 or higher  

## Installation  
Download the repository:  

	git clone https://github.com/icwells/compOncDB.git  

### Compiling scripts:
Any missing Go packages will be downloaded and installed when running install.sh.  

	cd compOncDB/  
	./install.sh  

### Testing the Installation  
Create a database named "testDataBase" in MySQL (CREATE DATABASE testDataBase;) and run the following:

	./test.sh

You will be prompted the enter your MySQL password twice. All of the output from the test scripts should begin with "ok".  

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

### Usage by Command  

#### Backup  
	./compOncDB backup

Backs up database to local machine. Must use root password. Output is written to current directory.  

#### New  
	./compOncDB new {-u username}  

	-u, --user="root"	MySQL username (default is root).  

Initializes new tables in new database. The database itself must be initialized manually.  
Make sure tableColumns.txt is in the bin/ directory.  

#### Upload  
	./compOncDB upload {-u username} --{type_from_list_below} -i infile

	-u, --user="root"	MySQL username (default is root).  
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
	--count				Recount species totals and update the Totals table.  
	--delete			Delete records from given table if column = value (must be root).  
	-c, --column="nil"	Column to be updated with given value if --eval column == value.
	-v, --value="nil"	Value to write to column if --eval column == value.
	-e, --eval="nil"	Searches tables for matches (table is automatically determined) (column operator value; valid operators: = <= >= > <).  
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
	-l, --level="Species"	Taxonomic level of taxon (or entries in taxon file)(default = Species).  
	-t, --taxa="nil"		Name of taxonomic unit to extract data for or path to file with single column of units.  
	--common				Indicates that common species name was given for taxa.  
	--count					Returns count of target records instead of printing entire records.  
	-e, --eval="nil"		Searches tables for matches (table is automatically determined) (column operator value; valid operators: = <= >= > <).  
	--table="nil"			Return matching rows from this table only.  

	-o outfile				Name of output file (writes to stdout if not given).  

Searches database for matches to given criteria. The taxonomy search is given special consideration since it is the most common search type.  
For most tables, the only valid operator for the eval flag is = (or ==). For searching the Totals or Life_history 
tables, valid operations also include less than (or equal to) (</<=) and greater than (or equal to) (>/>=).  

#### Test  
	./compOncDB test {paths_to_input_files}/ {--search --eval operation -o path_to_output file}

	-u, --user="root"				MySQL username (default is root).
	-i, --infile="nil"				Path to input file (if using).
	-o, --outfile="nil"				Name of output file (writes to stdout if not given).
	--tables=TABLES					Path tableColumns.txt file.
	--taxonomy=TAXONOMY				Path to taxonomy file.
	--diagnosis=DIAGNOSIS			Path to extracted diganoses file.
	--lifehistory=LIFEHISTORY		Path to life history data.
	--denominators=DENOMINATORS		Path to file conataining non-cancer totals.
	--search						Search for matches using above commands.

Runs tests using "testDataBase" (make sure it has been created in MySQL first).

## parseRecords  
The parseRecords utility has been provided to help with parsing input files before uploading them to the database.  
The extract functionality will use a regular expression based approach to identify diagnosis information.  
The merge command will merge the output of the diagnosis command and taxonomy information (output of the [Kestrel](https://github.com/icwells/Kestrel) tool) 
to make a csv file ready to upload to the MySQL database.  

### Usage  

	./parseRecords command -s service_name -i infile -o outfile

	extract							Extract diagnosis data from infile.  
	merge							Merges taxonomy and diagnosis info with infile.

	--help							Show help.  
	-i, --infile=INFILE				Path to input file (required).  
	-o, --outfile=OUTFILE			Path to output file (required).  
	-s, --service=SERVICE			Database/service name (required).  
	-d, --dict="cancerdict.tsv"		Path to dictionary of cancer terms (used with the extract command).  
	-t, --taxa="nil"				Path to kestrel output (used with the merge command).  
	-d, --diagnoses="nil"			Path to diagnosis data (used with the merge command).  
