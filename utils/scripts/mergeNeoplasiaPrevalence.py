'''Merges tissue specific neoplasia prevalence into one file.'''

from argparse import ArgumentParser
from datetime import datetime
from glob import glob
import os
import pandas as pd
import unixpath

class Merger():

	def __init__(self, args):
		self.avgs = ["PropMalignant", "PropBenign", "AverageAge(months)", "AvgAgeNeoplasia(months)"]
		self.ints = ["TotalRecords", "RecordsWithDenominators", "TotalNeoplasia", "NeoplasiaWithDenominators", "MalignancyKnown", "Malignant",  "Benign", 
					"Male", "Female", "MaleNeoplasia", "FemaleNeoplasia", "Necropsies"]
		self.df = None
		self.floats = ["NeoplasiaPrevalence", "MalignancyPrevalence", "BenignPrevalence"]
		self.infiles = glob(unixpath.checkDir(args.i) + "*.csv")
		self.name = args.n
		self.outfile = args.o
		self.tissues = {}
		self.__setDF__()

	def __getIndex__(self, uid, idx=-1):
		# Returns index of last matching id at idx
		ret = -1
		matches = self.df.loc[self.df["taxa_id"] == uid]
		if len(matches) > 0:
			ret = matches.index[idx] + 1
		return ret

	def __insert__(self, idx, row):
		# Inserts new row at index
		if idx >= 0:
			# Create a lists of index halves
			head = [*range(0, idx, 1)] 
			tail = [*range(idx, self.df.shape[0], 1)] 
			# Increment the value of lower half by 1 and update dataframe index
			tail = [x.__add__(1) for x in tail] 
			self.df.index = head + tail 
			# Insert a row at the end 
			self.df.loc[idx] = row
			self.df = self.df.sort_index()

	def __setDF__(self):
		# Initializes dataframe
		print("\n\tInitializing dataframe with {}...".format(os.path.split(self.infiles[0])[1]))
		self.df = pd.read_csv(self.infiles[0])
		for idx, row in self.df.iterrows():
			if row["Location"] != "all" and row["Location"] != self.name:
				#self.tissues[row["taxa_id"]] = 0
				if row["TotalRecords"] == 0:
					self.df.drop(idx)
					'''# Overwrite blank records
					row["Location"] = self.name
					row["#Sources"] = "NA"
				else:
					self.tissues[row["taxa_id"]] += 1
					series = row.copy()
					series["Location"] = self.name
					series["#Sources"] = "NA"
					self.__insert__(self.__getIndex__(series["taxa_id"], 0), series)'''

	def __addRow__(self, row):
		# Adds row to dataframe
		#self.tissues[row["taxa_id"]] += 1
		self.__insert__(self.__getIndex__(row["taxa_id"]), row)
		'''tid = row["taxa_id"]
		tmp = self.df.loc[self.df["taxa_id"] == tid]
		idx = tmp.loc["Location"] == self.name].index
		total = self.df[idx]
		for i in self.ints:
			total[i] += row[i]
		for i in self.avgs:
			total[i] += row[i]'''

	def mergeFiles(self):
		# Merges all infile data
		for infile in self.infiles[1:]:
			print("\tReading {}...".format(os.path.split(infile)[1]))
			df = pd.read_csv(infile)
			for _, row in df.iterrows():
				if row["Location"] != "all" and row["TotalRecords"] != 0:
					self.__addRow__(row)

	'''def __setRates__(self):
		# Sets proportional rates
		for _, row in self.df.iterrows():
			if row["Location"] == self.name:
				d = self.tissues[row["taxa_id"]]
				if d > 1:
					for i in self.avgs:
						row[i] = row[i]/d
				idx = self.df.loc[self.df["taxa_id"] == row["taxa_id"] & self.df["Location"] == "all"].index
				total = self.df[idx, "RecordsWithDenominators"]
				for i in self.floats:
					row[i] = row[i]/total'''

	def write(self):
		# Writes dict to file
		#self.__setRates__(self)
		self.df.to_csv(self.outfile)

def main():
	start = datetime.now()
	parser = ArgumentParser("Merges tissue specific neoplasia prevalence into one file.")
	parser.add_argument("-i", help = "Path to input directory.")
	parser.add_argument("-n", help = "Name of merged location.")
	parser.add_argument("-o", help = "Path to output file.")
	m = Merger(parser.parse_args())
	m.mergeFiles()
	m.write()
	print(("\tTotal runtime: {}\n").format(datetime.now() - start))

if __name__ == "__main__":
	main()
