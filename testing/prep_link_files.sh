#!/bin/sh
#
# prepare files to test hard linking
#
touch a b 
mkdir test_dir
touch test_dir/c test_dir/d
chmod 555 test_dir
