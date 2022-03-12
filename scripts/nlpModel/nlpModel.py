'''Defines TensorFlow model for comparative oncology record diagnosis.'''

from argparse import ArgumentParser
from datetime import datetime
from diagnosisPrediction import loadDiagnoses, Predictor
from formatInput import Formatter
import matplotlib.pyplot as plt
import numpy as np
from os.path import isfile
import pandas as pd
import pickle
from random import shuffle
import tensorflow as tf
import tensorflow_hub as hub
from unixpath import checkDir

DIAGNOSIS = "diagnosisModel"
ENCODING = "typeEncodings.csv"
INFILE = "diagnoses.csv"
NEOPLASIA = "neoplasiaModel"

def shuffleText(val):
	# Returns string with shuffled sentence order
	try:
		val = val.split(".")
		shuffle(val)
		return ".".join(val)
	except AttributeError:
		return val

class Classifier():

	def __init__(self, diag):
		self.batch_size = 128
		self.columns = []
		self.diag = diag
		self.hub = "https://tfhub.dev/google/nnlm-en-dim50-with-normalization/2"
		self.labels_test = {}
		self.labels_train = {}
		self.model = None
		self.test = []
		self.train = []
		self.training_size = 20000
		if self.diag:
			self.epochs = 15
			self.outdir = DIAGNOSIS
			self.types, self.locations = loadDiagnoses(ENCODING)
		else:
			self.epochs = 5
			self.outdir = NEOPLASIA
		# Make sure outdir exsits before saving model so plots can be saved there
		checkDir(self.outdir, True)
		self.plot = "{}/modelPlot.png".format(self.outdir)
		self.__getDataFrame__()

	def __formatData__(self, df, values):
		# Tokenizes training and testing data
		print("\tFormatting labels...")
		self.train, self.test = values[:self.training_size], values[self.training_size:]
		self.columns = list(df.columns)
		for i in self.columns:
			col = np.asarray(df.pop(i)).astype(np.int32)
			self.labels_train[i] = col[:self.training_size].reshape((-1,1))
			self.labels_test[i] = col[self.training_size:].reshape((-1,1))

	def __augmentText__(self, df):
		# Randomly shuffles sentences in comments
		print("\tAugmenting data...")
		l = 2
		mp = df.copy()
		if not self.diag:
			mp.drop(mp[mp["Masspresent"] != 1].index, inplace = True)
			l = 4
		for i in range(l):
			cp = mp.copy()
			cp["Comments"] = cp["Comments"].apply(shuffleText)
			df = df.append(cp)
		return df

	def __getDataFrame__(self):
		# Reads dataframe and splits into training and testing datasets
		print("\n\tReading input file...")
		df = pd.read_csv(INFILE, delimiter = ",")
		if self.diag:
			# Remove non-cancer records and previously modeled fields
			df.drop(df[df["Masspresent"] != 1].index, inplace = True)
			df.pop("Masspresent")
			df.pop("Hyperplasia")
		else:
			# Remove cancer specific values
			df.pop("Type")
			df.pop("Location")
		df = self.__augmentText__(df)
		df = df.sample(frac = 1).reset_index(drop = True)
		values = df.pop("Comments").apply(str)
		self.__formatData__(df, values)

#-----------------------------------------------------------------------------

	def __outputLayer__(self, name, parent_layer, units = 1, activation = "sigmoid"):
		# Returns new output node
		return tf.keras.layers.Dense(units = units, activation = activation, name = name)(parent_layer)

	def __multiOutputModel__(self):
		# Defines multiple-output model
		outputs = []
		leaky_relu = tf.keras.layers.LeakyReLU(alpha=0.01)
		input_layer = tf.keras.layers.Input(shape = [], dtype = tf.string)
		hub_layer = hub.KerasLayer(self.hub, input_shape = [], dtype = tf.string, trainable = True)(input_layer)
		dense = tf.keras.layers.Dense(16, activation = leaky_relu)(hub_layer)
		flattened = tf.keras.layers.Flatten()(dense)
		if self.diag:
			outputs.append(self.__outputLayer__("Location", flattened, len(self.locations), "softmax"))
			outputs.append(self.__outputLayer__("Type", flattened, len(self.types), "softmax"))
		else:
			# Get single masspresent output layer
			outputs.append(self.__outputLayer__("Masspresent", flattened))
			outputs.append(self.__outputLayer__("Hyperplasia", flattened))
		# Define the model with the input layer and a list of output layers
		return tf.keras.Model(inputs = input_layer, outputs = outputs, name = self.outdir)

#-----------------------------------------------------------------------------

	def __plot__(self, history, metric):
		# Plots results
		labels = []
		plt.xlabel("Epochs")
		plt.ylabel(metric)
		if self.diag:
			for i in self.columns:
				name = "{}_{}".format(i, metric)
				val = "val_" + name
				plt.plot(history.history[name], label = name)
				plt.plot(history.history[val], label = val)
				labels.extend([name, val])
			# Reduce plot size so legend is not covering it
			plt.tight_layout(rect=[0, 0, 0.65, 0.65])
			plt.legend(labels, loc = "center left", bbox_to_anchor = (1, 0.5))
		else:
			plt.plot(history.history[metric])
			plt.plot(history.history['val_'+metric])
			plt.legend([metric, 'val_'+metric])
		plt.savefig("{}/{}.svg".format(self.outdir, metric), format = "svg")
		# Clear plot
		plt.clf()

	def __getLoss__(self):
		# Returns loss estimation for each output column
		if self.diag:
			ret = {}
			for i in self.columns:
				ret[i] = tf.keras.losses.SparseCategoricalCrossentropy()
		else:
			ret = tf.keras.losses.BinaryCrossentropy()
		return ret
  
	def trainModel(self):
		# Trains species name classifier
		print("\tTraining model...")
		self.model = self.__multiOutputModel__()
		tf.keras.utils.plot_model(self.model, self.plot, show_shapes = True)
		self.model.compile(loss = self.__getLoss__(), optimizer = "adam", metrics = ["accuracy"])
		print(self.model.summary())
		history = self.model.fit(self.train, self.labels_train,
			epochs = self.epochs, 
			batch_size = self.batch_size, 
			validation_data = (self.test, self.labels_test), 
			verbose = 1
		)
		self.__plot__(history, "accuracy")
		self.__plot__(history, "loss")
		#print(self.model.evaluate(self.test, self.labels_test))

	def save(self):
		# Stores model in outdir
		self.model.save(self.outdir)

def main():
	start = datetime.now()
	parser = ArgumentParser("Defines TensorFlow model for comparative oncology record diagnosis.")
	parser.add_argument("--diagnosis", action = "store_true", default = False, help = "Trains diagnosis identification model. Trains cancer record identification model by default.")
	parser.add_argument("-i", help = "Path to unformatted training data. Pre-formats the data only. Run again without infile argument to train the model.")
	parser.add_argument("-o", help = "Path to output file (for predictions only).")
	args = parser.parse_args()
	if args.o:
		# Predict diagnosis values
		Predictor(args.i, args.o, args.diagnosis, ENCODING, DIAGNOSIS, NEOPLASIA)
	elif args.i:
		# Format input
		Formatter(args.i, INFILE, ENCODING)
	else:
		# Train model
		c = Classifier(args.diagnosis)
		c.trainModel()
		c.save()
	print(("\tTotal runtime: {}\n").format(datetime.now() - start))

if __name__ == "__main__":
	main()
