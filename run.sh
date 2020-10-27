#!/bin/bash

##############################################################################
# Runs golang web application
##############################################################################

HOST=127.0.0.1
PORT=8080
LOG="../serverLog.txt"
PID=".codb_pid.txt"

killProc () {
	if [ -f $PID ]; then
		echo "Stopping server process."
		kill -9 $(cat $PID)
		rm $PID
	fi
}

runHost () {
	# Start gunicorn with nohup and append stderr and stdout to serverLog
	echo "Starting production server on $HOST:$PORT..."
	cd codb/
	nohup go run *.go -h $HOST -p $PORT > $LOG 2>&1 &
	# Save process ids for easy termination later
	echo $! > ../$PID
}

helpText () {
	echo "Runs hosting server for the comparative oncology database."
	echo ""
	echo "start	Kills running processes and starts new server"
	echo "stop	Kills process using pid in $PID."
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
