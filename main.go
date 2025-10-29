package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
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
	r.GET("/log", showLog)

	/*r.GET("/Person/:id", showPerson)
	r.GET("/EditPerson/:id", editPerson)
	r.POST("/update_person", savePerson)
	r.GET("/SetPassword/:id", setPassword)
	r.POST("/set_password", setPassword2)
	r.GET("/DelPerson/:id", delPerson)
	r.GET("/UnlockPerson/:id", unlockPerson)*/

	// Start server, on non-default port
	fmt.Println("Running on port 8222")
	r.Run(":8222")
}
