#!/bin/bash

##############################################################################
# Creates new mysql user with basic select privileges.
##############################################################################

ROOT=root
SCRIPT=newUser.sql


newUser () {
	mysql -u $ROOT -p -e "CREATE USER '${USERNAME}'@'%s' IDENTIFIED BY '${USERNAME}';
GRANT SELECT ON comparativeOncology.Records TO '${USERNAME}'@'%';
GRANT SELECT ON comparativeOncology.Update_time TO '${USERNAME}'@'%';"
}

helpText () {
	echo ""
	echo "Creates new mysql user with basic select privileges. Must be run locally."
	echo "Usage: ./newUser.sh {username}"
	echo ""
	echo "help		Prints help text."
}

if [ $# -eq 1 ]; then
	USERNAME=$1
	newUser
elif [ $1 = "help" ]; then
	helpText
else
	helpText
fi

echo ""
