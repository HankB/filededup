package main

/** Core logic for file deduplication via hard linking
 */

import (
	"bytes"
	"crypto/md5"
	"database/sql"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

/* return smaller of two arguments
 */
func min(a, b int64) int64 {
	if a <= b {
		return a
	}
	return b
}

/* Calculate the hash for the given file
 */
func getHash(filename string) []byte {
	f, err := os.Open(filename)
	if err != nil {
		warnings++
		printf(priWarn, "getHash(): %v\n", err)
		return []byte{0}
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		warnings++
		printf(priWarn, "getHash(): %v\n", err)
		return []byte{0}
	}
	return h.Sum(nil)
}

//var dbName = "filelist.db"
var dbName = ":memory:"
var db *sql.DB
var err error
var hashes map[string][]byte

/* Insert the file into the database
 */
func insertFile(filePath string, length int64, hash []byte) {
	_, err = db.Exec(`insert into files (length, filename, hash) 
				values(?, ?, ?)`, length, filePath, hash)
	if err != nil {
		log.Fatal(err)
	}
}

/* Update the hash for a file already in the database
 */
func updateHash(filePath string, hash []byte) {
	result, err := db.Exec(`update files set hash = ? where filename = ?`,
		hash, filePath)
	if err != nil {
		warnings++
		printf(priWarn, "updateHash(1): %v\n", err)
	} else {
		rowCount, err := result.RowsAffected()
		if err != nil {
			warnings++
			printf(priWarn, "updateHash(2): %v\n", err)
		} else {
			if rowCount != 1 {
				warnings++
				printf(priWarn, "updateHash(3): %d rows affected\n", rowCount)
			}
		}
	}
}

/* Compare two files byte by byte to see if they are identical.
   Files are same length and have matching hashes.
*/
func compareByteByByte(f1, f2 string, len int64) bool {
	const blocksize int64 = 4096
	buf1 := make([]byte, blocksize)
	buf2 := make([]byte, blocksize)

	file1, err := os.Open(f1)
	if err != nil {
		warnings++
		printf(priWarn, "compareByteByByte(1): %v\n", err)
		return false
	}
	defer file1.Close()

	file2, err := os.Open(f2)
	if err != nil {
		warnings++
		printf(priWarn, "compareByteByByte(2): %v\n", err)
		return false
	}
	defer file2.Close()

	// read both a block at a time
	var bytesRead int64
	for bytesRead = 0; bytesRead < len; bytesRead += min(blocksize, len-bytesRead) {
		read1, err := file1.Read(buf1[0:min(blocksize, len-bytesRead)])
		if err != nil {
			warnings++
			printf(priWarn, "compareByteByByte(3): %s: %v\n", f1, err)
		}
		if int64(read1) < min(blocksize, len-bytesRead) {
			warnings++
			printf(priWarn, "compareByteByByte(4): expected %d got %d bytes\n",
				min(blocksize, len-bytesRead), read1)
		}

		read2, err := file2.Read(buf2[0:min(blocksize, len-bytesRead)])
		if err != nil {
			warnings++
			printf(priWarn, "compareByteByByte(5): %s: %v\n", f2, err)
		}
		if int64(read2) < min(blocksize, len-bytesRead) {
			warnings++
			printf(priWarn, "compareByteByByte(6): expected %d got %d bytes\n",
				min(blocksize, len-bytesRead), read2)
		}

		if bytes.Compare(buf1[0:read1], buf2[0:read2]) != 0 { // matching byes?
			return false
		}
	}
	return true
}

/* Check to see if a file matches (identical contents) to something
   in the database. Update the hash for any files in the database that
   need to be checked and do not already have the hash calculated.
*/
func findMatch(filepath string, info os.FileInfo) (bool, string, []byte) {
	// search the database for files with matching length
	rows, err := db.Query(`SELECT length, filename, hash linkCount
							FROM files
							WHERE length=?`,
		info.Size())
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var hashCandidate []byte     // hash of file we're trying to match
	var hashes map[string][]byte // hashes calculated for possible matches
	// check to see if any were found
	for rows.Next() {
		var possMatchLen int64 // do we really need this?
		var possMatchFilename string
		var possMatchHash []byte
		if err := rows.Scan(&possMatchLen, &possMatchFilename, &possMatchHash); err != nil {
			log.Fatal(err)
		}

		// check to see if the possMatchFilename and filepath are already linked
		filepathInfo, err := os.Stat(possMatchFilename)
		if err != nil {
			log.Fatal(err)
		}
		if !os.SameFile(filepathInfo, info) {
			if possMatchHash == nil { //need hash for the possible match?
				possMatchHash = getHash(possMatchFilename)
				// save hashes to update after the present query is closed
				if hashes == nil {
					hashes = make(map[string][]byte)
				}
				hashes[possMatchFilename] = possMatchHash // update later
			}
			if hashCandidate == nil {
				hashCandidate = getHash(filepath)
			}
			if bytes.Compare(hashCandidate, possMatchHash) == 0 { // matching hash?
				if compareByteByByte(filepath, possMatchFilename, info.Size()) { // verify match
					return true, possMatchFilename, possMatchHash
				}
			}
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	rows.Close()

	if hashes != nil {
		for key, value := range hashes {
			updateHash(key, value)
		}
	}

	return false, "", hashCandidate
}

/* Copy file attributes fro file to link
func copyFileAtr(fromFile, toFile string) error {
	return error(1)
}
*/

/* replace newName with link to oldName
 */
func replaceWithLink(oldName, newName string) {
	// first link to a temporary name. Linking with newName still existing
	// will fail. However if newName is one of multiple hard links, the
	// link to temporary and rename will fail.
	for i := 0; i < 999; i++ {
		tmpName := newName + strconv.Itoa(i)
		err := os.Link(oldName, tmpName)
		if err != nil {
			if !os.IsExist(err) {
				printf(priWarn, "Cannot link\"%s\" to \"%s\" :  %s\n",
					oldName, tmpName, err.Error())
				return
			}
		} else {
			// rename temporary to newFile which will overwrite with the link
			err = os.Rename(tmpName, newName)
			if err != nil {
				printf(priWarn, "Cannot rename\"%s\" to \"%s\" :  %s\n",
					tmpName, newName, err.Error())
			}
			return
		}
	}
	printf(priWarn, "Cannot create temp filename \"%s[0:999]\"\n", newName)
}

// callback from Walk()
func myWalkFunc(path string, info os.FileInfo, err error) error {
	if info.Mode()&os.ModeType != 0 { // not a regular file?
		printf(priInfo, "skipping: \"%s\"\n", path)
	} else {
		//fmt.Printf("checking len: %d name: %s\n", info.Size(), path)
		found, matchPath, hash := findMatch(path, info)
		filesConsidered++
		//fmt.Printf("%t, %s, %x\n\n", found, matchPath, hash)
		if found {
			printf(priInfo, "replacing \"%s\" with link to \"%s\"\n", path, matchPath)
			filesLinked++
			bytesSaved += uint64(info.Size())

			if !options.Trial {
				replaceWithLink(matchPath, path)
			}
		} else {
			printf(priInfo, "no match for \"%s\"\n", path)
			insertFile(path, info.Size(), hash)
		}
	}
	return nil
}

/* init the database or go up in flames
 */
func initDataBase(flavor, location string) {
	if location != ":memory:" {
		os.Remove(location) // remove previous copy
	}

	// open the database
	db, err = sql.Open(flavor, location)
	if err != nil {
		log.Fatal(err)
	}

	sqlCreateStmt := `
	create table files (
		length integer not null,
		filename text not null,
		hash blob defult null,
		linkCount integer default 1);
	`
	// create the table
	_, err = db.Exec(sqlCreateStmt)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func closeDataBase() {
	db.Close()
}

var filesConsidered uint64
var filesLinked uint64
var warnings uint64
var bytesSaved uint64

func main() {
	// initi stats here for the benefit of unit tests
	filesConsidered = 0
	filesLinked = 0
	warnings = 0
	bytesSaved = 0

	parseArgs()
	initDataBase("sqlite3", dbName)
	defer closeDataBase()
	filepath.Walk(options.Directory, myWalkFunc)
	if options.Summary {
		printf(priCritcl, "Verbosity %d, Directory \"%s\", Trial %t, Summary %t\n",
			len(options.Verbose), options.Directory, options.Trial, options.Summary)
		printf(priCritcl, "%d files %d linked, %d bytes saved, %d warnings\n",
			filesConsidered, filesLinked, bytesSaved, warnings)
	}
}
