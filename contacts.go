// Page handlers for contacts

package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Page showing contacts
func showContacts(c *gin.Context) {

	// Show the page as a table
	c.HTML(http.StatusOK,
		"contacts.html",
		gin.H{"contacts": 0})
}
