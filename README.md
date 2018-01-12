# filededup

File deduplication utility. (ssh://git@oak:10022/HankB/filededup.git)

## Summary
Deduplicate files in a directory tree by hard linking identical files. Files
are determined to be duplicates by checking in order

1. file length
2. MD5 hash
3. byte by byte comparison

## Status
* Three tests (length, MD5 hash, byte by byte comparison) working for a small
set of test files
* Hard linking not yet implemented.
* Formal tests not yet implemented. (Tests manually performed against sample files)

## Testing

`go test` from the project root to run Go tests.
`go run go run filededup.go` to run the main applidcation which is a sort of manual test.


## Issues
* Need to address files already hard linked (short cut comparison.)
* Need to address file permissions - cannot link a file we cannot delete.
* Need to implement rename, link, delete and accommodate permissions.
* Command line arguments - directory, verbosity, version.

## Requirements

* `go get github.com/mattn/go-sqlite3`

## Database
* length int - required
* filename text - required
* hash blob - not required, calculated when needed.
* linkCount - default to 1 and incremented for each hard link
## Strategy
Iterate through all files and for each candidate:
1. Query the database for files with matching length. If none match,
insert the candidate into the database.
2. For each database row that matches (e.g. same length file) compare
the MD5 hash for the file. If that matches, perform a byte by byte
compare to verify that the files are identical. If no files in the
database match hash or byte comparison, insert the candidate and
hash into the database.
3. For a candidate that matches length, hash and comparison, link it to
the matching record in the database and increment the link count for the
matching record.