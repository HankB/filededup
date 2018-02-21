# filededup

File deduplication utility. (ssh://git@oak:10022/HankB/filededup.git)

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

* Three comparisons (length, MD5 hash, byte by byte comparison) working for a small set of test files
* test as root.

## Testing

* `go test` from the project root to run Go tests. (Or run individual tests from within VS code)
* `go run filededup.go filededup_util.go -s some-test-dir` to run the main application which is a sort of manual test. Note that if this is executed without the `-d` argument it will link files in the `sample-files` directory.

## Issues/Challenges

* Need to address file permissions - cannot link a file we cannot delete. Note: Permissions and ownership match the file linked to.
* Command line arguments - directory, verbosity, version.
* file permissions: file/directory permissions may prevent replacing a fle with a link. 
* Changing file ownership requires root priveledge.
* race conditions matching and then performing the link/replace operation.

## Requirements

(Relative to a default Ubuntu 16.04 desktop install.)

* `sudo apt install sqlite3`
* `go get github.com/mattn/go-sqlite3`
* `go get github.com/jessevdk/go-flags`
* `sudo apt install sqlite3` (Debian, Ubuntu, derivatives)

## Database

* length int - required
* filename text - required
* hash blob - not required, calculated when needed.
* linkCount - default to 1 and incremented for each hard link.

## Strategy

Iterate through all files and for each candidate:

1. Query the database for files with matching length. If none match, insert the candidate into the database.
1. For each database row that matches (e.g. same length file) compare the MD5 hash for the file. If that matches, perform a byte by byte compare to verify that the files are identical. If no files in the database match hash or byte comparison, insert the candidate and hash into the database.
1. For a candidate that matches length, hash and comparison, link it to the matching record in the database and increment the link count for the matching record.