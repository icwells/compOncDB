'''Plots historgams of record counts.'''

from argparse import ArgumentParser
from datetime import datetime
from matplotlib import pyplot
import os
from unixpath import *

class Histograms():

	def __init__(self, args):
		pyplot.style.use("seaborn-deep")
		self.axes = setAxes()
		self.columns = [["Masspresent", "Necropsy"], ["Approved"]]
		self.combinations = [["Approved", "Necropsy"], ["Approved", "Masspresent"]]
		self.id = "ID"
		self.label = ["True", "False", "NA"]
		self.legend = "upper right"
		self.outdir = unixpath.checkDir(args.o, True)
		self.records = {}
		print()
		for idx, i in enumerate([args.p, args.n, args.s]):
			self.__setTable__(i, self.columns[idx])
		self.__plotHistograms__()

	def __newRecord__(self):
		# Returns empty record dict
		return {"Infant": None, "Masspresent": None, "Necropsy": None, "Approved": None}

	def __setTable__(self, infile, columns):
		# Returns list of table columns
		first = True
		print("\tReading {}...".format(os.path.split(infile)[1]))
		for line in unixpath.readFile(infile, header = True, d = ","):
			if not first:
				pid = line[head[self.id]]
				if pid not in self.records.keys():
					self.records[pid] = self.__newRecord__()
				for i in columns:
					val = line[head[i]]
					if val:
						try:
							self.records[pid][i] = int(val)
						except:
							pass
			else:
				head = line
				first = False

	def __plot__(self, name, l):
		# Adds histogram to figure pane
		print(("\tPlotting {}...").format(name))
		fig, ax = pyplot.subplots(nrows = 1, ncols = 1)

		ax.hist([ ], label = self.label)
		ax.set(title = name, ylabel = "Frequency", xlabel = self.axes[k].label)
		#ax.set_xlim(0)
		ax.legend(loc=self.legend)
		fig.savefig(("{}{}.{}.svg").format(self.outdir, name, datetime.now().strftime("%Y-%m-%d")))

	def __plotHistograms__(self):
		# Plots histograms by related fields
		for i in self.combinations:
			l = []
			for k in self.records:
				r = self.records[k]
				if r[i[0]] and r[i[1]]:
					l.append([r[i[0]], r[i[1]]])
			#self.__trimLists__(k)
			self.__plot__("{} {}".format(r[i[0]], r[i[1]]), l)

def main():
	start = datetime.now()
	parser = ArgumentParser("Plots historgams of record counts.")
	parser.add_argument("-d", help = "Path to diagnosis table.")
	#parser.add_argument("-p", help = "Path to patient table.")
	parser.add_argument("-o", help = "Path to output directory.")
	parser.add_argument("-s", help = "Path to source table.")
	Histograms(parser.parse_args())
	print(("\tTotal runtime: {}\n").format(datetime.now() - start))

if __name__ == "__main__":
	main()
