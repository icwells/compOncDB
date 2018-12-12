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

## Usage  
Once compiled, the compOncDB program can be used by giving it a base command and the appropriate flags.  
The program will prompt for a mysql password for the given username.  

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
	--eval="nil"		Searches tables for matches (table is automatically determined) (column operator value; valid operators: = <= >= > <).  

Update or delete existing records from the database.  

#### Extract  
	./compOncDB extract {-u username} {--flags...} {-o outfile}

	-u, --user="root"	MySQL username (default is root).  
	-d, --dump="nil"	Name of table to dump (writes all data from table to output file).  
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
	--eval="nil"			Searches tables for matches (table is automatically determined) (column operator value; valid operators: = <= >= > <).  
	--table="nil"			Return matching rows from this table only.  

	-o outfile				Name of output file (writes to stdout if not given).  

Searches database for matches to given criteria. The taxonomy search is given special consideration since it is the most common search type.  
For most tablesm the only valid operator for the eval flag is = (or ==). For searching the Totals or Life_history 
tables, valid operations also include less than (or equal to) (</<=) and greater than (or equal to) (>/>=).  

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
