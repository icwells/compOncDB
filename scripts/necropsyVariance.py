'''Identifies species which have significant discrepancies between number of necropsies and non-necropsies'''

from argparse import ArgumentParser
from datetime import datetime
from math import sqrt
import operator
import os
import unixpath

class Record():

	def __init__(self, records, cancer, prev):
		self.cancer = cancer
		self.prevalence = prev
		self.records = records

	def toList(self):
		# Returns list of record values
		return [str(self.records), str(self.cancer), str(self.prevalence)]

class Species():

	def __init__(self, tid, species):
		self.id = tid
		self.necropsy = None
		self.other = None
		self.difference = 0
		self.variance = None
		self.species = species
		self.total = 0

	def toList(self):
		# Returns list of record values
		ret = [self.id, self.species, str(self.total)]
		ret.extend(self.necropsy.toList())
		ret.extend(self.other.toList())
		ret.append(str(self.variance))
		ret.append(str(self.difference))
		return ret

	def setSignificance(self):
		# Stores absolute value of difference between number of records
		self.total = self.necropsy.records + self.other.records
		self.difference = abs(self.necropsy.prevalence - self.other.prevalence)
		if self.necropsy.cancer > self.other.cancer:
			n = self.necropsy.cancer
			p = self.necropsy.prevalence
			x = self.other.prevalence
		else:
			n = self.other.cancer
			p = self.other.prevalence
			x = self.necropsy.prevalence
		self.variance = 2 * sqrt(n * p * (1 - p)) / n

	def setOther(self, records, cancer, prev):
		# Stores non-necropsy total
		if cancer > 0:
			self.other = Record(records, cancer, prev)

	def setNecropsy(self, records, cancer, prev):
		# Stores necropsy total
		if cancer > 0:
			self.necropsy = Record(records, cancer, prev)

class NecropsyVariance():

	def __init__(self, args):
		for i in [args.i, args.n]:
			unixpath.checkFile(i)
		self.min = 50
		self.outfile = args.o
		self.records = {}
		self.rows = []
		print()
		self.__setCounts__(args.i, False)
		self.__setCounts__(args.n, True)
		self.__filter__()
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
							self.records[tid] = Species(tid, sp)
					if tid in self.records.keys():
						rec = self.records[tid]
						n = int(i[header["NeoplasiaWithDenominators"]])
						p = float(i[header["NeoplasiaPrevalence"]])
						if nec:
							rec.setNecropsy(total, n, p)
							if rec.necropsy and rec.other:
								# Calculate significance and store
								rec.setSignificance()
						else:
							rec.setOther(total, n, p)
			else:
				header = i
				first = False

	def __filter__(self):
		# Removes species which are missing a records class and sorts remaining records
		for k in self.records.keys():
			r = self.records[k]
			if r.necropsy and r.other:
				self.rows.append(self.records[k])
		self.rows.sort(key=operator.attrgetter("difference"), reverse=True)

	def __write__(self):
		# Writes records to file
		print("\tWriting records to file...")
		with open(self.outfile, "w") as out:
			out.write("taxa_id,Species,TotalRecords,NecropsyRecords,NecropsyNeoplasia,NecropsyPrevalence,NonNecropsyRecords,\
NonNecropsyNeoplasia,NonNecropsyPrevalence,StandardDeviation,Difference\n")
			for i in self.rows:
				out.write(",".join(i.toList()) + "\n")

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
