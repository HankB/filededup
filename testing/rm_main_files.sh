#!/bin/sh
#
# remove files to test program operation
#

# first report all linked files
find some-test-dir -type f -printf '%n %p\n' | awk '$1 > 1{$1="";print}'

rm -r some-test-dir