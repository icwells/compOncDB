'''Plots historgams of record counts.'''

from argparse import ArgumentParser
from datetime import datetime
from matplotlib import pyplot
import numpy as np
import os
import unixpath

class Histograms():

	def __init__(self, args):
		pyplot.style.use("seaborn-deep")
		self.approved ="Approved"
		self.columns = [["Infant", "Castrated"], ["Masspresent", "Necropsy", "Metastasis"], ["Approved", "Zoo"]]
		self.fields = ["Infant", "Castrated", "Masspresent", "Necropsy", "Metastasis", "Zoo"]
		self.id = "ID"
		self.label = [self.approved, "All"]
		self.legend = "upper left"
		self.outdir = unixpath.checkDir(args.o, True)
		self.records = {}
		print()
		for idx, i in enumerate([args.p, args.d, args.s]):
			self.__setTable__(i, self.columns[idx])
		self.__barPlot__()

	def __newRecord__(self):
		# Returns empty record dict
		ret = {}
		for i in self.fields:
			ret[i] = None
		return ret

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

	def __plot__(self, l):
		# Adds histogram to figure pane
		print("\tGenerating plot...")
		fig, ax = pyplot.subplots(nrows = 1, ncols = 1)
		width = 0.4
		ax.bar(self.fields, l[0], width)
		ax.bar(self.fields, l[1], width, bottom = l[0])
		ax.set(title = "Zoo Approvals", ylabel = "Number of Records")
		ax.legend(loc=self.legend, labels = self.label)
		fig.savefig(("{}ZooApproval.{}.svg").format(self.outdir, datetime.now().strftime("%Y-%m-%d")))

	def __barPlot__(self):
		# Plots histograms by related fields
		l = []
		for i in range(len(self.label)):
			l.append([])
			for j in range(len(self.fields)):
				l[i].append(0)
		for k in self.records:
			r = self.records[k]
			for idx, i in enumerate(self.fields):
				if r[i] == 1:
					if r[self.approved] == 1:
						l[0][idx] += 1
					else:
						l[1][idx] += 1
		self.__plot__(l)

def main():
	start = datetime.now()
	parser = ArgumentParser("Plots historgams of record counts.")
	parser.add_argument("d", help = "Path to diagnosis table.")
	parser.add_argument("p", help = "Path to patient table.")
	parser.add_argument("s", help = "Path to source table.")
	parser.add_argument("o", help = "Path to output directory.")
	Histograms(parser.parse_args())
	print(("\tTotal runtime: {}\n").format(datetime.now() - start))

if __name__ == "__main__":
	main()
