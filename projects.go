package main

import (
	"fmt"
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

	// Get filter from query string (all, active, inactive)
	filter := c.Query("filter")
	if filter == "" {
		filter = "active"
	}

	// Get all projects
	allProjects := getProjects()

	// Filter projects based on filter parameter
	var filteredProjects []Project
	switch filter {
	case "active":
		for _, p := range allProjects {
			if p.Active {
				filteredProjects = append(filteredProjects, p)
			}
		}
	case "inactive":
		for _, p := range allProjects {
			if !p.Active {
				filteredProjects = append(filteredProjects, p)
			}
		}
	default: // "all"
		filteredProjects = allProjects
	}

	// Show the page as a table
	c.HTML(http.StatusOK,
		"projects.html",
		gin.H{
			"projects": filteredProjects,
			"filter":   filter,
			"current":  "projects",
		})
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

	// Fetch project and related work entries
	project := getProject(id)
	entries := getWorkEntriesForProject(id)

	totalHours := 0.0
	for _, e := range entries {
		totalHours += e.Hours
	}

	// Show the page
	c.HTML(http.StatusOK,
		"project.html",
		gin.H{
			"p":          project,
			"entries":    entries,
			"totalCount": len(entries),
			"totalHours": totalHours,
			"current":    "projects",
		})
}

// Page to edit a project (or create new one if id is 0)
func editProject(c *gin.Context) {

	// Get project ID
	idStr := c.Param("id")
	if idStr == "" {
		idStr = "0"
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid project ID")
		return
	}

	var p Project
	if id == 0 {
		// New project - create empty project
		p = Project{Id: 0, Client: "", Name: "", Description: "", Category: "", Active: true}
	} else {
		// Existing project - get from database
		p = getProject(id)
	}

	// Show the edit page
	c.HTML(http.StatusOK,
		"edit_project.html",
		gin.H{"project": p, "current": "projects"})
}

// Handle form submission to save a project
func saveProjectForm(c *gin.Context) {

	// Get project ID from form
	idStr := c.PostForm("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid project ID")
		return
	}

	// Get form values
	p := Project{
		Id:          id,
		Client:      c.PostForm("client"),
		Name:        c.PostForm("name"),
		Description: c.PostForm("description"),
		Category:    c.PostForm("category"),
		Active:      c.PostForm("active") == "on" || c.PostForm("active") == "true",
	}

	// Save the project
	savedId := saveProject(p)

	// Redirect to the project page
	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/project/%d", savedId))
}

// Handle deletion of a project
func deleteProjectHandler(c *gin.Context) {

	// Get project ID from URL
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid project ID")
		return
	}

	// Delete the project (and all child records)
	deleteProject(id)

	// Redirect to projects list
	c.Redirect(http.StatusSeeOther, "/projects")
}
