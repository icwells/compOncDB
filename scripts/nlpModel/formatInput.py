'''Preformats input file for nlpModeler.'''

import re
from unixpath import readFile

class Formatter():

	def __init__(self, infile, outfile, encodingfile):
		self.count = 1
		self.encodingfile = encodingfile
		self.header = {}
		self.infile = infile
		self.outfile = outfile
		self.types = {"NA-NA": 0}
		self.__formatFile__()
		self.__writeDicts__()

	def __writeDicts__(self):
		# Writes locations and types to file
		with open(self.encodingfile, "w") as out:
			for k in self.types.keys():
				out.write("{},{}\n".format(k, self.types[k]))

	def __formatLine__(self, line):
		# Replaces punctuation, splits compound locations and types, and encodes paired locations and types as integers
		rows = []
		if line[self.header["Comments"]] == "NA" or line[self.header["Comments"]] == "n/a. n/a.":
			if line[self.header["Masspresent"]] == "1":
				# Skip records where diagnosis info is not in comments
				return rows
		line[self.header["Comments"]] = re.sub(r"[^\w\s]", "", line[self.header["Comments"]])
		# Split compound locations and types
		loc = line[self.header["Location"]].split(";")
		for idx, i in enumerate(line[self.header["Type"]].split(";")):
			# Combine location and type as one key
			t = "-".join([loc[idx], i])
			if t not in self.types.keys():
				self.types[t] = self.count
				self.count += 1
			row = line[:self.header["Type"]]
			row.append(str(self.types[t]))
			rows.append(row)
		return rows

	def __formatFile__(self):
		# Reads input file, formats values, and writes to output
		first = True
		print("\n\tFormatting input file...")
		with open(self.outfile, "w") as out:
			for line in readFile(self.infile, header = True, d = ","):
				if not first:
					for i in self.__formatLine__(line):
						out.write(",".join(i) + "\n")
				else:
					self.header = line
					row = []
					for i in range(len(self.header)):
						row.append(-1)
					for k in self.header.keys():
						row[self.header[k]] = k
					# Omit Location column
					out.write(",".join(row[:-1]) + "\n")
					first = False
