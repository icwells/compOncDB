'''Produces updated files using commands in config file'''

from argparse import ArgumentParser
from datetime import datetime
from getpass import getpass
from multiprocessing import Pool, cpu_count
import os
import unixpath

class Command():

	def __init__(self, com, d, user, pw):
		self.command = com
		self.directory = d
		self.password = getpass(prompt = "\n\tEnter MySQL password: ")
		self.user = user

class Updater():

	def __init__(self, args):
		self.commands = []
		self.config = "config.csv"
		self.outdir = unixpath.checkDir(ags.o, True)
		self.password = None
		self.user = args.u

	def __setConfig__(self):
		# Stores input file values
		first = True
		for i in unixpath.ReadFile(self.config, d = ","):
			if not first:
				self.commands.append(Command(i[header["Command"]], i[header["Directory"]], user, self.password))
			else:
				header = i
				first = False

def main():
	start = datetime.now()
	parser = ArgumentParser("Produces updated files using commands in config file.")
	parser.add_argument("-o", help = "Path to output directory.")
	parser.add_argument("-u", help = "Mysql username.")
	u = Updater(parser.parse_args())
	cpu = len(u.commands)
	if cpu > cpu_count()/2:
		cpu = cpu_count()
	pool = Pool(processes = cpu)
	print("\n\tIssuing commands...")
	for i,_ in enumerate(pool.imap_unordered(func, fasta), 1):


	print(("\tTotal runtime: {}\n").format(datetime.now() - start))

if __name__ == "__main__":
	main()
