'''Creates scatter plots fro life history variables.'''

from argparse import ArgumentParser
from datetime import datetime
from itertools import combinations
from matplotlib import pyplot
from pandas import DataFrame, read_csv
import unixpath

class plotter():

	def __init__(self, args):
		print("\n\tReading input file...")
		self.df = read_csv(args.i, delimiter = ",", header = 0, index_col = 0)
		self.outdir = unixpath.checkDir(args.o, True)
		self.fields = [["female_maturity", "male_maturity", "Gestation", "Weaning", "Infancy"],
					["litter_size", "litters_year", "interbirth_interval", "max_longevity", "metabolic_rate"],
					["birth_weight", "weaning_weight", "adult_weight", "growth_rate"]]

	def __trimUnits__(self, n):
		# Removes units from name
		return n[:n.find("(")]

	def __getColumns__(self, x, y):
		# Returns paired values if both fields are >= 0
		ret = [[], []]
		xvals = self.df[x].tolist()
		yvals = self.df[y].tolist()
		for idx, i in enumerate(xvals):
			if i > 0 and yvals[idx] > 0:
				ret[0].append(i)
				ret[1].append(yvals[idx])
		return ret

	def __plot__(self, x, y):
		# Plots pair of columns and saves to csv
		print(("\tPlotting {} and {}...").format(x, y))
		vals = self.__getColumns__(x, y)
		fig, ax = pyplot.subplots(nrows = 1, ncols = 1)
		ax.scatter(vals[0], vals[1])
		ax.set(title = ("{} vs. {}").format(x, y), ylabel = y, xlabel = x)
		pyplot.xscale("log")
		pyplot.yscale("log")
		fig.savefig(("{}{}-{}.svg").format(self.outdir, self.__trimUnits__(x), self.__trimUnits__(y)))
		pyplot.close(fig)

	def getPlots(self):
		# Plots column pairs
		for col in self.fields:
			pairs = combinations(col, 2)
			for i in pairs:
				self.__plot__(i[0], i[1])

def main():
	start = datetime.now()
	parser = ArgumentParser("Creates scatter plots fro life history variables.")
	parser.add_argument("-i", help = "Path to life history table.")
	parser.add_argument("-o", help = "Path to output directory.")
	p = plotter(parser.parse_args())
	p.getPlots()
	print(("\tTotal runtime: {}\n").format(datetime.now() - start))

if __name__ == "__main__":
	main()
