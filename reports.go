// Page handlers for log of activity on projects

package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Page showing reports menu
func showReports(c *gin.Context) {
	c.HTML(http.StatusOK, "reports.html", gin.H{"log": 0})
}
