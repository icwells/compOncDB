'''Defines TensorFlow model for common/scientifc name classifier'''

from argparse import ArgumentParser
from datetime import datetime
import matplotlib.pyplot as plt
import numpy as np
import tensorflow as tf
import tensorflow_hub as hub
#from tensorflow.keras.preprocessing.text import Tokenizer
#from tensorflow.keras.preprocessing.sequence import pad_sequences
from unixpath import readFile

class Classifier():

	def __init__(self, args):
		self.db = None
		self.embedding = None
		self.epochs = 20
		self.infile = args.i
		self.labels_test = []
		self.labels_train = []
		self.model = None
		self.module = 
		self.outfile = "diagnosisModel"
		self.test = []
		self.train = []
		self.training_size = 10000
		self.__connect__()
		self.__getDataSets__()

	def __getDataSets__(self):
		# Extracts common and scientific lists from database
		labels = []
		comments = []
		train = []
		first = True
		print("\n\tReading training data...")
		for i in readFile(self.infile, d = ","):
			if not first:
				comments.append(i[0])
			else:
				header = i
		# Get training and testing sets
		'''for i in names:
			# Split labels and terms after shuffling
			labels.append(i[0])
			train.append(i[1])
		self.labels_train = np.array(labels[:self.training_size])
		self.labels_test = np.array(labels[self.training_size:])
		self.train = np.array(train[:self.training_size])
		self.test = np.array(train[self.training_size:])'''

	def __generate_embeddings__(self, text, random_projection_matrix=None):
	  # Beam will run this function in different processes that need to
	  # import hub and load embed_fn (if not previously loaded)
	  if self.embedding is None:
		self.embedding = hub.load(self.module)
	  embedding = self.embedding(text).numpy()
	  if random_projection_matrix is not None:
		embedding = embedding.dot(random_projection_matrix)
	  return text, embedding

	def __plot__(self, history, metric):
		# Plots results
		plt.plot(history.history[metric])
		plt.plot(history.history['val_'+metric])
		plt.xlabel("Epochs")
		plt.ylabel(metric)
		plt.legend([metric, 'val_'+metric])
		plt.savefig("{}.svg".format(metric), format="svg")
		# Clear plot
		plt.clf()
  
	def trainModel(self):
		# Trains species name classifier
		print("\tTraining model...")
		hub_layer = hub.KerasLayer(self.hub, input_shape=[], dtype=tf.string, trainable=True)
		self.model = tf.keras.Sequential([
			hub_layer,
			tf.keras.layers.Dense(32),
			tf.keras.layers.Dense(16, activation='relu'),
			tf.keras.layers.Dense(1, activation="sigmoid")
		])
		self.model.compile(loss='binary_crossentropy', optimizer='adam', metrics=['accuracy'])
		print(self.model.summary())
		history = self.model.fit(self.train, self.labels_train, 
				epochs = self.epochs, 
				batch_size = 512, 
				validation_data = (self.test, self.labels_test), 
				verbose = 2
		)
		print(self.model.evaluate(self.test, self.labels_test))
		self.__plot__(history, "accuracy")
		self.__plot__(history, "loss")

	def save(self):
		# Stores model in outfile
		self.model.save(self.outfile)

def main():
	start = datetime.now()
	parser = ArgumentParser("")
	parser.add_argument("i", help = "Path to training data.")
	c = Classifier(parser.parse_args())
	c.trainModel()
	c.save()
	print(("\tTotal runtime: {}\n").format(datetime.now() - start))

if __name__ == "__main__":
	main()
