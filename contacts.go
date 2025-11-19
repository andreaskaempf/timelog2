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

// Page showing one contact, with all the projects linked to
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

	// Fetch linked projects for this contact
	contactProjects := getProjectsForContact(id)

	// Get all active projects that the contact is not yet linked to, for the dropdown (to link new ones)
	allProjects := getProjects()
	newProjects := []Project{}
	for _, p := range allProjects {
		if !p.Active {
			continue
		}
		already := false
		for _, q := range contactProjects {
			if p.Id == q.Id {
				already = true
				break
			}
		}
		if !already {
			newProjects = append(newProjects, p)
		}
	}

	// Show the page
	c.HTML(http.StatusOK,
		"contact.html",
		gin.H{"c": contact, "projects": contactProjects, "newProjects": newProjects, "current": "contacts"})
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

// Handle adding a project link to a contact
func addContactProjectLink(c *gin.Context) {

	// Get contact ID from URL
	contactIdStr := c.Param("contact_id")
	contactId, err := strconv.Atoi(contactIdStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid contact ID")
		return
	}

	// Get project ID from form
	projectIdStr := c.PostForm("project_id")
	projectId, err := strconv.Atoi(projectIdStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid project ID")
		return
	}

	// Add the link
	addProjectContact(projectId, contactId)

	// Redirect back to the contact page
	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/contact/%d", contactId))
}

// Handle removing a project link from a contact
func deleteContactProjectLink(c *gin.Context) {

	// Get contact ID and project ID from URL query string
	contactIdStr := c.Query("cid")
	contactId, err := strconv.Atoi(contactIdStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid contact ID")
		return
	}

	projectIdStr := c.Query("pid")
	projectId, err := strconv.Atoi(projectIdStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid project ID")
		return
	}

	// Delete the link
	deleteProjectContact(projectId, contactId)

	// Redirect back to the contact page
	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/contact/%d", contactId))
}
