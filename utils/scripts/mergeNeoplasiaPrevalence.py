'''Merges tissue specific neoplasia prevalence into one file.'''

from argparse import ArgumentParser
from datetime import datetime
from glob import glob
import os
import unixpath

D = ","

class Species():

	def __init__(self, h, line):
		self.all = line.strip()
		self.h = h
		self.tissues = []
		self.total = None

	def __setTotal__(self, row):
		# Initializes total list


	def add(self, row):
		# Adds new line to totals and appends to tissue list
		loc = row[h["Location"]]

		if not self.total:
			self.__setTotal__()
		

	def string(self):
		# Formats species records for writing	
		ret = []
		ret.append(self.all)
		if self.total:
			ret.append(",".join(self.total))
			for k in self.tissues.keys():
				ret.append(self.tissues[k])
		return "\n".join(ret)

class Merger():

	def __init__(self, args):
		self.header = ""
		self.indir = unixpath.checkDir(args.i)
		self.infiles = glob(self.indir + "*.csv")
		self.outfile = args.o
		self.records = {}

	def __getHeader__(self, row):
		# Returns header dict
		ret = {}
		for idx, i in enumerate(row):
			ret[i] = idx
		return ret

	def mergeFiles(self):
		# Merges all infile data
		h = {}
		for infile in self.infiles:
			first = True
			with open(infile) as f:
				for line in f:
					row = line.strip().split(D)
					if not first:
						uid = row[h["taxa_id"]]
						if uid not in self.records.keys():
							# Store in ordered dict with all tissues row at top
							self.records[uid] = Species(h, line)
						elif row[h["TotalRecords"]] != "0":
							# Store additional tissues in order
							self.records[uid].add(row)						
					else:
						if not h:
							# Set header once since all are the same
							self.header = line
							h = self.__getHeader__(row)
						first = False

	def write(self):
		# Writes dict to file
		with open(self.outfile, "w") as out:
			out.write(self.header)
			for uid in self.records.keys():
				out.write(self.records[uid].string())

def main():
	start = datetime.now()
	parser = ArgumentParser("Merges tissue specific neoplasia prevalence into one file.")
	parser.add_argument("-i", help = "Path to input directory.")
	parser.add_argument("-o", help = "Path to output file.")
	m = Merger(parser.parse_args())
	m.mergeFiles()
	m.write()
	print(("\tTotal runtime: {}\n").format(datetime.now() - start))

if __name__ == "__main__":
	main()
