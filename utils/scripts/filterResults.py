'''Filters search results based on given account ids'''

from argparse import ArgumentParser
from datetime import datetime
import os
import unixpath

class Filter():
	def __init__(self, args):
		self.infile = args.i
		self.accountfile = args.a
		self.outfile = args.o
		self.accounts = {}
		self.__setAccounts__()

	def __setAccounts__(self):
		# Reads accounts into dict
		with open(self.accountfile, "r") as f:
			for line in f:
				line = line.strip()
				s = line.split("|")
				self.accounts[s[0].strip()] = s[1].strip()

	def filterResults(self):
		# Filters result file based on account id
		first = True
		with open(self.outfile, "w") as out:
			with open(self.infile, "r") as f:
				for line in f:
					line = line.strip()
					if not first:
						s = line.split(d)
						aid = s[idx]
						if aid in self.accounts.keys():
							s.append(self.accounts[aid])
							out.write(",".join(s) + "\n")
					else:
						out.write(line + ",submitter_name\n")
						d = unixpath.getDelim(line)
						for ind, i in enumerate(line.split(d)):
							if i == "account_id":
								idx = ind
								break
						first = False

def main():
	start = datetime.now()
	parser = ArgumentParser("Filters search results based on given account ids.")
	parser.add_argument("-i", help = "Path to search results file.")
	parser.add_argument("-a", help = "Path to account IDs file.")
	parser.add_argument("-o", help = "Path to output file.")
	args = parser.parse_args()
	f = Filter(args)
	f.filterResults()
	print(("\tTotal runtime: {}\n").format(datetime.now() - start))

if __name__ == "__main__":
	main()
