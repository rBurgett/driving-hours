package models

// DrivingLog maps date strings (YYYY-MM-DD) to DayEntry
type DrivingLog map[string]DayEntry

// DayEntry represents hours logged for a single day
type DayEntry struct {
	DayHours   float64 `json:"day_hours"`
	NightHours float64 `json:"night_hours"`
}

// HasEntry checks if there's an entry for the given date
func (d DrivingLog) HasEntry(date string) bool {
	entry, exists := d[date]
	if !exists {
		return false
	}
	return entry.DayHours > 0 || entry.NightHours > 0
}

// GetEntry returns the entry for a date, or zero values if not found
func (d DrivingLog) GetEntry(date string) DayEntry {
	if entry, exists := d[date]; exists {
		return entry
	}
	return DayEntry{}
}
