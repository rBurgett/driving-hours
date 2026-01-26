package utils

import (
	"time"
)

// GetGreeting returns a time-based greeting
func GetGreeting() string {
	hour := time.Now().Hour()

	switch {
	case hour >= 5 && hour < 12:
		return "Good morning"
	case hour >= 12 && hour < 17:
		return "Good afternoon"
	case hour >= 17 && hour < 21:
		return "Good evening"
	default:
		return "Good night"
	}
}

// CalendarDay represents a day in the calendar
type CalendarDay struct {
	Day          int
	Date         string
	IsOtherMonth bool
	IsToday      bool
	HasEntry     bool
	Entry        interface{}
}

// CalendarData contains all data needed to render a calendar
type CalendarData struct {
	Year      int
	Month     int
	MonthName string
	PrevMonth int
	PrevYear  int
	NextMonth int
	NextYear  int
	Days      []CalendarDay
}

// GetCalendarData generates calendar data for the given month and year
func GetCalendarData(year, month int, hasEntry func(string) bool, getEntry func(string) interface{}) CalendarData {
	if month < 1 || month > 12 {
		now := time.Now()
		year = now.Year()
		month = int(now.Month())
	}

	firstOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	today := time.Now().Format("2006-01-02")

	// Calculate previous and next month
	prevMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local).AddDate(0, -1, 0)
	nextMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, 0)

	// Start from the Sunday of the week containing the first day
	startDay := firstOfMonth
	for startDay.Weekday() != time.Sunday {
		startDay = startDay.AddDate(0, 0, -1)
	}

	// Generate 6 weeks (42 days) to fill the calendar grid
	var days []CalendarDay
	current := startDay

	for i := 0; i < 42; i++ {
		dateStr := current.Format("2006-01-02")
		isOtherMonth := current.Month() != time.Month(month)

		day := CalendarDay{
			Day:          current.Day(),
			Date:         dateStr,
			IsOtherMonth: isOtherMonth,
			IsToday:      dateStr == today,
		}

		if !isOtherMonth && hasEntry != nil {
			day.HasEntry = hasEntry(dateStr)
			if day.HasEntry && getEntry != nil {
				day.Entry = getEntry(dateStr)
			}
		}

		days = append(days, day)
		current = current.AddDate(0, 0, 1)

		// Stop if we've covered the last day of the month and completed the week
		if current.After(lastOfMonth) && current.Weekday() == time.Sunday {
			break
		}
	}

	return CalendarData{
		Year:      year,
		Month:     month,
		MonthName: firstOfMonth.Month().String(),
		PrevMonth: int(prevMonth.Month()),
		PrevYear:  prevMonth.Year(),
		NextMonth: int(nextMonth.Month()),
		NextYear:  nextMonth.Year(),
		Days:      days,
	}
}
