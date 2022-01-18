'''Preformats input file for nlpModeler.'''

from random import shuffle
import re
from unixpath import readFile

class Formatter():

	def __init__(self, infile, outfile, encodingfile):
		self.cancer = []
		self.encodingfile = encodingfile
		self.header = {}
		self.infile = infile
		self.lcount = 1
		self.locations = {"NA": 0}
		self.nas = []
		self.noncancer = []
		self.outfile = outfile
		self.tcount = 1
		self.total = 0
		self.types = {"NA": 0}
		self.__formatFile__()
		self.__writeLists__()
		self.__writeDicts__()

	def __writeDicts__(self):
		# Writes locations and types to file
		with open(self.encodingfile, "w") as out:
			for k in self.locations.keys():
				out.write("Location,{},{}\n".format(k, self.locations[k]))
			for k in self.types.keys():
				out.write("Type,{},{}\n".format(k, self.types[k]))

	def __formatOutput__(self):
		# Returns randomly sorted list, truncates records with na for comments
		shuffle(self.nas)
		ret = self.cancer
		ret.extend(self.noncancer)
		ret.extend(self.nas[:int(len(ret) * 0.05)])
		shuffle(ret)
		print("\tFormatted {} of {} records.".format(len(ret), self.total))
		return ret

	def __writeLists__(self):
		# Writes data to file
		with open(self.outfile, "w") as out:
			header = []
			for i in range(len(self.header)):
				header.append(-1)
			for k in self.header.keys():
				header[self.header[k]] = k
			# Omit Location column
			out.write(",".join(header) + "\n")
			for i in self.__formatOutput__():
				out.write(",".join(i) + "\n")

	def __formatLine__(self, line):
		# Replaces punctuation, splits compound locations and types, and encodes paired locations and types as integers
		rows = []
		go = True
		if line[self.header["Comments"]] == "NA" or line[self.header["Comments"]] == "n/a. n/a.":
			go = False
			if line[self.header["Masspresent"]] != "1":
				# Skip records where diagnosis info is not in comments
				line[-1] = "0"
				line[-2] = "0"
				self.nas.append(line)
		if go:
			line[self.header["Comments"]] = re.sub(r"[^\w\s]", "", line[self.header["Comments"]])
			# Split compound locations and types; store only one for now
			loc = line[self.header["Location"]].split(";")
			for idx, i in enumerate(line[self.header["Type"]].split(";")):
				# Combine location and type as one key
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
			if line[self.header["Masspresent"]] == "1":
				self.cancer.append(row)
			else:
				self.noncancer.append(row)

	def __formatFile__(self):
		# Reads input file, formats values, and writes to output
		first = True
		print("\n\tFormatting input file...")
		for line in readFile(self.infile, header = True, d = ","):
			if not first:
				self.total += 1
				self.__formatLine__(line)
			else:
				self.header = line
				first = False
