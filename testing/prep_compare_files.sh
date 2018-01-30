#!/bin/sh
# 
# prepare files for comparison
#
# cmpfile.4096-1
# cmpfile.4096-2 - identical, one block
# cmpfile.4096-3 - different, one block
dd if=/dev/urandom of=cmpfile.4096-1 bs=4096 count=1 >/dev/null 2>&1
cp cmpfile.4096-1 cmpfile.4096-2
dd if=/dev/urandom of=cmpfile.4096-3 bs=4096 count=1 >/dev/null 2>&1
#
# cmpfile.1024-1
# cmpfile.1024-2 - identical, < one block
# cmpfile.1024-3 - different, < one block
dd if=/dev/urandom of=cmpfile.1024-1 bs=1024 count=1 >/dev/null 2>&1
dd if=/dev/urandom of=cmpfile.1024-3 bs=1024 count=1 >/dev/null 2>&1
cp cmpfile.1024-1 cmpfile.1024-2
#
# cmpfile.5120-1
# cmpfile.5120-2 identical > one block
# cmpfile.5120-3 different > one block
cp cmpfile.4096-1 cmpfile.5120-1; cat cmpfile.1024-1 >>cmpfile.5120-1
cp cmpfile.5120-1 cmpfile.5120-2
cp cmpfile.4096-1 cmpfile.5120-3; cat cmpfile.1024-3 >>cmpfile.5120-3
