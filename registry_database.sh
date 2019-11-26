#! /bin/bash
if [ "$EUID" -ne 0 ]
  then echo "Please run as root"
  exit
fi
source .env
# create database and user
mysql -e "CREATE DATABASE $DB"
mysql -e "GRANT ALL PRIVILEGES ON $DB.* TO $DB@localhost IDENTIFIED BY '$PASSWORD!'"