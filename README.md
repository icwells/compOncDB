# The compOncDB scripts are 
This program is meant specifically for managing the comparative oncology database for the Maley lab at Arizona State University.

Copyright 2018 by Shawn Rupp

## Installation
Download the repository:

git clone https://github.com/icwells/compOncDB.git

Most of the scripts are written in python3, but several contain Cython modules which
must be compiled. Cython can be installed from the pypi repository or via Miniconda 
(it is installed by default with the full Anaconda package).

### To install with Miniconda:
conda install cython

### Compiling scripts:
cd compOncDB/

./install.sh

## Please refer to compOncDBReadMe.pdf for more detailed instructions on running the program
