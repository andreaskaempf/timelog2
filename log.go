// Page handlers for log of activity on projects

package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Page showing activity on projects
func showLog(c *gin.Context) {

	// Show the page as a table
	c.HTML(http.StatusOK,
		"log.html",
		gin.H{"log": 0})
}
