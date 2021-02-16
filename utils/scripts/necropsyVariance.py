'''Identifies species which have significant discrepancies between number of necropsies and non-necropsies'''

from argparse import ArgumentParser
from datetime import datetime
from numpy import std
import os
import unixpath

class Record():

	def __init__(self, species):
		self.difference = None
		self.other = None
		self.necropsy = None
		self.significant = 0
		self.species = species

	def toList(self):
		# Returns list of record values
		return [self.species, str(self.other + self.necropsy), str(self.other), str(self.necropsy), str(self.difference), str(self.significant)]

	def setDifference(self):
		# Stores absolute value of difference between number of records
		if self.other and self.necropsy:
			self.difference = abs(self.other - self.necropsy)

	def setOther(self, val):
		# Stores non-necropsy total
		self.other = val

	def setNecropsy(self, val):
		# Stores necropsy total
		self.necropsy = val		

class NecropsyVariance():

	def __init__(self, args):
		for i in [args.i, args.n]:
			unixpath.checkFile(i)
		self.col = "RecordsWithDenominators"
		self.min = 50
		self.outfile = args.o
		self.records = {}
		self.variance = []
		print()
		self.__setCounts__(args.i, False)
		self.__setCounts__(args.n, True)
		# Calculate two standard deviations
		self.sd = std(self.variance)
		self.__filter___()
		self.__write__()

	def __setCounts__(self, infile, nec):
		# Gets species counts
		first = True
		print("\tReading {}...".format(os.path.split(infile)[1]))
		for i in unixpath.readFile(infile, header = True, d = ","):
			if not first:
				total = int(i[header["RecordsWithDenominators"]])
				if self.min < total:
					tid = i[header["taxa_id"]]
					if tid not in self.records.keys():
						sp = i[header["Species"]]
						if sp != "NA":
							self.records[tid] = Record(sp)
					if tid in self.records.keys():
						n = int(i[header[self.col]])
						if nec:
							self.records[tid].setNecropsy(n)
							# Calculate difference and store
							self.records[tid].setDifference()
							if self.records[tid].difference:
								self.variance.append(self.records[tid].difference)
						else:
							self.records[tid].setOther(n)
			else:
				header = i
				first = False

	def __filter___(self):
		# Identifies species with significant deviations in necropsy/non-necropsy counts
		rm = []
		sd2 = self.sd * 2
		print("\tFiltering records...")
		for k in self.records.keys():
			if not self.records[k].difference:
				rm.append(k)
			elif self.records[k].difference > sd2:
				self.records[k].significant = 2
			elif self.records[k].difference > self.sd:
				self.records[k].significant = 1
		for k in rm:
			self.records.pop(k)

	def __write__(self):
		# Writes records to file
		sd = str(self.sd)
		print("\tWriting records to file...")
		with open(self.outfile, "w") as out:
			out.write("taxa_id,Species,TotalRecords,NonNecropsyRecords,NecropsyRecords,Difference,Significance,StandardDeviation\n")
			for k in self.records.keys():
				if self.records[k].significant > 0:
					row = [k]
					row.extend(self.records[k].toList())
					row.append(sd)
					out.write(",".join(row) + "\n")

def main():
	start = datetime.now()
	parser = ArgumentParser("Identifies species which have significant discrepancies between number of necropsies and non-necropsies.")
	parser.add_argument("-i", help = "Path to neoplasia prevalence file for non-necropsy records.")
	parser.add_argument("-n", help = "Path to neoplasia prevalence file for necropsy records.")
	parser.add_argument("-o", help = "Path to output file.")
	NecropsyVariance(parser.parse_args())
	print(("\tTotal runtime: {}\n").format(datetime.now() - start))

if __name__ == "__main__":
	main()
