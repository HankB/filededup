package main

import (
	"bytes"
	"crypto/md5"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

/* return largest of two arguments
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
		log.Fatal(err)
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return h.Sum(nil)
}

var dbName = "filelist.db"
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
		log.Fatal(err)
	} else {
		rowCount, err := result.RowsAffected()
		if err != nil {
			log.Fatal(err)
		} else {
			if rowCount != 1 {
				log.Fatal("hash update affected %d", rowCount)
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
		log.Fatal(err)
	}
	defer file1.Close()

	file2, err := os.Open(f2)
	if err != nil {
		log.Fatal(err)
	}
	defer file2.Close()

	// read both a block at a time
	var bytesRead int64
	for bytesRead = 0; bytesRead < len; bytesRead += min(blocksize, len-bytesRead) {
		read1, err := file1.Read(buf1[0:min(blocksize, len-bytesRead)])
		if err != nil {
			log.Fatal("bytes read:", read1, " ", err)
		}
		if int64(read1) < min(blocksize, len-bytesRead) {
			log.Fatal("Expected ", min(blocksize, len-bytesRead), " bytes got ", read1)
		}

		read2, err := file2.Read(buf2[0:min(blocksize, len-bytesRead)])
		if err != nil {
			log.Fatal(err)
		}
		if int64(read2) < min(blocksize, len-bytesRead) {
			log.Fatal("Expected ", min(blocksize, len-bytesRead), " bytes got ", read2)
		}

		//fmt.Printf("comparing bytes read:%d, %d: %d\n", read1, read2, bytes.Compare(buf1[0:read1], buf2[0:read2]))
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
func findMatch(filepath string, length int64) (bool, string, []byte) {
	// search the database for files with matching length
	fmt.Printf("matching %s len %d\n", filepath, length)
	rows, err := db.Query(`SELECT length, filename, hash linkCount
							FROM files
							WHERE length=?`,
		length)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// check to see if any were found
	more := rows.Next()
	var hashCandidate []byte = nil
	if more {
		// found similar files - need to calculate the hash of the candidate
		hashCandidate = getHash(filepath) // need hash
		for more {
			var possMatchLen int64 // do we really need this?
			var possMatchFilename string
			var possMatchHash []byte
			if err := rows.Scan(&possMatchLen, &possMatchFilename, &possMatchHash); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("possible: %s is %d\n", possMatchFilename, length)
			if possMatchHash == nil { //need hash for the possible match?
				possMatchHash = getHash(possMatchFilename)
				//updateHash(possMatchFilename, possMatchHash)
			}
			if bytes.Compare(hashCandidate, possMatchHash) == 0 { // matching hash?
				if compareByteByByte(filepath, possMatchFilename, length) { // verify match
					return true, possMatchFilename, possMatchHash
				}
			}
			more = rows.Next()
		}
	} else {
		fmt.Printf("no length matches\n")
		return false, "", nil // no same length files
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return false, "", hashCandidate
}

// callback from Walk()
func myWalkFunc(path string, info os.FileInfo, err error) error {
	if info.Mode()&os.ModeType != 0 { // not a regular file?
		fmt.Printf("other: %s\n", path)
	} else {
		/*
			hash := getHash(path)
			fmt.Printf("len: %d name: %s hash %x\n", info.Size(), path, hash)
			tx, err := db.Begin()
			if err != nil {
				log.Fatal(err)
			}

			stmt, err := tx.Prepare("insert into files(length, filename, hash) values(?, ?, ?)")
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()
			_, err = stmt.Exec(info.Size(), path, hash)
			if err != nil {
				log.Fatal(err)
			}
			tx.Commit()
		*/
		fmt.Printf("checking len: %d name: %s\n", info.Size(), path)
		found, matchPath, hash := findMatch(path, info.Size())
		fmt.Printf("%t, %s, %x\n\n", found, matchPath, hash)
		if found {
			// TODO link files
		} else {
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
func main() {
	initDataBase("sqlite3", dbName)
	defer closeDataBase()
	filepath.Walk("./sample-files", myWalkFunc)

}
