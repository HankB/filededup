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

/*
func getLengthMatches(filename string) sql.Rows {

}
*/

var dbName = "filelist.db"
var db *sql.DB
var err error

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
		hash blob,
		links integer);
	`
	// create the table
	_, err = db.Exec(sqlCreateStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlCreateStmt)
		return
	}

	filepath.Walk("./sample-files", myWalkFunc)

}
