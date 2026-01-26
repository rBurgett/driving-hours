package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"driving-hours/internal/auth"
	"driving-hours/internal/models"
	"driving-hours/internal/storage"
	"driving-hours/internal/templates"
)

type AdminHandler struct {
	storage  storage.Storage
	sessions *auth.SessionManager
	renderer *templates.Renderer
}

func NewAdminHandler(s storage.Storage, sm *auth.SessionManager, r *templates.Renderer) *AdminHandler {
	return &AdminHandler{
		storage:  s,
		sessions: sm,
		renderer: r,
	}
}

func (h *AdminHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r)

	drivers, err := h.storage.GetDrivers()
	if err != nil {
		http.Error(w, "Failed to load drivers", http.StatusInternalServerError)
		return
	}

	h.renderer.Render(w, r, "admin/dashboard.html", templates.Data{
		"Title":   "Admin Dashboard",
		"User":    user,
		"Drivers": drivers,
	})
}

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r)

	users, err := h.storage.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to load users", http.StatusInternalServerError)
		return
	}

	h.renderer.Render(w, r, "admin/users.html", templates.Data{
		"Title": "Manage Users",
		"User":  user,
		"Users": users,
	})
}

func (h *AdminHandler) NewUserForm(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r)

	h.renderer.Render(w, r, "admin/user_form.html", templates.Data{
		"Title":            "Create User",
		"User":             user,
		"IsNew":            true,
		"CanChangePassword": true,
		"EditUser":         &models.User{Role: models.RoleDriver},
	})
}

func (h *AdminHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r)

	email := r.FormValue("email")
	name := r.FormValue("name")
	password := r.FormValue("password")
	roleStr := r.FormValue("role")
	dayHoursStr := r.FormValue("required_day_hours")
	nightHoursStr := r.FormValue("required_night_hours")

	// Parse and validate role
	role := models.Role(roleStr)
	if role != models.RoleAdmin && role != models.RoleDriver {
		role = models.RoleDriver
	}

	// Validation
	var errors []string
	if email == "" {
		errors = append(errors, "Email is required")
	}
	if name == "" {
		errors = append(errors, "Name is required")
	}
	if password == "" {
		errors = append(errors, "Password is required")
	}

	dayHours, _ := strconv.ParseFloat(dayHoursStr, 64)
	nightHours, _ := strconv.ParseFloat(nightHoursStr, 64)

	if len(errors) > 0 {
		h.renderer.Render(w, r, "admin/user_form.html", templates.Data{
			"Title":             "Create User",
			"User":              user,
			"IsNew":             true,
			"CanChangePassword": true,
			"Errors":            errors,
			"EditUser": &models.User{
				Email:              email,
				Name:               name,
				Role:               role,
				RequiredDayHours:   dayHours,
				RequiredNightHours: nightHours,
			},
		})
		return
	}

	// Check if email already exists
	existing, _ := h.storage.GetUserByEmail(email)
	if existing != nil {
		h.renderer.Render(w, r, "admin/user_form.html", templates.Data{
			"Title":             "Create User",
			"User":              user,
			"IsNew":             true,
			"CanChangePassword": true,
			"Errors":            []string{"Email already in use"},
			"EditUser": &models.User{
				Email:              email,
				Name:               name,
				Role:               role,
				RequiredDayHours:   dayHours,
				RequiredNightHours: nightHours,
			},
		})
		return
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	now := time.Now()
	newUser := &models.User{
		ID:                 uuid.New().String(),
		Email:              email,
		Name:               name,
		PasswordHash:       hash,
		Role:               role,
		RequiredDayHours:   dayHours,
		RequiredNightHours: nightHours,
		CreatedAt:          now,
		UpdatedAt:          now,
		DrivingLog:         make(models.DrivingLog),
	}

	if err := h.storage.SaveUser(newUser); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

func (h *AdminHandler) ViewDriver(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r)
	driverID := chi.URLParam(r, "id")

	driver, err := h.storage.GetUser(driverID)
	if err != nil || driver == nil {
		http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
		return
	}

	h.renderer.Render(w, r, "admin/driver_stats.html", templates.Data{
		"Title":  driver.Name + " - Statistics",
		"User":   user,
		"Driver": driver,
	})
}

func (h *AdminHandler) EditUserForm(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r)
	editUserID := chi.URLParam(r, "id")

	editUser, err := h.storage.GetUser(editUserID)
	if err != nil || editUser == nil {
		http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
		return
	}

	// Admins cannot change another admin's password
	canChangePassword := !editUser.IsAdmin()

	h.renderer.Render(w, r, "admin/user_form.html", templates.Data{
		"Title":             "Edit " + editUser.Name,
		"User":              user,
		"IsNew":             false,
		"CanChangePassword": canChangePassword,
		"EditUser":          editUser,
	})
}

