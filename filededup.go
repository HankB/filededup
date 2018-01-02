package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path/filepath"
)

var dbName = "filelist.db"
var db	*sql.DB
var err error

// callback from Walk()
func myWalkFunc(path string, info os.FileInfo, err error) error {
	if info.Mode() & os.ModeType != 0 {	// not a regular file?
		print("other: ", path, "\n")
	} else {
		print("len: ", info.Size(), " ", path, "\n")
		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		stmt, err := tx.Prepare("insert into files(length, filename) values(?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		_, err = stmt.Exec(info.Size(), path)
		if err != nil {
			log.Fatal(err)
		}
		tx.Commit()
	}
	return nil
}


func main() {
	os.Remove(dbName)		// remove pevious copy

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
		hash blob );
	`
	// create the table
	_, err = db.Exec(sqlCreateStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlCreateStmt)
		return
	}

	filepath.Walk("./", myWalkFunc)



}