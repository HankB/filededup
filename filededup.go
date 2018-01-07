package main

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

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

func findMatch(filepath string, length int) (bool, string) {
	rows, err := db.Query(`SELECT length, filename, hash linkCount
							FROM files
							WHERE length=?`,
		length)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	more := rows.Next()
	if more { // did we get any results?
		hashCandidate = getHash(filepath) // need hash 
		for rows.Next() {
			var possiblMatch string
			if err := rows.Scan(&possiblMatch); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s is %d\n", possiblMatch, length)
		}
	} else {
		// no same length files
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return false, ""
}

// callback from Walk()
func myWalkFunc(path string, info os.FileInfo, err error) error {
	if info.Mode()&os.ModeType != 0 { // not a regular file?
		fmt.Printf("other: %s\n", path)
	} else {
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
