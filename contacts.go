// Contacts pages

package main

import (
	//"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Page showing list of all contacts
func showContacts(c *gin.Context) {

	// Make sure logged in, and check if administrator
	/*sess := loadSession(c)
	if sess == nil || sess["user_id"] == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/Login")
	}
	admin := sess["admin"] == "Y"
	user := sess["user_name"]
	uid := parseInt(sess["user_id"])*/

	// Get filter from query string (all, active, inactive)
	filter := c.Query("filter")
	if filter == "" {
		filter = "active"
	}

	// Get all contacts
	allContacts := getContacts()

	// Show the page as a table
	c.HTML(http.StatusOK,
		"contacts.html",
		gin.H{"contacts": allContacts})
}

// Page showing one contact
func showContact(c *gin.Context) {

	// Get contact ID from URL
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid contact ID")
		return
	}

	// Fetch contact
	contact := getContact(id)

	// Show the page
	c.HTML(http.StatusOK,
		"contact.html",
		gin.H{"contact": contact})
}

/*
// Page to edit a contact (or create new one if id is 0)
func editContact(c *gin.Context) {

	// Get contact ID from query string
	idStr := c.Query("id")
	if idStr == "" {
		idStr = "0"
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid contact ID")
		return
	}

	var p Contact
	if id == 0 {
		// New contact - create empty contact
		p = Contact{Id: 0, Client: "", Name: "", Description: "", Category: "", Active: true}
	} else {
		// Existing contact - get from database
		p = getContact(id)
	}

	// Show the edit page
	c.HTML(http.StatusOK,
		"edit_contact.html",
		gin.H{"contact": p})
}

// Handle form submission to save a contact
func saveContactForm(c *gin.Context) {

	// Get contact ID from form
	idStr := c.PostForm("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid contact ID")
		return
	}

	// Get form values
	p := Contact{
		Id:          id,
		Client:      c.PostForm("client"),
		Name:        c.PostForm("name"),
		Description: c.PostForm("description"),
		Category:    c.PostForm("category"),
		Active:      c.PostForm("active") == "on" || c.PostForm("active") == "true",
	}

	// Save the contact
	savedId := saveContact(p)

	// Redirect to the contact page
	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/contact/%d", savedId))
}

// Handle deletion of a contact
func deleteContactHandler(c *gin.Context) {

	// Get contact ID from URL
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid contact ID")
		return
	}

	// Delete the contact (and all child records)
	deleteContact(id)

	// Redirect to contacts list
	c.Redirect(http.StatusSeeOther, "/contacts")
}
*/
