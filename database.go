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
	"fmt"

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
//                    G E N E R A L   U T I L I T I E S             //
//------------------------------------------------------------------//

// Get the maximum ID from a table
// Returns 0 if the table is empty or has no rows
// Note: tableName should be validated by the caller to prevent SQL injection
func getMaxId(tableName string) int {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	// Execute query to get maximum ID
	var maxId sql.NullInt64
	query := fmt.Sprintf("select max(id) from %s", tableName)
	err := db.QueryRow(query).Scan(&maxId)
	if err != nil {
		panic("getMaxId: " + err.Error())
	}

	// Return 0 if NULL (empty table), otherwise return the max ID
	if !maxId.Valid {
		return 0
	}
	return int(maxId.Int64)
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

// Get one project by ID
func getProject(id int) Project {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	// Execute query to get one project
	var p Project
	err := db.QueryRow("select id, client, name, description, category, active from project where id = ?", id).
		Scan(&p.Id, &p.Client, &p.Name, &p.Description, &p.Category, &p.Active)
	if err != nil {
		if err == sql.ErrNoRows {
			panic("getProject: project with id " + fmt.Sprintf("%d", id) + " not found")
		}
		panic("getProject: " + err.Error())
	}

	// Return project
	return p
}

// Save a project (insert if Id is zero, update if Id is nonzero)
// Returns the project ID
func saveProject(p Project) int {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	if p.Id == 0 {
		// Get next ID
		nextId := getMaxId("project") + 1
		p.Id = nextId

		// Insert new project
		_, err := db.Exec("insert into project (id, client, name, description, category, active) values (?, ?, ?, ?, ?, ?)",
			p.Id, p.Client, p.Name, p.Description, p.Category, p.Active)
		if err != nil {
			panic("saveProject insert: " + err.Error())
		}
	} else {
		// Update existing project
		_, err := db.Exec("update project set client=?, name=?, description=?, category=?, active=? where id=?",
			p.Client, p.Name, p.Description, p.Category, p.Active, p.Id)
		if err != nil {
			panic("saveProject update: " + err.Error())
		}
	}
	return p.Id
}
