# filededup

File deduplication utility. Find identical files and hard link them together.

## Contributing

If you look at this and see opportunities for improvement Feel free to let me know. If you find bugs, file an issue! If you have questions, feel free to email me or file an issue. I would generally be open to PRs as long as they improve the result in some way. Thanks!

## Motivation

My present backup srategy uses `rsync` to

  1. Perform a complete backup of target directories on the first day of every month. These backups are kept indefinitely.
  1. Perform an incremental backup (`--link-dest` option) every day past the first day of the month. These relative backups overwrite the previous month's of the same day of he month.

This results in a lot of duplicate files since the vast majority of the files in the target directories never change or change infrequently. In addition my normal photo workflow results in multiple copies of each picture. The incremental backups produce a lot of hard links (to the first of the month copy.)

If `rsync` is used to copy the backup dataset to another file system, the hard links cannot be preserved. Trying to do so brings `rsync` to a crawl. This causes the dataset to approximately double the disk requirements. I sought a way to recover this 'lost' disk space.

### Solution

I searched for disk file deduplicators for Linux and found one that seemed to do what I wanted: FSlint (http://www.pixelbeat.org/fslint/) I ruled out any that deleted duplicates. I tried using this and it worked if given sufficient time. On my several terabyte backup set it ran for days. It also didn't seem to do anything until near the end of the process. If it was interrupted no deduplication was performed. I wanted something more performant and decided to write it myself. One option was to write this in a shell script. Most of the heavy lifting (navigating directories, calculating hashes, comparing files, hard linking) would performed by well optimized and tested programs and shell. However I was more interested in using `go` which I had used for some other projects. This would provide an opportunity to use `go` (AKA `golang`) for interfacing with the file system and file I/O.

### Performance

"premature optimization is the root of all evil (or at least most of it) in programming." - Knuth

Of course the selection of algorithm can have a huge impact on performance and can be difficult to restructure. I chose to implement this as an incremental process - hard linking as soon as duplicates were discovered - for the following reasons.

 * If the program is interrupted part of the way through, not all processing is wasted. 
 * Information on only one of the duplicates is required to be stored for future comparisons.

 The choice of comparisons is also geared to performance.

 * File length is checked first. This is easily available from the filesystem with minimal overhead. If lengths do not match, further checks are not required.
 * Possible matches are checked to see if they are already hard links to the same file before performing more expensive tests. They could not be linked anyway since they are already linked.
 * The file hash is not calculated until it is needed. For some files it may never be calculated.
 * Byte by byte comparison is only performed if previous criteria are satisfied. This provides an absolute guarantee that the files match but the 'results' cannot be stored as with the hash.

 The MD5 sum was selected as the hash because it is one I had heard about. A brief search seemed to indicate that it provides a reasonable compromise between computational requirements and likelihood of a hash collision. Further, there are command line tools to calculate the hash for disk files and which are handy for testing. Early on I implemented a minimal MD5 calculator in `go` and found it to perform similarly to the shell command `md5sum`.

 One concession to programming efficiency vs. computational efficiency was the choice of Sqlite3 to store the information for each file. Reasons for this included:

 * Ability to examine the data store following a run for testing and debugging purposes.
 * Ability to locate this in a disk file or in memory (or in a disk file mounted in memory - `tmpfs`.)
 * A chance to explore the `go` SQL API.

 A custom data structure would probably perform better but no profiling has yet been performed to demonstrate that this is a bottleneck.

 I suspect that the likely bottleneck for locally mounted filesystems will be disk I/O. No consideration has been made to try to localize disk access. It seems likely that is best left to the filesystem. This will have been tested on both ZFS (Linux) and EXT4 before I consider it 'ready to ship.'

## Summary

Deduplicate files in a directory tree by hard linking identical files. Files
are determined to be duplicates by checking in order

1. file length
1. MD5 hash
1. byte by byte comparison

## Command line arguments

``` shell
hbarta@olive:~/Documents/go-work/src/oak/HankB/filededup$ go run filededup.go filededup_util.go -h
Usage:
  filededup [OPTIONS]

Application Options:
  -v, --verbose  Verbose output
  -d, --dir=     Directory to start (default: .)
  -t, --trial    report instead of performing operations
  -s, --summary  print summary of operations

Help Options:
  -h, --help     Show this help message
```

### Verbosity

    (default) Report only fatal errors (to STDERR via `log.Fatal()`)
    --verbosity / -v => Report all errors including those that do not terminate operation.
    -vv => Report all operations. Will produce at least one line per file examined.

## Status

At present the unit tests pass and this has been tested on small data sets. Testing is in progress on larger data sets.

Punchlist:

* Test as root.
* Remove `linkCount` from database.
* Implement version.
* report counts of hash calculations and file comparisons in summary.

## Testing

* `go test` from the project root to run `go` tests. (Or run individual tests from within VS code)
* `go run filededup.go filededup_util.go -s some-test-dir` to run the main application which is a sort of manual test. Note that if this is executed without the `-d` argument it will link files in the project directory.
* Further testing is performed on larget data sets utilizing the ZFS filesystem. This is chosen because I have a new backup server not yet commissioned on which to test. Since ZFS implements snapshots I can snapshot the filesystem, test and revert to test again.

## Issues/Challenges

* Need to address file permissions - cannot link a file we cannot delete. Note: Permissions and ownership match the file linked to.
* Command line arguments - directory, verbosity, version. (version not yet implemented)
* file permissions: file/directory permissions may prevent replacing a fle with a link. (In this case an error will be reported given sufficient verbosity.)
* Changing file ownership requires root priveledge.
* race conditions matching and then performing the link/replace operation.

## Requirements

(Relative to a default Ubuntu 16.04 desktop install.)

* `sudo apt install sqlite3` (Debian, Ubuntu, derivatives) This is required for `go test`, not required to build and execute.
* `go get github.com/HankB/filededup` will pull in dependencies. They can also be installed using the following two commands.
* `go get github.com/mattn/go-sqlite3`
* `go get github.com/jessevdk/go-flags`

### Operating system

This has been developed and (minimally) tested on Ubuntu (16.04 LTS and 17.10) and Debian Stretch. In theory the `go` code should run on Microsoft Windows but the unit tests use Bourne shell scripts to set up the test environment. Windows would also require NTFS or other file system that supports hard links. This should work easily on Mac OS if the Bourne shell is installed and available.

## Database

* length int - required
* filename text - required
* hash blob - not required, calculated when needed.
* linkCount - default to 1 and incremented for each hard link. (not implemented)

## Strategy

Iterate through all files and for each candidate:

1. Query the database for files with matching length. If none match, insert the candidate into the database.
1. For each database row that matches (e.g. same length file) compare the MD5 hash for the file. If that matches, perform a byte by byte compare to verify that the files are identical. If no files in the database match hash or byte comparison, insert the candidate and hash into the database.
1. For a candidate that matches length, hash and comparison, link it to the matching record in the database and increment the link count for the matching record.
