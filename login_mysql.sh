#! /bin/bash
source .env
mysql -u $DB -D $DB --password=$PASSWORD