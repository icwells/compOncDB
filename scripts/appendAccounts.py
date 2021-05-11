'''Appends account names to records.'''

from argparse import ArgumentParser
from datetime import datetime
import os
import unixpath

class Accounts():

	def __init__(self, args):
		self.accounts = {}
		self.__setAccounts__(args.a)
		self.__appendAccounts__(args.i, args.o)

	def __setAccounts__(self, infile):
		# Stores account names in dict
		first = True
		print("\n\tReading accounts file...")
		for i in unixpath.readFile(infile, d = ",", header = True):
			if not first:
				self.accounts[i[head["account_id"]]] = i[head["submitter_name"]]
			else:
				head = i
				first = False

	def __appendAccounts__(self, infile, outfile):
		# Appends account name to records
		first = True
		with open(outfile, "w") as out:
			for i in unixpath.readFile(infile, d = ",", header = True):
				if not first:
					aid = i[head["account_id"]]
					if aid in self.accounts.keys():
						i.append(self.accounts[aid])
					out.write(",".join(i) + "\n")
				else:
					head = i
					header = ["" for n in range(len(head))]
					for k in head.keys():
						header[head[k]] = k
					header.append("Account")
					out.write(",".join(header) + "\n")
					first = False

def main():
	start = datetime.now()
	parser = ArgumentParser("Appends account names to records.")
	parser.add_argument("-a", help = "Path to accounts file.")
	parser.add_argument("-i", help = "Path to records file.")
	parser.add_argument("-o", help = "Path to output file.")
	Accounts(parser.parse_args())
	print(("\tTotal runtime: {}\n").format(datetime.now() - start))

if __name__ == "__main__":
	main()
