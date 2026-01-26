package handlers

import (
	"net/http"
	"strconv"
	"time"

	"driving-hours/internal/auth"
	"driving-hours/internal/models"
	"driving-hours/internal/storage"
	"driving-hours/internal/templates"
	"driving-hours/internal/utils"
)

type DriverHandler struct {
	storage  storage.Storage
	renderer *templates.Renderer
}

func NewDriverHandler(s storage.Storage, r *templates.Renderer) *DriverHandler {
	return &DriverHandler{
		storage:  s,
		renderer: r,
	}
}

func (h *DriverHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r)

	// Get month/year from query params
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))

	if year == 0 {
		year = time.Now().Year()
	}
	if month == 0 {
		month = int(time.Now().Month())
	}

	// Generate calendar data
	calendar := utils.GetCalendarData(year, month,
		func(date string) bool {
			return user.DrivingLog.HasEntry(date)
		},
		func(date string) interface{} {
			return user.DrivingLog.GetEntry(date)
		},
	)

	// Check for fireworks flag from session flash
	showFireworks := r.URL.Query().Get("celebrate") == "1"

	h.renderer.Render(w, r, "driver/dashboard.html", templates.Data{
		"Title":         "Dashboard",
		"User":          user,
		"Greeting":      utils.GetGreeting(),
		"Calendar":      calendar,
		"Today":         time.Now().Format("2006-01-02"),
		"ShowFireworks": showFireworks,
	})
}

func (h *DriverHandler) LogHours(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r)

	date := r.FormValue("date")
	deleteEntry := r.FormValue("delete") == "1"
	dayHoursStr := r.FormValue("day_hours")
	dayMinutesStr := r.FormValue("day_minutes")
	nightHoursStr := r.FormValue("night_hours")
	nightMinutesStr := r.FormValue("night_minutes")

	if date == "" {
		http.Redirect(w, r, "/driver?error=date_required", http.StatusSeeOther)
		return
	}

	// Initialize driving log if nil
	if user.DrivingLog == nil {
		user.DrivingLog = make(models.DrivingLog)
	}

	// Handle delete action
	if deleteEntry {
		delete(user.DrivingLog, date)
		if err := h.storage.SaveUser(user); err != nil {
			http.Error(w, "Failed to delete entry", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/driver", http.StatusSeeOther)
		return
	}

	// Parse hours and minutes
	dayHours, _ := strconv.ParseFloat(dayHoursStr, 64)
	dayMinutes, _ := strconv.ParseFloat(dayMinutesStr, 64)
	nightHours, _ := strconv.ParseFloat(nightHoursStr, 64)
	nightMinutes, _ := strconv.ParseFloat(nightMinutesStr, 64)

	// Convert to decimal hours
	totalDayHours := dayHours + (dayMinutes / 60)
	totalNightHours := nightHours + (nightMinutes / 60)

	// Set entry (replaces existing) or delete if zero
	if totalDayHours > 0 || totalNightHours > 0 {
		user.DrivingLog[date] = models.DayEntry{
			DayHours:   totalDayHours,
			NightHours: totalNightHours,
		}
	} else {
		delete(user.DrivingLog, date)
	}

	// Save user
	if err := h.storage.SaveUser(user); err != nil {
		http.Error(w, "Failed to save hours", http.StatusInternalServerError)
		return
	}

	// Only celebrate if hours were actually logged
	if totalDayHours > 0 || totalNightHours > 0 {
		http.Redirect(w, r, "/driver?celebrate=1", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/driver", http.StatusSeeOther)
	}
}

func (h *DriverHandler) Profile(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r)

	h.renderer.Render(w, r, "driver/profile.html", templates.Data{
		"Title": "Profile",
		"User":  user,
	})
}

func (h *DriverHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r)

	name := r.FormValue("name")
	currentPassword := r.FormValue("current_password")
	newPassword := r.FormValue("new_password")

	var errors []string
	var success string

	if name == "" {
		errors = append(errors, "Name is required")
	}

	// If changing password, validate current password
	if newPassword != "" {
		if currentPassword == "" {
			errors = append(errors, "Current password is required to set a new password")
		} else {
			valid, _ := auth.VerifyPassword(currentPassword, user.PasswordHash)
			if !valid {
				errors = append(errors, "Current password is incorrect")
			}
		}
	}

	if len(errors) > 0 {
		h.renderer.Render(w, r, "driver/profile.html", templates.Data{
			"Title":  "Profile",
			"User":   user,
			"Errors": errors,
			"Name":   name,
		})
		return
	}

	user.Name = name

	if newPassword != "" {
		hash, err := auth.HashPassword(newPassword)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}
		user.PasswordHash = hash
		success = "Profile and password updated successfully"
	} else {
		success = "Profile updated successfully"
	}

	if err := h.storage.SaveUser(user); err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	h.renderer.Render(w, r, "driver/profile.html", templates.Data{
		"Title":   "Profile",
		"User":    user,
		"Success": success,
	})
}
