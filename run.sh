#!/bin/bash

##############################################################################
# Runs golang web application
##############################################################################

HOST=127.0.0.1
PORT=8080
LOG="../serverLog.txt"

killProc () {
	#if lsof -t -i:8080 | grep "*" > /dev/null; then
		echo "Stopping server process."
		kill -9 $(lsof -t -i:8080)
	#fi
}

runHost () {
	# Start gunicorn with nohup and append stderr and stdout to serverLog
	echo "Starting production server on $HOST:$PORT..."
	cd codb/
	nohup go run *.go -h $HOST -p $PORT > $LOG 2>&1 &
}

helpText () {
	echo "Runs hosting server for the comparative oncology database."
	echo ""
	echo "start	Kills running processes and starts new server on port 8080."
	echo "stop	Kills process running on port 8080."
	echo "help	Prints help text and exits."
	echo ""
}

if [ $# -eq 0 ]; then
	helpText
elif [ $1 = "help" ]; then
	helpText
elif [ $1 = "stop" ]; then
	killProc
elif [ $1 = "start" ]; then
	killProc
	runHost
fi
