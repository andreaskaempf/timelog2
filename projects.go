package main

import (
	//"fmt"
	"net/http"
	"strconv"

	//"sort"
	//"strings"

	"github.com/gin-gonic/gin"
)

// Page showing list of all projects
func showProjects(c *gin.Context) {

	// Make sure logged in, and check if administrator
	/*sess := loadSession(c)
	if sess == nil || sess["user_id"] == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/Login")
	}
	admin := sess["admin"] == "Y"
	user := sess["user_name"]
	uid := parseInt(sess["user_id"])*/

	// Show the page as a table
	c.HTML(http.StatusOK,
		"projects.html",
		gin.H{"projects": getProjects()})
}

// Page showing one project
func showProject(c *gin.Context) {

	// Get project ID from URL
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid project ID")
		return
	}

	// Show the page
	c.HTML(http.StatusOK,
		"project.html",
		gin.H{"project": getProject(id)})
}
