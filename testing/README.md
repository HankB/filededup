# Testing

Reviewing the project after not touching it for 4 years and with plans to offer the option for a persistent database, more testing is in order. 

## Bulk test

Run the program on a large data set and perform some analyses to determine if everything seems to be working and no files are leaking. The test will be run on the entire (ZFS) dataset to avoid complications with identical files outside the scope.

1. Copy data using `rsync` to break any existing hard links.
1. Or revert to the snapshot created in the next step.
1. Snapshot the dataset to provide a known starting point.
1. Capture a list of all files using `find "${DIR}" -type f -printf '%n %p\n'` or similar. Perhaps wise to capture lists with and w/out link counts.
1. Run the program as is.
1. Capture output with and without link counts.
1. Verify that all files remain.
1. Review the link counts to determine if any further validation can be performed.
1. Perform the programming changes to persist the database. 
1. Repeat the testing process above with and without persistence. Also perform some runs that are interrupted during execution and restarted.