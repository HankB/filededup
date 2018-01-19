#!/bin/sh
#
# prepare files to test files already hard linked. 
# Link is not preserved by git.
#

echo x > sample-files/x
ln sample-files/x sample-files/y
