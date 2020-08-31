'''Extracts primate diagnosis data and formats into update file'''

from argparse import ArgumentParser
from datetime import datetime
import os
import pandas as pd
import unixpath

class DiagnosisExtractor():

	def __init__(self, args):
		unixpath.checkFile(args.i)
		self.df = pd.read_excel(args.i)
		self.header = ["ID","Masspresent","Hyperplasia","Metastasis","primary_tumor","Malignant","Type","Location"]
		self.outfile = args.o

	def __getMalignant__(self, i):
		# Extracts malignant/benign
		mb = str(i["Malignant or B9"]).lower()
		if "malignant" in mb:
			i["Malignant"] = "1"
		elif "benign" in mb: 
			i["Malignant"] = "0" 

	def __checkNeoplasia__(self, i):
		# Removes cancer info if record does not contain neoplasia
		d = str(i["Diagnosis"])
		if "NOT NEOPLASIA" in d:
			i["Masspresent"] = "0"
			i["Hyperplasia"] = "0"
			i["Metastasis"] = "0"
			i["primary_tumor"] = "0"
			i["Malignant"] = "0" 
			i["Type"] = "NA"
			i["Location"] = "NA"
			return False
		return True

	def extractDiagnosis(self):
		# Extracts diagnosis dat for update file
		with open(self.outfile, "w") as out:
			out.write(",".join(self.header) + "\n")
			for _, i in self.df.iterrows():
				if self.__checkNeoplasia__(i):
					self.__getMalignant__(i)
				row = []
				go = False
				for h in self.header:
					v = str(i[h])
					if v == "nan":
						v = "NA"
					else:
						go = True
					row.append(v)
				if go:
					out.write(",".join(row) + "\n")

def main():
	start = datetime.now()
	parser = ArgumentParser("Extracts primate diagnosis data and formats into update file.")
	parser.add_argument("-i", help = "Path to input xlsx file.")
	parser.add_argument("-o", help = "Path to output csv.")
	d = DiagnosisExtractor(parser.parse_args())
	print("\n\tGenerating update file...")
	d.extractDiagnosis()
	print(("\tTotal runtime: {}\n").format(datetime.now() - start))

if __name__ == "__main__":
	main()
