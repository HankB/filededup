package main

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"bytes"

	_ "github.com/mattn/go-sqlite3"
)
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
/* Insert the file into the database
*/
func insertFile(filePath string, length int64, hash []byte) {
	_, err = db.Exec(`insert into files (length, filename, hash) 
				values(?, ?, ?)`, length, filePath, hash)
	if err != nil {
		log.Fatal(err);
	} 
}

/* Update the hash for a file already in the database
*/
func updateHash(filePath string, hash []byte) {
	result, err := db.Exec(`update files set hash = ? where filename = ?`,
			hash, filePath)
	if err != nil {
		log.Fatal(err);
	} else {
		rowCount, err := result.RowsAffected()
		if err != nil {
			log.Fatal(err);
		} else {
			if rowCount != 1 {
				log.Fatal("hash update affected %d", rowCount);
			}
		}
	}
}

/* Check to see if a file matches (identical contents) to something
   in the database. Update the hash for any files in the database that 
   need to be checked and do not already have the hash calculated.
*/
func findMatch(filepath string, length int64) (bool, string, []byte) {
	// search the database for files with matching length
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
	if more { 
		// found similar files - need to calculate the hash of the candidate
		hashCandidate := getHash(filepath) // need hash 
		for rows.Next() {
			var possMatchLen int64	// do we really need this?
			var possMatchFilename string
			var possMatchHash []byte
			if err := rows.Scan(&possMatchLen, &possMatchFilename, &possMatchHash); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("possible: %s is %d\n", possMatchFilename, length)
			if possMatchHash == nil { //need hash for the possible match?
				possMatchHash = getHash(possMatchFilename)
				updateHash(possMatchFilename, possMatchHash)
			}
			if bytes.Compare(hashCandidate, possMatchHash) == 0 { // matching hash?
				// TODO perform byte by byte check.
				return true, possMatchFilename, possMatchHash
			}
		}
	} else {
		return false, "", nil	// no same length files
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return false, "", nil
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
		if !found {
			insertFile(path, info.Size(), hash)
		}
	}
	return nil
}

func main() {
	os.Remove(dbName) // remove pevious copy

	// open the database
	db, err = sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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
		log.Printf("%q: %s\n", err, sqlCreateStmt)
		return
	}

	filepath.Walk("./sample-files", myWalkFunc)

}
