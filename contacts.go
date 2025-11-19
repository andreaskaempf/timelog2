// Contacts pages

package main

import (
	"fmt"
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
		gin.H{"contacts": allContacts, "current": "contacts"})
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
		gin.H{"c": contact, "current": "contacts"})
}

// Page to edit a contact (or create new one if id is 0)
func editContact(c *gin.Context) {

	// Get contact ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid contact ID")
		return
	}

	var cont Contact // new contact is blank by default
	if id > 0 {      // Existing contact - get from database
		cont = getContact(id)
	}

	// Show the edit page
	c.HTML(http.StatusOK,
		"edit_contact.html",
		gin.H{"c": cont, "current": "contacts"})
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
	cont := Contact{
		Id:        id,
		FirstName: c.PostForm("first_name"),
		LastName:  c.PostForm("last_name"),
		Company:   c.PostForm("company"),
		Title:     c.PostForm("title"),
		Source:    c.PostForm("source"),
		Phones:    c.PostForm("phones"),
		Emails:    c.PostForm("emails"),
		Address:   c.PostForm("address"),
		Comments:  c.PostForm("comments"),
		Active:    c.PostForm("active") == "on" || c.PostForm("active") == "true",
	}

	// Save the contact (TODO: is ID assigned for new contacts?)
	savedId := saveContact(cont)

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

	// Delete the contact (TODO: all child records)
	deleteContact(id)

	// Redirect to contacts list
	c.Redirect(http.StatusSeeOther, "/contacts")
}
