// Page handlers for log of activity on projects

package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// LogEntryWithSubtotals contains a work entry and subtotal information
type LogEntryWithSubtotals struct {
	Work           Work
	ShowDayTotal   bool
	DayTotal       float64
	ShowWeekTotal  bool
	WeekTotal      float64
	ShowMonthTotal bool
	MonthTotal     float64
	DayLabel       string
	WeekLabel      string
	MonthLabel     string
}

// Page showing activity on projects
func showLog(c *gin.Context) {

	// Get all work entries
	entries := getWorkEntries()

	// Process entries and calculate subtotals
	logEntries := []LogEntryWithSubtotals{}
	var lastDate, lastWeek, lastMonth string
	var dayTotal, weekTotal, monthTotal float64

	for i, w := range entries {
		entry := LogEntryWithSubtotals{Work: w}

		// Parse date (assuming format YYYY-MM-DD)
		date, err := time.Parse("2006-01-02", w.WorkDate)
		if err != nil {
			// Skip entries with invalid dates
			fmt.Printf("showLog: invalid date \"%s\"\n", w.WorkDate)
			continue
		}

		// Calculate day, week, and month labels
		dayLabel := date.Format("2006-01-02")
		year, week := date.ISOWeek()
		weekLabel := fmt.Sprintf("%d-W%02d", year, week)
		monthLabel := date.Format("2006-01")

		// Check if we're starting a new day
		if dayLabel != lastDate && lastDate != "" {
			// Show day total on last entry of previous day
			if len(logEntries) > 0 {
				logEntries[len(logEntries)-1].ShowDayTotal = true
				logEntries[len(logEntries)-1].DayTotal = dayTotal
				logEntries[len(logEntries)-1].DayLabel = lastDate
			}
			dayTotal = 0
		}
		dayTotal += w.Hours

		// Check if we're starting a new week
		if weekLabel != lastWeek && lastWeek != "" {
			// Show week total on last entry of previous week
			if len(logEntries) > 0 {
				logEntries[len(logEntries)-1].ShowWeekTotal = true
				logEntries[len(logEntries)-1].WeekTotal = weekTotal
				logEntries[len(logEntries)-1].WeekLabel = lastWeek
			}
			weekTotal = 0
		}
		weekTotal += w.Hours

		// Check if we're starting a new month
		if monthLabel != lastMonth && lastMonth != "" {
			// Show month total on last entry of previous month
			if len(logEntries) > 0 {
				logEntries[len(logEntries)-1].ShowMonthTotal = true
				logEntries[len(logEntries)-1].MonthTotal = monthTotal
				logEntries[len(logEntries)-1].MonthLabel = lastMonth
			}
			monthTotal = 0
		}
		monthTotal += w.Hours

		lastDate = dayLabel
		lastWeek = weekLabel
		lastMonth = monthLabel

		logEntries = append(logEntries, entry)

		// If this is the last entry, show all totals
		if i == len(entries)-1 {
			logEntries[len(logEntries)-1].ShowDayTotal = true
			logEntries[len(logEntries)-1].DayTotal = dayTotal
			logEntries[len(logEntries)-1].DayLabel = dayLabel

			logEntries[len(logEntries)-1].ShowWeekTotal = true
			logEntries[len(logEntries)-1].WeekTotal = weekTotal
			logEntries[len(logEntries)-1].WeekLabel = weekLabel

			logEntries[len(logEntries)-1].ShowMonthTotal = true
			logEntries[len(logEntries)-1].MonthTotal = monthTotal
			logEntries[len(logEntries)-1].MonthLabel = monthLabel
		}
	}

	// Show the page
	c.HTML(http.StatusOK,
		"log.html",
		gin.H{"entries": logEntries})
}

// Page showing one work entry detail
func showWorkEntry(c *gin.Context) {

	// Get work entry ID from URL
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid work entry ID")
		return
	}

	// Show the page
	c.HTML(http.StatusOK,
		"work_entry.html",
		gin.H{"work": getWorkEntry(id)})
}

// Page to create/edit a work entry
func editWork(c *gin.Context) {

	// ID from URL param (consistent with /edit_log/:id)
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid work entry ID")
		return
	}

	var w Work
	if id == 0 {
		// New entry defaults
		w = Work{
			Id:        0,
			WorkDate:  time.Now().Format("2006-01-02"),
			Billable:  true,
			Hours:     1.0,
			ProjectId: 0,
		}
	} else {
		w = getWorkEntry(id)
	}

	// Get active projects for dropdown
	activeProjects := []Project{}
	for _, p := range getProjects() {
		if p.Active {
			activeProjects = append(activeProjects, p)
		}
	}

	c.HTML(http.StatusOK, "edit_work.html", gin.H{
		"work":     w,
		"projects": activeProjects,
	})
}

// Handle save of a work entry
func saveWorkForm(c *gin.Context) {

	// Parse fields
	id, _ := strconv.Atoi(c.PostForm("id"))
	projectId, err := strconv.Atoi(c.PostForm("project_id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid project")
		return
	}
	workDate := c.PostForm("work_date")
	hoursStr := c.PostForm("hours")
	hours, err := strconv.ParseFloat(hoursStr, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid hours")
		return
	}
	billable := c.PostForm("billable") == "on" || c.PostForm("billable") == "true"
	description := c.PostForm("description")

	w := Work{
		Id:          id,
		ProjectId:   projectId,
		WorkDate:    workDate,
		Hours:       hours,
		Billable:    billable,
		Description: description,
	}

	savedId := saveWork(w)

	// Redirect to work entry detail
	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/work_entry/%d", savedId))
}

// Handle deletion of a work entry
func deleteWorkHandler(c *gin.Context) {
	// Get work entry ID from URL
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid work entry ID")
		return
	}
	// Delete and redirect to log
	deleteWork(id)
	c.Redirect(http.StatusSeeOther, "/log")
}
