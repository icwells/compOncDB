'''Ranks specfic accounts by number of records.'''

from argparse import ArgumentParser
from datetime import datetime
import os
import unixpath

class AccountCounter():

	def __init__(self, args):
		for i in [args.a, args.i]:
			unixpath.checkFile(i)
		self.accounts = {}
		self.accountfile = args.a
		self.counts = {}
		self.infile = args.i
		self.outfile = args.o
		self.__setAccounts__()
		self.__countRecords__()
		self.__write__()

	def __setAccounts__(self):
		# Reads account names into dict
		print("\n\tReading accounts file...")
		first = True
		for i in unixpath.readFile(self.accountfile, header = True, d = ","):
			if not first:
				if i[head["submitter_name"]] != "NA":
					self.accounts[i[head["account_id"]]] = i[head["submitter_name"]]
			else:
				head = i
				first = False

	def __countRecords__(self):
		# Counts records for each account id
		print("\tReading records file...")
		first = True
		for i in unixpath.readFile(self.infile, header = True, d = ","):
			if not first:
				a = i[head["account_id"]]
				if a != "-1":
					if i[head["Zoo"]] == "1" or i[head["Institute"]] == "1":
						if a not in self.counts.keys():
							# account_id: #records, approved
							self.counts[a] = [0, i[head["Approved"]]]
						self.counts[a][0] += 1
			else:
				head = i
				first = False

	def __write__(self):
		# Writes account names and number of records to file
		print("\tWriting output...")
		with open(self.outfile, "w") as out:
			out.write("Account Name,Number of Records,Approved\n")
			for k in self.counts.keys():
				if k in self.accounts.keys():
					i = self.counts[k]
					row = [self.accounts[k], str(i[0]), i[1]]
					out.write(",".join(row) + "\n")

def main():
	start = datetime.now()
	parser = ArgumentParser("Ranks specfic accounts by number of records.")
	parser.add_argument("-a", help = "Path to accounts file")
	parser.add_argument("-i", help = "Path to records file")
	parser.add_argument("-o", help = "Path to output file")
	AccountCounter(parser.parse_args())
	print(("\tTotal runtime: {}\n").format(datetime.now() - start))

if __name__ == "__main__":
	main()
