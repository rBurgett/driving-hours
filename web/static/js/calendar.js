// Calendar interaction for driver dashboard
document.addEventListener('DOMContentLoaded', function() {
    const calendarDays = document.querySelectorAll('.calendar-day:not(.other-month)');
    const dateInput = document.getElementById('date');
    const dayHoursInput = document.getElementById('day_hours');
    const dayMinutesInput = document.getElementById('day_minutes');
    const nightHoursInput = document.getElementById('night_hours');
    const nightMinutesInput = document.getElementById('night_minutes');
    const form = document.getElementById('log-form');

    if (!dateInput || !form) return;

    calendarDays.forEach(day => {
        day.addEventListener('click', function() {
            const date = this.dataset.date;
            if (!date) return;

            // Set the date
            dateInput.value = date;

            // If the day has an entry, parse the title for existing values
            if (this.classList.contains('has-entry') && this.title) {
                const match = this.title.match(/Day: ([\d.]+)h, Night: ([\d.]+)h/);
                if (match) {
                    const dayHours = parseFloat(match[1]);
                    const nightHours = parseFloat(match[2]);

                    // Convert to hours and minutes
                    const dayH = Math.floor(dayHours);
                    const dayM = Math.round((dayHours - dayH) * 60);
                    const nightH = Math.floor(nightHours);
                    const nightM = Math.round((nightHours - nightH) * 60);

                    if (dayHoursInput) dayHoursInput.value = dayH;
                    if (dayMinutesInput) dayMinutesInput.value = dayM;
                    if (nightHoursInput) nightHoursInput.value = nightH;
                    if (nightMinutesInput) nightMinutesInput.value = nightM;
                }
            } else {
                // Reset to zero for new entries
                if (dayHoursInput) dayHoursInput.value = 0;
                if (dayMinutesInput) dayMinutesInput.value = 0;
                if (nightHoursInput) nightHoursInput.value = 0;
                if (nightMinutesInput) nightMinutesInput.value = 0;
            }

            // Scroll to the form
            form.scrollIntoView({ behavior: 'smooth', block: 'start' });
        });
    });
});
