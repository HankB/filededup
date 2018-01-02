# filededup

File deduplication utility.

## Summary
Deduplicate files in a directory tree by hard linking identical files. Files
are determined to be dupolicates by checking in order

1. file length
2. MD5 hash
3. byte by byte comparison

## Database
* length int - required
* filename text - required
* hash blob - not required, calculated when needed.
* link count - default to 1 and incremented for each hard link
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