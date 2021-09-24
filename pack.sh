#!/bin/sh
make &&
	tar -zcvf go-stocks-$(date --iso-8601).tar.gz config.toml go-stocks templates static
