'''Plots histogram of species totals'''

from argparse import ArgumentParser
from datetime import datetime
from matplotlib import pyplot, ticker
import os
import unixpath

class SpeciesTotals():

	def __init__(self, args):
		unixpath.checkFile(args.a)
		unixpath.checkFile(args.n)
		pyplot.style.use("seaborn-deep")
		self.bins = 10
		self.col = "RecordsWithDenominators"
		self.label = ["all records", "necropsies"]
		self.legend = "upper right"
		self.max = 100
		self.min = 20
		self.outfile = os.path.join(os.path.split(args.a)[0], "speciesTotals.svg")
		print()
		self.all = self.__setCounts__(args.a)
		self.necropsy = self.__setCounts__(args.n)

	def plot(self):
		# Plots counts and writes to csv
		print("\tPlotting species counts...")
		fig, ax = pyplot.subplots(nrows = 1, ncols = 1)
		ax.hist([self.all, self.necropsy], self.bins, label = self.label)
		ax.set(title = "Species Totals", ylabel = "Number of Species", xlabel = "Total Records with Denominators")
		ax.set_xlim(self.min, self.max)
		ax.legend(loc=self.legend)
		fig.savefig(self.outfile)

	def __setCounts__(self, infile):
		# Gets species counts
		l = []
		first = True
		print("\tReading {}...".format(os.path.split(infile)[1]))
		for i in unixpath.readFile(infile, header = True, d = ","):
			if not first:
				n = int(i[header[self.col]])
				if self.min <= n <= self.max:
					l.append(n)
			else:
				header = i
				first = False
		l.sort()
		return l

def main():
	start = datetime.now()
	parser = ArgumentParser("Plots histogram of species totals.")
	parser.add_argument("-a", help = "Path to neoplasia prevalence file for all records.")
	parser.add_argument("-n", help = "Path to neoplasia prevalence file for neoplasia records.")
	s = SpeciesTotals(parser.parse_args())
	s.plot()
	print(("\tTotal runtime: {}\n").format(datetime.now() - start))

if __name__ == "__main__":
	main()
