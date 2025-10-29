# Time & Contacts System

This will be a Go web app to manage projects, time on projects, and contacts.
Database schema is based on the old 
[Python version](https://github.com/andreaskaempf/timelog), but completely
rewritten in Go.

To install
* Clone this repository
* `cd timelog2`
* If you have an existing database from the old timelog system, make a
  couple of changes to bring to new schema, otherwise create a new timelog.db
  using schema.txt (TODO: instructions)
* `go get` to install dependencies
* Download [Bulma](https://bulma.io) and install it into the static directory
* `go build` to build executable
* `./timelog2` to start the app server
* Browse to http://localhost:8222

AK, Oct-Nov 2025
