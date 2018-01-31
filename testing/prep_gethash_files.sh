#!/bin/sh
# 
# prepare files for calculating hashes. In this case only files that cannot be
# opened and/or read.
#
touch foo       # create file
chmod 000 foo   # remove all permissions