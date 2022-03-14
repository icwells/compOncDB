'''Returns predictions for neoplasia and diagnosis models.'''

from argparse import ArgumentParser
from datetime import datetime
import numpy as np
import tensorflow as tf
from unixpath import readFile

def toList(d):
	# Converts dict to ordered list
	ret = ["" for k in range(len(d.keys()))]
	for k in d.keys():
		ret[k] = d[k]
	return ret

def loadDiagnoses(encoding):
	# Loads types and locations lists
	types, locations = {}, {}
	for i in readFile(encoding, header = False, d = ","):
		if i[0] == "Type":
			types[int(i[2])] = i[1]
		else:
			locations[int(i[2])] = i[1]
	return toList(types), toList(locations)

class Predictor():

	def __init__(self, infile, outfile, diag, encoding, diagnosis, neoplasia):
		print("\n\tLoading NLP model...")
		self.comments = []
		self.header = None
		self.ids = []
		self.infile = infile
		self.outfile = outfile
		self.res = {}
		self.__getComments__()
		if diag:
			self.header = "ID,Comments,Location,Lscore,Type,Tscore"
			self.model = tf.keras.models.load_model(diagnosis)
			self.types, self.locations = loadDiagnoses(encoding)
			self.__predictDiagnoses__()
		else:
			self.header = "ID,Comments,Neoplasia,Hyperplasia"
			self.model = tf.keras.models.load_model(neoplasia)
			self.__predictNeoplasia__()
		self.__write__()

	def __write__(self):
		# Writes to output file
		print("\tWriting to file...")
		with open(self.outfile, "w") as out:
			if self.header:
				out.write(self.header + "\n")
			for k in self.res:
				row = ",".join(self.res[k])
				out.write("{}\n".format(row))

	def __getComments__(self):
		# Reads in single column of input names
		print("\tReading input file...")
		with open(self.infile, "r") as f:
			for line in f:
				line = line.split(",")
				self.ids.append(line[0].strip())
				self.comments.append(line[1].strip())

	def __predictDiagnoses__(self):
		# Predicts location and tumor types
		print("\tClassifying neoplasia records...")
		for ldx, label in enumerate(self.model.predict(self.comments)):
			for idx, i in enumerate(label):
				pid = self.ids[idx]
				if pid not in self.res.keys():
					self.res[pid] = [pid, self.comments[idx]]
				ind = np.argmax(i)
				try:
					# Append prediction and score
					if ldx == 0:
						val = [self.types[ind], str(i[ind])]
					elif ldx == 1:
						val = [self.locations[ind], str(i[ind])]
				except IndexError:
					val = ["NA", "-1"]
				self.res[pid].extend(val)

	def __predictNeoplasia__(self):
		# Predicts whether name is common/scientific
		print("\tClassifying neoplasia records...")
		for idx, i in enumerate(self.model.predict(self.comments)):
			pid = self.ids[idx]
			self.res[pid] = [pid, self.comments[idx], str(i[0]), str(i[1])]
