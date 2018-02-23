#!/bin/sh
#
# prepare files to test program operation
#
cp -R sample-files some-test-dir
cd some-test-dir

# start with osme multiply linked files
echo "hello world"  > a
cp a b
ln a c 
ln a d
ln b e
ln b f
