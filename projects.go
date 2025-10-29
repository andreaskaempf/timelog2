package main

import (
	//"fmt"
	"net/http"
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
