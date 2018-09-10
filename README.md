# Go-based MySQL database CRUD program  
This program is meant specifically for managing the comparative oncology database for the Maley lab at Arizona State University.  
It is currently being developed, so many features are not yet available.   

Copyright 2018 by Shawn Rupp

## Dependencies  
Go version 1.10 or higher  
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

### Uasage by Command  

#### Backup  
	./compOncDB backup

Backs up database to local machine. Must use root password. Output is written to current directory.  

#### New  
	./compOncDB new {-u username}  

	-u, --user="root"	MySQL username (default is root).  

Initializes new tables in new database. The database itself must be initialized manually.  
Make sure tableColumns.txt is in the bin/ directory.  

#### Upload  
	./compOncDB upload {-u username} --{type_from_list_below} infile

	-u, --user="root"	MySQL username (default is root).  
	--taxa				Load taxonomy tables from Kestrel output to update taxonomy table.  
	--common			Additionally extract common names from Kestrel output to update common name tables.  
	--lh				Upload life history info from merged life history table to the database.  
	--accounts			Extract account info from input file and update database.  
	--diagnosis			Extract diagnosis info from input file and update database.  
	--patient			Upload patient info from input table to database.  

	infile	Path to appropriate input file (Required).  

Uploads data from input files to appropriate tables. Only one flag may be given to indicate the type of 
input data, and therefore which tables must be updated. The only exception is for the --common flag which 
may be given in addition to the --taxa flag. This indicates that the Kestrel search terms are comman 
names and should be uploaded to the common names table. The input for --accounts, --diagnosis, and --patient
are all the same file which must in the format of uploadTemplate.csv.  

#### Update  
	./compOncDB update {-u username} {infile}

	-u, --user="root"	MySQL username (default is root).  
	--delete			Delete records from given table if column = value.  


	infile	Path to appropriate input file.  

Update or delete existing records from the database.  

#### Extract  
	./compOncDB extract {-u username} {--flags...} outfile

	-u, --user="root"	MySQL username (default is root).  
	-d, --dump="nil"	Name of table to dump (writes all data from table to output file).  
	--cancerRate		Calculates cancer rates for species with greater than min entries.  
	-m, --min=50		Minimum number of entries required for calculations (default = 50).  
	--necropsy			Extract only necropsy records (extracts all matches by default).  

	outfile				Name of output file (writes to stdout if not given).  

Extract data from the database and perform optional analyses.  
