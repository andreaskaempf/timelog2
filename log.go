// Page handlers for log of activity on projects

package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CalendarDay represents one day and its work entries
type CalendarDay struct {
	Date    string
	Day     int
	Entries []Work
}

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

// Page: monthly calendar of work entries
func showCalendar(c *gin.Context) {
	// Parse year and month from query; default to current
	now := time.Now()
	year, _ := strconv.Atoi(c.DefaultQuery("year", fmt.Sprintf("%04d", now.Year())))
	monthInt, _ := strconv.Atoi(c.DefaultQuery("month", fmt.Sprintf("%02d", int(now.Month()))))
	if monthInt < 1 || monthInt > 12 {
		monthInt = int(now.Month())
	}
	month := time.Month(monthInt)

	firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	// Compute last day
	firstNextMonth := firstOfMonth.AddDate(0, 1, 0)
	lastOfMonth := firstNextMonth.AddDate(0, 0, -1)
	daysInMonth := lastOfMonth.Day()

	// Range strings
	startDate := firstOfMonth.Format("2006-01-02")
	endDate := lastOfMonth.Format("2006-01-02")

	// Fetch entries in range
	entries := getWorkEntriesBetween(startDate, endDate)

	// Bucket entries by date
	dayMap := map[string][]Work{}
	for _, w := range entries {
		dayMap[w.WorkDate] = append(dayMap[w.WorkDate], w)
	}

	// Prepare days slice
	days := make([]CalendarDay, 0, daysInMonth)
	for d := 1; d <= daysInMonth; d++ {
		cur := time.Date(year, month, d, 0, 0, 0, 0, time.Local)
		ds := cur.Format("2006-01-02")
		days = append(days, CalendarDay{
			Date:    ds,
			Day:     d,
			Entries: dayMap[ds],
		})
	}

	// Weekday of first (0=Sunday .. 6=Saturday)
	startWeekday := int(firstOfMonth.Weekday())

	// Next/prev month routing
	prev := firstOfMonth.AddDate(0, -1, 0)
	next := firstOfMonth.AddDate(0, 1, 0)

	// Build 6x7 grid of weeks with optional days
	weeks := make([][]*CalendarDay, 0, 6)
	curWeek := make([]*CalendarDay, 7)
	// Fill leading blanks
	for i := 0; i < startWeekday; i++ {
		curWeek[i] = nil
	}
	col := startWeekday
	for i := 0; i < len(days); i++ {
		d := days[i]
		curWeek[col] = &d
		col++
		if col == 7 {
			weeks = append(weeks, curWeek)
			curWeek = make([]*CalendarDay, 7)
			col = 0
		}
	}
	if col != 0 {
		weeks = append(weeks, curWeek)
	}

	// Build project color map (stable palette)
	palette := []string{
		"is-primary", "is-link", "is-info", "is-success", "is-warning", "is-danger",
		"is-dark", "is-black", "is-primary is-light", "is-link is-light", "is-info is-light",
		"is-success is-light", "is-warning is-light", "is-danger is-light",
	}
	colors := map[int]string{}
	for _, w := range entries {
		if _, ok := colors[w.ProjectId]; !ok {
			idx := w.ProjectId % len(palette)
			if idx < 0 {
				idx = -idx
			}
			colors[w.ProjectId] = palette[idx]
		}
	}

	c.HTML(http.StatusOK, "calendar.html", gin.H{
		"year":         year,
		"month":        int(month),
		"monthName":    firstOfMonth.Format("January 2006"),
		"daysInMonth":  daysInMonth,
		"startWeekday": startWeekday,
		"days":         days,
		"weeks":        weeks,
		"prevYear":     prev.Year(),
		"prevMonth":    int(prev.Month()),
		"nextYear":     next.Year(),
		"nextMonth":    int(next.Month()),
		"colors":       colors,
	})
}
