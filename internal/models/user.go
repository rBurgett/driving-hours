package models

import (
	"time"
)

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleDriver Role = "driver"
)

type User struct {
	ID                 string     `json:"id"`
	Email              string     `json:"email"`
	Name               string     `json:"name"`
	PasswordHash       string     `json:"password_hash"`
	Role               Role       `json:"role"`
	RequiredDayHours   float64    `json:"required_day_hours,omitempty"`
	RequiredNightHours float64    `json:"required_night_hours,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DrivingLog         DrivingLog `json:"driving_log,omitempty"`
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) IsDriver() bool {
	return u.Role == RoleDriver
}

func (u *User) TotalDayHours() float64 {
	var total float64
	for _, entry := range u.DrivingLog {
		total += entry.DayHours
	}
	return total
}

func (u *User) TotalNightHours() float64 {
	var total float64
	for _, entry := range u.DrivingLog {
		total += entry.NightHours
	}
	return total
}

func (u *User) TotalHours() float64 {
	return u.TotalDayHours() + u.TotalNightHours()
}

func (u *User) DayProgress() float64 {
	if u.RequiredDayHours <= 0 {
		return 0
	}
	progress := (u.TotalDayHours() / u.RequiredDayHours) * 100
	if progress > 100 {
		return 100
	}
	return progress
}

func (u *User) NightProgress() float64 {
	if u.RequiredNightHours <= 0 {
		return 0
	}
	progress := (u.TotalNightHours() / u.RequiredNightHours) * 100
	if progress > 100 {
		return 100
	}
	return progress
}

func (u *User) WeeklyAverage() float64 {
	now := time.Now()
	cutoff := now.AddDate(0, 0, -28)

	var total float64
	for dateStr, entry := range u.DrivingLog {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}
		if date.After(cutoff) && !date.After(now) {
			total += entry.DayHours + entry.NightHours
		}
	}
	return total / 4
}
