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
		self.columns = []
		self.db = None
		self.df = None
		self.embedding = 32
		self.epochs = 20
		self.hub = "https://tfhub.dev/google/nnlm-en-dim50/2"
		self.labels_test = []
		self.labels_train = []
		self.locations = {}
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
		# Loads locations and types dict
		for i in readFile(ENCODING, header = False, d = ","):
			if i[0] == "Location":
				self.locations[int(i[2])] = i[1]
			else:
				self.types[int(i[2])] = i[1]

	def __getTokenizer__(self):
		# Tokenizes training and testing data
		print("\tTokenizing input data...")
		values = self.df.pop("Comments").apply(str)
		self.columns = list(self.df.columns)
		for i in self.df.columns:
			self.df[[i]] = self.df[[i]].astype(np.int32)
		train, self.labels_train = values[:self.training_size], self.df[:self.training_size]
		test, self.labels_test = values[self.training_size:], self.df[self.training_size:]
		tokenizer = Tokenizer(num_words = self.vocab_size, oov_token = self.oov)
		tokenizer.fit_on_texts(values)
		self.train = np.array(pad_sequences(tokenizer.texts_to_sequences(train), maxlen = self.maxlen, padding = self.padding, truncating = self.padding))
		self.test = np.array(pad_sequences(tokenizer.texts_to_sequences(test), maxlen = self.maxlen, padding = self.padding, truncating = self.padding))

	def __getDataFrame__(self):
		# Reads dataframe and splits into training and testing datasets
		print("\n\tReading input file...")
		self.df = pd.read_csv(INFILE, delimiter = ",")
		# Randomly shuffle dataframe
		self.df = self.df.sample(frac = 1).reset_index(drop = True)
		'''values = self.df.pop("Comments").apply(str)
		self.train, self.test = train_test_split(values, test_size = 1 - self.training_size/len(values), random_state = 1)
		self.labels_train, self.labels_test = train_test_split(self.df, test_size = 1 - self.training_size/len(self.df), random_state = 1)'''
		self.__getTokenizer__()

	def __plot__(self, history, metric):
		# Plots results
		for i in self.columns:
			name = "{}_{}".format(name, metric)
			plt.plot(history.history[name])
			plt.plot(history.history["val_" + name])
			plt.xlabel("Epochs")
			plt.ylabel(name)
			plt.legend([name, "val_" + name])
		plt.savefig("{}/{}.svg".format(self.outdir, metric), format="svg")
		# Clear plot
		plt.clf()

	'''def __getFeatureColumns__(self):
		# Returns feature columns
		types = tf.feature_column.categorical_column_with_vocabulary_list(key = "Type", vocabulary_list = self.df["Type"].unique()),
		locations = tf.feature_column.categorical_column_with_vocabulary_list(key = "Location", vocabulary_list = self.df["Location"].unique())
		return [
			tf.feature_column.categorical_column_with_vocabulary_list(key = "Comments", vocabulary_list = self.df["Comments"].unique()),
			tf.feature_column.numeric_column(key = "Masspresent", dtype = tf.int32),
			tf.feature_column.numeric_column(key = "Hyperplasia", dtype = tf.int32),
			tf.feature_column.numeric_column(key = "Necropsy", dtype = tf.int32),
			tf.feature_column.numeric_column(key = "Metastasis", dtype = tf.int32),
			tf.feature_column.numeric_column(key = "primary_tumor", dtype = tf.int32),
			# Cross types and locations to weigh combinations of each
			tf.feature_column.crossed_column(keys=[types, locations], hash_bucket_size = 100)
		]

	def __lstmModel__(self):
		# Defines bidirectional lstm model
		print("\tBuilding bidirectional LSTM model...")
		return tf.keras.Sequential([
			tf.keras.layers.Embedding(self.vocab_size, self.embedding, input_length = self.maxlen),
			tf.keras.layers.Bidirectional(tf.keras.layers.LSTM(64, return_sequences=True)),
			tf.keras.layers.Bidirectional(tf.keras.layers.LSTM(32)),
			tf.keras.layers.Dense(16, activation='relu'),
			tf.keras.layers.Dense(1, activation="sigmoid")
		])'''

	def __outputLayer__(self, name, input_layer):
		# Returns new output node
		return tf.keras.layers.Dense(units = "1", name = name)(input_layer)

	def __multiOutputModel__(self):
		# Defines multiple-output model
		outputs = []
		#input_layer = hub.KerasLayer(self.hub, input_shape=[], dtype=tf.string, trainable=True)
		input_layer = tf.keras.layers.Input(shape = (self.maxlen, ))
		#embed = tf.keras.layers.Embedding(self.vocab_size, self.embedding, input_length = self.maxlen),
		# Add 2 bidirectional LSTMs
		bidirectional = tf.keras.layers.Bidirectional(
			tf.keras.layers.LSTM(64, return_sequences=True)
		)(input_layer)
		first_dense = tf.keras.layers.Dense(units='128', activation='relu')(bidirectional)
		second_dense = tf.keras.layers.Dense(units='32', activation='sigmoid')(first_dense)
		for i in self.columns:
			outputs.append(self.__outputLayer__(i, second_dense))
		# Define the model with the input layer and a list of output layers
		return tf.keras.Model(inputs = input_layer, outputs = outputs)
  
	def trainModel(self):
		# Trains species name classifier
		print("\tTraining model...")
		self.model = self.__multiOutputModel__()
		tf.keras.utils.plot_model(self.model, "{}/model_plot.png".format(self.outdir), show_shapes=True)
		self.model.compile(loss='binary_crossentropy', optimizer='adam', metrics=['accuracy'])
		print(self.model.summary())
		history = self.model.fit(self.train, self.labels_train,
			epochs = self.epochs, 
			batch_size = 512, 
			validation_data = (self.test, self.labels_test), 
			verbose = 2
		)
		self.__plot__(history, "accuracy")
		self.__plot__(history, "loss")
		print(self.model.evaluate(self.test, self.labels_test))

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
