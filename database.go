// database.go
//
// Data model for the time & contacts system, including structure definitions
// for all tables, and functions to retrieve or update data in the database.
// All database functions should be in this file.
//
// TODO: Remove panics with more graceful messages?

package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// Connect to database
// Don't forget to "defer db.Close() after calling this
func dbConnect() *sql.DB {
	db, err := sql.Open("sqlite3", "./timelog.db")
	if err != nil {
		panic("dbConnect: " + err.Error())
	}
	return db
}

//------------------------------------------------------------------//
//                          P R O J E C T S                         //
//------------------------------------------------------------------//

// Record format for one project
type Project struct {
	Id          int
	Client      string
	Name        string
	Description string
	Category    string // Billable, CD, IP, Training, Absent, Other
	Active      bool
}

// Get a list of all projects, sorted by client, name
func getProjects() []Project {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	// Execute query to get all projects
	rows, err := db.Query("select id, client, name, description, category, active from project order by client, name")
	if err != nil {
		panic("getProjects query: " + err.Error())
	}
	defer rows.Close()

	// Collect into a list
	pp := []Project{}
	for rows.Next() {
		p := Project{}
		err := rows.Scan(&p.Id, &p.Client, &p.Name, &p.Description, &p.Category, &p.Active)
		if err != nil {
			panic("getProjects next: " + err.Error())
		}
		pp = append(pp, p)
	}
	if rows.Err() != nil {
		panic("getProjects exit: " + err.Error())
	}

	// Return list
	return pp
}
