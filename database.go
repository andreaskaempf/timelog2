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
	"strconv"

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

// Delete a project and all its child records (work and project_contact)
func deleteProject(id int) {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		panic("deleteProject begin: " + err.Error())
	}

	// Delete all work records for this project
	_, err = tx.Exec("delete from work where project_id = ?", id)
	if err != nil {
		tx.Rollback()
		panic("deleteProject work: " + err.Error())
	}

	// Delete all project_contact records for this project
	_, err = tx.Exec("delete from project_contact where project_id = ?", id)
	if err != nil {
		tx.Rollback()
		panic("deleteProject project_contact: " + err.Error())
	}

	// Delete the project itself
	_, err = tx.Exec("delete from project where id = ?", id)
	if err != nil {
		tx.Rollback()
		panic("deleteProject project: " + err.Error())
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		panic("deleteProject commit: " + err.Error())
	}
}

//------------------------------------------------------------------//
//                            W O R K                                //
//------------------------------------------------------------------//

// Record format for one work entry
type Work struct {
	Id          int
	ProjectId   int
	WorkDate    string // date as string
	Hours       float64
	Billable    bool
	Description string
	// Joined fields from project
	ProjectName string
	Client      string
}

// Cutoff year for work entries
const cutoffDate = "2025-01-01"

// Get all work entries, sorted by date (increasing)
func getWorkEntries() []Work {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	// Execute query to get all work entries with project info
	query := `select w.id, w.project_id, w.work_date, w.hours, w.billable, w.description,
	          p.name as project_name, p.client
	          from work w
	          left join project p on w.project_id = p.id
	          where w.work_date >= ?
	          order by w.work_date, w.id`
	rows, err := db.Query(query, cutoffDate)
	if err != nil {
		panic("getWorkEntries query: " + err.Error())
	}
	defer rows.Close()

	// Collect into a list
	ww := []Work{}
	for rows.Next() {
		w := Work{}
		var hrs, billable string
		err := rows.Scan(&w.Id, &w.ProjectId, &w.WorkDate, &hrs, &billable, &w.Description,
			&w.ProjectName, &w.Client)
		if err != nil {
			panic("getWorkEntries next: " + err.Error())
		}

		// Convert some fields
		if len(w.WorkDate) > 10 { // fix dates that include time (e.g., "2025-01-01T12:00:00")
			w.WorkDate = w.WorkDate[:10]
		}
		w.Hours, err = strconv.ParseFloat(hrs, 64)
		if err != nil {
			fmt.Printf("getWorkEntries: invalid hours \"%s\"\n", hrs)
			w.Hours = 0
		}
		w.Billable = billable == "1"

		// Add to list
		ww = append(ww, w)
	}
	if rows.Err() != nil {
		panic("getWorkEntries exit: " + err.Error())
	}

	// Return list
	fmt.Printf("getWorkEntries: %d rows\n", len(ww))
	return ww
}

// Get one work entry by ID
func getWorkEntry(id int) Work {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	// Execute query to get one work entry with project info
	query := `select w.id, w.project_id, w.work_date, w.hours, w.billable, w.description,
	          p.name as project_name, p.client
	          from work w
	          left join project p on w.project_id = p.id
	          where w.id = ?`
	var w Work
	var workDate sql.NullString
	var hours sql.NullFloat64
	var billable sql.NullBool
	var description sql.NullString
	var projectName sql.NullString
	var client sql.NullString

	err := db.QueryRow(query, id).Scan(&w.Id, &w.ProjectId, &workDate, &hours, &billable, &description,
		&projectName, &client)
	if err != nil {
		if err == sql.ErrNoRows {
			panic("getWorkEntry: work entry with id " + fmt.Sprintf("%d", id) + " not found")
		}
		panic("getWorkEntry: " + err.Error())
	}

	if workDate.Valid {
		w.WorkDate = workDate.String
	}
	if hours.Valid {
		w.Hours = hours.Float64
	}
	if billable.Valid {
		w.Billable = billable.Bool
	}
	if description.Valid {
		w.Description = description.String
	}
	if projectName.Valid {
		w.ProjectName = projectName.String
	}
	if client.Valid {
		w.Client = client.String
	}

	// Return work entry
	return w
}

// Get work entries between dates [startDate, endDate] inclusive, sorted by date
func getWorkEntriesBetween(startDate, endDate string) []Work {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	// Query with project info
	query := `select w.id, w.project_id, w.work_date, w.hours, w.billable, w.description,
	          p.name as project_name, p.client
	          from work w
	          left join project p on w.project_id = p.id
	          where w.work_date >= ? and w.work_date <= ?
	          order by w.work_date, w.id`
	rows, err := db.Query(query, startDate, endDate)
	if err != nil {
		panic("getWorkEntriesBetween query: " + err.Error())
	}
	defer rows.Close()

	entries := []Work{}
	for rows.Next() {
		w := Work{}
		var hrs, billable string
		err := rows.Scan(&w.Id, &w.ProjectId, &w.WorkDate, &hrs, &billable, &w.Description, &w.ProjectName, &w.Client)
		if err != nil {
			panic("getWorkEntriesBetween next: " + err.Error())
		}
		if len(w.WorkDate) > 10 {
			w.WorkDate = w.WorkDate[:10]
		}
		w.Hours, err = strconv.ParseFloat(hrs, 64)
		if err != nil {
			w.Hours = 0
		}
		w.Billable = billable == "1"
		entries = append(entries, w)
	}
	if rows.Err() != nil {
		panic("getWorkEntriesBetween exit: " + err.Error())
	}
	return entries
}

// Delete one work entry by ID
func deleteWork(id int) {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	_, err := db.Exec("delete from work where id = ?", id)
	if err != nil {
		panic("deleteWork: " + err.Error())
	}
}

// Save a work entry (insert if Id is zero, update if Id is nonzero)
// Returns the work ID
func saveWork(w Work) int {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	if w.Id == 0 {
		// Get next ID
		nextId := getMaxId("work") + 1
		w.Id = nextId

		// Insert new work entry
		_, err := db.Exec("insert into work (id, project_id, work_date, hours, billable, description) values (?, ?, ?, ?, ?, ?)",
			w.Id, w.ProjectId, w.WorkDate, w.Hours, w.Billable, w.Description)
		if err != nil {
			panic("saveWork insert: " + err.Error())
		}
	} else {
		// Update existing work entry
		_, err := db.Exec("update work set project_id=?, work_date=?, hours=?, billable=?, description=? where id=?",
			w.ProjectId, w.WorkDate, w.Hours, w.Billable, w.Description, w.Id)
		if err != nil {
			panic("saveWork update: " + err.Error())
		}
	}
	return w.Id
}