func (h *AdminHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r)
	editUserID := chi.URLParam(r, "id")

	editUser, err := h.storage.GetUser(editUserID)
	if err != nil || editUser == nil {
		http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
		return
	}

	// Admins cannot change another admin's password
	canChangePassword := !editUser.IsAdmin()

	email := r.FormValue("email")
	name := r.FormValue("name")
	password := r.FormValue("password")
	dayHoursStr := r.FormValue("required_day_hours")
	nightHoursStr := r.FormValue("required_night_hours")

	// Validation
	var errors []string
	if email == "" {
		errors = append(errors, "Email is required")
	}
	if name == "" {
		errors = append(errors, "Name is required")
	}

	dayHours, _ := strconv.ParseFloat(dayHoursStr, 64)
	nightHours, _ := strconv.ParseFloat(nightHoursStr, 64)

	if len(errors) > 0 {
		editUser.Email = email
		editUser.Name = name
		editUser.RequiredDayHours = dayHours
		editUser.RequiredNightHours = nightHours

		h.renderer.Render(w, r, "admin/user_form.html", templates.Data{
			"Title":             "Edit " + editUser.Name,
			"User":              user,
			"IsNew":             false,
			"CanChangePassword": canChangePassword,
			"Errors":            errors,
			"EditUser":          editUser,
		})
		return
	}

	// Check if email is taken by another user
	existing, _ := h.storage.GetUserByEmail(email)
	if existing != nil && existing.ID != editUser.ID {
		h.renderer.Render(w, r, "admin/user_form.html", templates.Data{
			"Title":             "Edit " + editUser.Name,
			"User":              user,
			"IsNew":             false,
			"CanChangePassword": canChangePassword,
			"Errors":            []string{"Email already in use"},
			"EditUser":          editUser,
		})
		return
	}

	editUser.Email = email
	editUser.Name = name
	editUser.RequiredDayHours = dayHours
	editUser.RequiredNightHours = nightHours

	// Update password if provided (only for drivers, not other admins)
	if password != "" && canChangePassword {
		hash, err := auth.HashPassword(password)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}
		editUser.PasswordHash = hash
	}

	if err := h.storage.SaveUser(editUser); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

func (h *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r)
	deleteUserID := chi.URLParam(r, "id")
	deletingSelf := deleteUserID == user.ID

	deleteUser, err := h.storage.GetUser(deleteUserID)
	if err != nil || deleteUser == nil {
		http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
		return
	}

	// Prevent deleting self if last admin
	if deletingSelf && user.IsAdmin() {
		users, err := h.storage.GetAllUsers()
		if err != nil {
			http.Error(w, "Failed to check admin count", http.StatusInternalServerError)
			return
		}
		adminCount := 0
		for _, u := range users {
			if u.IsAdmin() {
				adminCount++
			}
		}
		if adminCount <= 1 {
			http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
			return
		}
	}

	if err := h.storage.DeleteUser(deleteUserID); err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	// Log out if deleting self
	if deletingSelf {
		_ = h.sessions.DestroySession(w, r)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

func (h *AdminHandler) EditHoursForm(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r)
	driverID := chi.URLParam(r, "id")

	driver, err := h.storage.GetUser(driverID)
	if err != nil || driver == nil {
		http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
		return
	}

	h.renderer.Render(w, r, "admin/driver_hours.html", templates.Data{
		"Title":  driver.Name + " - Edit Hours",
		"User":   user,
		"Driver": driver,
	})
}

func (h *AdminHandler) UpdateHours(w http.ResponseWriter, r *http.Request) {
	driverID := chi.URLParam(r, "id")

	driver, err := h.storage.GetUser(driverID)
	if err != nil || driver == nil {
		http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
		return
	}

	date := r.FormValue("date")
	dayHoursStr := r.FormValue("day_hours")
	dayMinutesStr := r.FormValue("day_minutes")
	nightHoursStr := r.FormValue("night_hours")
	nightMinutesStr := r.FormValue("night_minutes")
	deleteEntry := r.FormValue("delete") == "1"

	if date == "" {
		http.Redirect(w, r, "/admin/users/"+driverID+"/hours", http.StatusSeeOther)
		return
	}

	if deleteEntry {
		delete(driver.DrivingLog, date)
	} else {
		dayHours, _ := strconv.ParseFloat(dayHoursStr, 64)
		dayMinutes, _ := strconv.ParseFloat(dayMinutesStr, 64)
		nightHours, _ := strconv.ParseFloat(nightHoursStr, 64)
		nightMinutes, _ := strconv.ParseFloat(nightMinutesStr, 64)

		totalDayHours := dayHours + (dayMinutes / 60)
		totalNightHours := nightHours + (nightMinutes / 60)

		if driver.DrivingLog == nil {
			driver.DrivingLog = make(models.DrivingLog)
		}

		if totalDayHours > 0 || totalNightHours > 0 {
			driver.DrivingLog[date] = models.DayEntry{
				DayHours:   totalDayHours,
				NightHours: totalNightHours,
			}
		} else {
			delete(driver.DrivingLog, date)
		}
	}

	if err := h.storage.SaveUser(driver); err != nil {
		http.Error(w, "Failed to update hours", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/users/"+driverID+"/hours", http.StatusSeeOther)
}

func (h *AdminHandler) Profile(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r)

	h.renderer.Render(w, r, "admin/profile.html", templates.Data{
		"Title": "Profile",
		"User":  user,
	})
}

func (h *AdminHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
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
		h.renderer.Render(w, r, "admin/profile.html", templates.Data{
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

	if err := h.storage.SaveAdmin(user); err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	h.renderer.Render(w, r, "admin/profile.html", templates.Data{
		"Title":   "Profile",
		"User":    user,
		"Success": success,
	})
}
