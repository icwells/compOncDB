'''Defines TensorFlow model for common/scientifc name classifier'''

from argparse import ArgumentParser
from datetime import datetime
from formatInput import Formatter
import matplotlib.pyplot as plt
import numpy as np
import pandas as pd
from sklearn.model_selection import train_test_split
import tensorflow as tf
from tensorflow.keras.preprocessing.text import Tokenizer
from tensorflow.keras.preprocessing.sequence import pad_sequences
import tensorflow_hub as hub
from unixpath import readFile

INFILE = "diagnoses.csv"
ENCODING = "typeEncodings.csv"

class Classifier():

	def __init__(self):
		#plt.style.use("seaborn-deep")
		self.columns = []
		self.db = None
		self.embedding = 32
		self.epochs = 5
		self.hub = "https://tfhub.dev/google/nnlm-en-dim50/2"
		self.labels_test = {}
		self.labels_train = {}
		self.maxlen = 150
		self.model = None
		self.oov = "<OOV>"
		self.outdir = "diagnosisModel"
		self.padding = "post"
		self.test = []
		self.train = []
		self.training_size = 10000
		self.types = {}
		self.vocab_size = 10000
		self.__loadDicts__()
		self.__getDataFrame__()

	def __loadDicts__(self):
		# Loads types dict
		for i in readFile(ENCODING, header = False, d = ","):
			self.types[int(i[1])] = i[0]

	def __getTokenizer__(self, df):
		# Tokenizes training and testing data
		print("\tTokenizing input data...")
		values = df.pop("Comments").apply(str)
		train, test = values[:self.training_size], values[self.training_size:]
		df.pop("Metastasis")
		df.pop("Necropsy")
		self.columns = list(df.columns)
		for i in self.columns:
			#df[[i]] = df[[i]].astype(np.int32)
			col = np.asarray(df.pop(i)).astype(np.int32)
			self.labels_train[i] = col[:self.training_size].reshape((-1,1))
			self.labels_test[i] = col[self.training_size:].reshape((-1,1))
		#self.labels_train, self.labels_test = df[:self.training_size], df[self.training_size:]
		tokenizer = Tokenizer(num_words = self.vocab_size, oov_token = self.oov)
		tokenizer.fit_on_texts(values)
		self.train = np.array(pad_sequences(tokenizer.texts_to_sequences(train), maxlen = self.maxlen, padding = self.padding, truncating = self.padding))
		self.test = np.array(pad_sequences(tokenizer.texts_to_sequences(test), maxlen = self.maxlen, padding = self.padding, truncating = self.padding))

	def __getDataFrame__(self):
		# Reads dataframe and splits into training and testing datasets
		print("\n\tReading input file...")
		df = pd.read_csv(INFILE, delimiter = ",")
		# Randomly shuffle dataframe
		df = df.sample(frac = 1).reset_index(drop = True)
		self.__getTokenizer__(df)

#-----------------------------------------------------------------------------

	def __plot__(self, history, metric):
		# Plots results
		labels = []
		plt.xlabel("Epochs")
		plt.ylabel(metric)
		for i in self.columns:
			name = "{}_{}".format(i, metric)
			val = "val_" + name
			plt.plot(history.history[name], label = name)
			plt.plot(history.history[val], label = val)
			labels.extend([name, val])
		# Reduce plot size so legend is not covering it
		plt.tight_layout(rect=[0,0,0.65,0.65])
		plt.legend(labels, loc='center left', bbox_to_anchor=(1, 0.5))
		plt.savefig("{}/{}.svg".format(self.outdir, metric), format="svg")
		# Clear plot
		plt.clf()

	def __outputLayer__(self, name, input_layer):
		# Returns new output node
		return tf.keras.layers.Dense(units = 1, activation = "sigmoid", name = name)(input_layer)

	def __typeLayer__(self, input_layer):
		# Returns output layer for types
		dense1 = tf.keras.layers.Dense(units = 128, activation = "relu")(input_layer)
		dense2 = tf.keras.layers.Dense(units = 256, activation = "relu")(dense1)
		flattened = tf.keras.layers.Flatten()(dense2)
		return tf.keras.layers.Dense(units = len(self.types.keys()), activation = "softmax", name = "Type")(flattened)

	def __multiOutputModel__(self):
		# Defines multiple-output model
		outputs = []
		input_layer = tf.keras.layers.Input(shape = (self.maxlen, 1, ))
		# Add 2 bidirectional LSTMs
		bidirectional = tf.keras.layers.Bidirectional(
			tf.keras.layers.LSTM(256, return_sequences=True, name = "forwardLSTM"),
			backward_layer = tf.keras.layers.LSTM(128, return_sequences=True, go_backwards = True, name = "backwardLSTM"),
			name = "BidirectionalLSTM"
		)(input_layer)
		dense1 = tf.keras.layers.Dense(units = 64, activation = "relu")(bidirectional)
		dense2 = tf.keras.layers.Dense(units = 32, activation = "relu")(dense1)
		dense3 = tf.keras.layers.Dense(units = 16, activation = "relu")(dense2)
		flattened = tf.keras.layers.Flatten()(dense3)
		for i in self.columns[:-1]:
			outputs.append(self.__outputLayer__(i, flattened))
		outputs.append(self.__typeLayer__(bidirectional))
		# Define the model with the input layer and a list of output layers
		return tf.keras.Model(inputs = input_layer, outputs = outputs, name = self.outdir)

	def __getLoss__(self):
		# Returns loss estimation for each output column
		ret = {}
		for i in self.columns:
			if i == "Type":
				ret[i] = tf.keras.losses.SparseCategoricalCrossentropy()
			else:
				ret[i] = tf.keras.losses.BinaryCrossentropy()
		return ret
  
	def trainModel(self):
		# Trains species name classifier
		print("\tTraining model...")
		self.model = self.__multiOutputModel__()
		tf.keras.utils.plot_model(self.model, "{}/model_plot.png".format(self.outdir), show_shapes=True)
		self.model.compile(loss = self.__getLoss__(), optimizer = 'adam', metrics = ['accuracy'])
		print(self.model.summary())
		history = self.model.fit(self.train, self.labels_train,
			epochs = self.epochs, 
			batch_size = 512, 
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
	parser = ArgumentParser("")
	parser.add_argument("-i", help = "Path to unformatted training data. Pre-formats the data only. Run again without infile argument to train the model.")
	args = parser.parse_args()
	if args.i:
		Formatter(args.i, INFILE, ENCODING)
	else:
		c = Classifier()
		c.trainModel()
		c.save()
	print(("\tTotal runtime: {}\n").format(datetime.now() - start))

if __name__ == "__main__":
	main()
