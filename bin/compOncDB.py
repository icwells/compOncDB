'''This script will direct commands for the comparative oncology database.'''

import argparse
import MySQLdb
import os
from datetime import datetime, date
from getpass import getpass
from dbIO import *
from version import version

HOST = "localhost"
DB = "comparativeOncology"

def backup():
	# Backup database to local Linux machine
	print(("\n\tBacking up {} database to loacl machine...").format(DB))
	f = open(os.devnull, 'w')
	sys.stdout = f
	try:
		p = Popen(split(("mysqldump -u root -p --result-file={}.{}.sql '{}'").format(DB, date.today(), DB)), stdout = f)
		# Wait for completion
		p.communicate()
		# Check for errors
		if(p.returncode != 0):
			raise
		print("\tBackup complete.\n")
	except:
		print("\tBackup failed.\n")

def connect(username):
	# Connects to mysql database
	password = getpass(prompt = "\tEnter MySQL password: ")
	try:
		# Connect to database
		db = MySQLdb.connect(HOST, username, password, DB)
		return db, password
	except:
		print("\n[Error] Cannot connect to database. Please check your username and password.\n")
		quit()

def printError(msg):
	# Prints error message and exits
	print(("\n\t[Error] {}. Exiting.\n").format(msg))
	quit()

def checkArgs(args):
	# Attempts to identify errors in the given arguments
	if not os.path.isfile("tableColumns.txt"):
		printError("Cannot find tableColumns.txt. Please be sure you are in the bin/.)
	if args.i and not os.path.isfile(args.i):
		printError(("Cannot find {}.").format(args.i))

def main():
	starttime = datetime.now()
	parser = argparse.ArgumentParser()
	parser.add_argument("-v", action = "store_true", 
help = "Prints version info and exits.")
	parser.add_argument("--backup", action = "store_true", help = "Backup database to local machine.")
	parser.add_argument("-u", default = "root", help = "MySQL username (default is 'root').")
	parser.add_argument("--new", action = "store_true", default = False,
 help = "Initializes new tables in new database (database must be made manually).")
	parser.add_argument("--dump", help = "Name of table to dump (writes all data from table to output file/stdout).")
	parser.add_argument("-i", help = "Path to input file. Data will be uploaded to the database.")
	parser.add_argument("-o", help = "Path to output file (if extracting data).")
	args = parser.parse_args()
	if args.v:
		version()
	checkArgs(args)
	db, pw = connect(args.u)
	cursor = db.cursor()
	if args.new == True:
		# Create new table
		newTable(cursordb, args.t)
	elif args.dump:
		# Extract table from db
		dumpTable(cursor, args.dump, args.o)
	if args.i:
		convertFF(args.i, db, args.t, columns, ids)
	db.close()
	print(("\tTotal runtime: {}\n").format(datetime.now()-starttime))

if __name__ == "__main__":
	main()
