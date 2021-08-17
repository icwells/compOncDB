'''Produces updated files using commands in config file'''

from argparse import ArgumentParser
from copy import deepcopy
from datetime import datetime
from getpass import getpass
from multiprocessing import Pool, cpu_count
import os
#from sys import stdout
import unixpath

class Command():

	def __init__(self, com, d, user, pw, outdir):
		self.command = com
		self.directory = ""
		if d != "None":
			self.directory = d
		self.outdir = outdir
		self.password = pw
		self.user = user
		self.__formatOutfile__()

	def __formatOutfile__(self):
		# Adds outdir and time stamp to outfile name
		stamp = datetime.now().strftime("%Y-%m-%d")
		idx = self.command.find("-o ")
		tail = self.command[idx:].split()
		if self.outdir != "":
			tail[1] = os.path.join(self.outdir, tail[1])
		tail[1] = tail[1].replace(".csv", ".{}.csv".format(stamp))
		self.command = self.command[:idx] + "-u {} --password {} ".format(self.user, self.password) + " ".join(tail)

	def run(self):
		# Runs command
		name = self.command[:self.command.find("-")].strip()
		print("\tCalling {}...".format(name))
		if self.directory != "":
			os.chdir(self.directory)
		res = unixpath.runProc(self.command)
		if not res:
			print("\n\t[Error] {} call failed.\n".format(name))

class Updater():

	def __init__(self, args):
		self.commands = []
		self.config = "config.txt"
		self.outdir = unixpath.checkDir(args.o, True)
		self.password = getpass(prompt = "\n\tEnter MySQL password: ")
		self.user = args.u
		self.__setConfig__()

	def __setConfig__(self):
		# Stores input file values
		first = True
		for i in unixpath.readFile(self.config, True, "\t"):
			if not first and len(i) > 0 and i[0] != "#":
				self.commands.append(Command(i[header["Command"]], i[header["Directory"]], self.user, self.password, self.outdir))
			elif first:
				header = deepcopy(i)
				first = False

	def runCommand(self, c):
		# Calls command.run
		c.run()

def main():
	start = datetime.now()
	parser = ArgumentParser("Produces updated files using commands in config file. config.txt must be tab seperated and in the same directoy as script.")
	parser.add_argument("-o", help = "Path to output directory.")
	parser.add_argument("-u", help = "Mysql username.")
	u = Updater(parser.parse_args())
	cpu = len(u.commands)
	if cpu > cpu_count():
		cpu = cpu_count()
	pool = Pool(processes = cpu)
	l = len(u.commands)
	print("\n\tIssuing commands...")
	for idx, _ in enumerate(pool.imap_unordered(u.runCommand, u.commands)):
		print("\t{} of {} commands have finished".format(idx+1, l))
	pool.close()
	pool.join()
	print()
	print(("\tTotal runtime: {}\n").format(datetime.now() - start))

if __name__ == "__main__":
	main()
