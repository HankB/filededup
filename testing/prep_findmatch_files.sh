#!/bin/sh
#
# prepare files to test files already hard linked. 
# Link is not preserved by git.
#

echo x > sample-files/z
echo x > sample-files/x
ln sample-files/x sample-files/y

# accidentally linking some files (e.g run the executable in the project
# directory) causes tests to fail. Fix that before proceeding

cp sample-files/empty sample-files/tmp && mv sample-files/tmp sample-files/empty
cp sample-files/'another file' sample-files/tmp && mv sample-files/tmp sample-files/'another file'
exit 0