'''These functions will upload and extract data from the ASUviral database.'''

from sys import stdout
import MySQLdb
import os
import re

def updateDB(db, table, columns, values):
	# Updates existing table
	cursor = db.cursor()
	# Conctruct SQL statement
	sql = ("INSERT INTO " + table + "(" + columns 
			+ ") VALUES(" + values + ");")
	try:
		# Insert new row
		cursor.execute(sql)
		db.commit()
		return 1
	except:
		db.rollback()
		err = values.split(",")[0].replace('"', '')
		if "ProteinID" in columns:
			err += ": " + values.split(",")[3].replace('"', '')
		print(("\tThere was an error uploading {}").format(err))
		return 0

def buildValues(gene, columns):
	# Organizes input data
	values = ""
	cols = columns.split(",")
	for i in cols:
		try:
			values += '"' + gene[i] + '",'
		except KeyError:
			# Append dash for missing info
			values += '"-",'
	return values[:-1]

def getIDs(db, table, token = "ID", new = False):
	# Retrieves list of ids numbers from database
	# Add to set to automatically skip repeats, return as list
	ids = set()
	cursor = db.cursor()
	if new == True:
		sql = ("SELECT {} FROM {} WHERE Date >= CURDATE();"
				).format(token, table)
	else:
		sql = ("SELECT {} FROM {};").format(token, table)
	cursor.execute(sql)
	result = cursor.fetchall()
	for i in result:	
		ids.add(i[0])
	return list(ids)

def getColumns(types = False):
	# Build dict of column statements
	columns = {}
	table = ""
	with open("tableColumns.txt") as col:
		for line in col:
			if line[0] == "#":
				# Get table names
				table = line[1:].strip()
			elif line.strip():
				# Get columns for given table
				if types == True:
					col = line.strip()
				else:
					col = line.split()[0]
				if table not in columns.keys():
					columns[table] = col
				else:
					columns[table] += "," + col
	return columns

def newTables(cursor, table):
	# Initializes new table
	print("\n\tInitializing new tables...")
	columns = getColumns(types = True)
	for i in columns.keys():
		cmd = ("CREATE TABLE {}({});").format(i, columns[i])
		try:
			cursor.execute(cmd)
			db.commit()
		except:
			db.rollback()
			print(("\t[Error] Creating table {}. Exiting.\n").format(i))
			quit()

#-----------------------------------------------------------------------------

def writeToCSV(outfile, header, results):
	# Wites table to file
	with open(outfile, "w") as out:
		out.write(header + "\n")
		for i in results:
			out.write((",").join(i) + "\n")

def dumpTable(cursor, table, outfile = None):
	# Dumps contents of table to outfile (if given)
	print(("\n\tExtracting from {}...").format(table))
	sql = ("SELECT * FROM {};").format(table)
	cursor.execute(sql)
	results = cursor.fetchall()
	if outfile:
		columns = getColumns()
		writeToCSV(outfile, columns[table], results)
	else:
		for i in results:
			print(i)

def extractRow(cursor, table, name, column):
	# Extract data from table
	sql = ('SELECT * FROM {} WHERE {} = "{}";').format(table, column, name)
	# Execute the SQL command
	cursor.execute(sql)
	# Fetch row from table
	results = cursor.fetchall()
	return results

def extractDNA(db, outdir):
	# Extracts dna sequences and accessions
	cursor = db.cursor()
	cdef str outfile
	cdef str sql
	acc = getBacAcc(cursor)
	print("\tExtracting nucleotide sequences in fasta format...")
	outfile = outdir + "viralRefSeq.fna"
	sql = 'SELECT Accession, DNA FROM Annotations;'
	# Execute the SQL command
	cursor.execute(sql)
	# Fetch row from table
	results = cursor.fetchall()
	with open(outfile, "w") as fasta:
		for row in results:
			try:
				if len(row[1]) > 2 and row[0] not in acc:
					# Skip entries with missing data
					fasta.write((">{}\n{}\n").format(row[0], row[1]))
					acc.append(row[0])
			except IndexError:
				pass
