'''Preformats input file for nlpModeler.'''

import re
from unixpath import readFile

class Formatter():

	def __init__(self, infile, outfile, encodingfile):
		self.encodingfile = encodingfile
		self.header = {}
		self.infile = infile
		self.lcount = 1
		self.locations = {"NA": 0}
		self.outfile = outfile
		self.tcount = 1
		self.types = {"NA": 0}
		self.__formatFile__()
		self.__writeDicts__()

	def __writeDicts__(self):
		# Writes locations and types to file
		with open(self.encodingfile, "w") as out:
			for k in self.locations.keys():
				out.write("Location,{},{}\n".format(k, self.locations[k]))
			for k in self.types.keys():
				out.write("Type,{},{}\n".format(k, self.types[k]))

	def __formatLine__(self, line):
		# Replaces punctuation, splits compound locations and types, and encodes locations and types with integers
		rows = []
		line[self.header["Comments"]] = re.sub(r"[^\w\s]", "", line[self.header["Comments"]])
		# Split compound locations and types
		loc = line[self.header["Location"]].split(";")
		for idx, i in enumerate(line[self.header["Type"]].split(";")):
			l = loc[idx]
			if l not in self.locations.keys():
				self.locations[l] = self.lcount
				self.lcount += 1
			if i not in self.types.keys():
				self.types[i] = self.tcount
				self.tcount += 1
			row = line[:self.header["Type"]]
			row.append(str(self.types[i]))
			row.append(str(self.locations[l]))
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
					out.write(",".join(row) + "\n")
					first = False
