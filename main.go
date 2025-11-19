package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	for _, p := range getProjects() {
		fmt.Printf("%+v\n", p)
	}

	// Make sure sessions directory exists under /tmp, create if necessary
	/*if !checkSessionsDir() {
		panic("No /tmp/sessions directory, could not create")
	}*/

	// Get host name, set to to dev mode if workstation
	host, err := os.Hostname()
	if err != nil {
		panic("Error getting host name: " + err.Error())
	}
	devMode := host == "brix"

	// Set to release mode if not on workstation
	if !devMode {
		fmt.Println("Host name is", host, "- running in production mode")
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router, initialize templates and location of static files
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")
	r.StaticFile("/favicon.ico", "static/favicon.ico")
	r.StaticFile("/robots.txt", "static/robots.txt")

	// Page routing for login/logout
	/*r.GET("/Login", showLogin)
	r.POST("/Login", doLogin)
	r.GET("/Logout", logout)*/

	// Project pages
	r.GET("/", showProjects)
	r.GET("/projects", showProjects)
	r.GET("/project/:id", showProject)
	r.GET("/edit_project/:id", editProject)
	r.POST("/save_project", saveProjectForm)
	r.GET("/delete_project/:id", deleteProjectHandler)

	// Work history
	r.GET("/log", showLog)
	r.GET("/edit_log/:id", editWork)
	r.POST("/save_work", saveWorkForm)
	r.GET("/work_entry/:id", showWorkEntry)
	r.GET("/delete_work/:id", deleteWorkHandler)

	// Contacts
	r.GET("/contacts", showContacts)
	r.GET("/contact/:id", showContact)
	r.GET("/edit_contact/:id", editContact)
	r.POST("/save_contact", saveContactForm)
	r.GET("/delete_contact/:id", deleteContactHandler)

	// Other pages
	r.GET("/reports", showReports)
	r.GET("/calendar", showCalendar)

	// Start server, on non-default port
	fmt.Println("Running on port 8222")
	r.Run(":8222")
}
