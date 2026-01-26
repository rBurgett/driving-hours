// Calendar and list view interaction for driver dashboard
document.addEventListener('DOMContentLoaded', function() {
    const calendarDays = document.querySelectorAll('.calendar-day:not(.other-month)');
    const dateInput = document.getElementById('date');
    const dayHoursInput = document.getElementById('day_hours');
    const dayMinutesInput = document.getElementById('day_minutes');
    const nightHoursInput = document.getElementById('night_hours');
    const nightMinutesInput = document.getElementById('night_minutes');
    const form = document.getElementById('log-form');
    const deleteBtn = document.getElementById('delete-btn');

    if (!dateInput || !form) return;

    // Tab switching
    const viewTabs = document.querySelectorAll('.view-tab');
    const calendarView = document.getElementById('calendar-view');
    const listView = document.getElementById('list-view');
    const currentViewInput = document.getElementById('current-view');

    function switchToView(view) {
        // Update active tab
        viewTabs.forEach(t => {
            t.classList.toggle('active', t.dataset.view === view);
        });

        // Show/hide views
        if (view === 'calendar') {
            calendarView.style.display = 'block';
            listView.style.display = 'none';
        } else {
            calendarView.style.display = 'none';
            listView.style.display = 'block';
        }

        // Update hidden input for form submission
        if (currentViewInput) {
            currentViewInput.value = view;
        }
    }

    viewTabs.forEach(tab => {
        tab.addEventListener('click', function() {
            switchToView(this.dataset.view);
        });
    });

    // Check URL for view parameter and restore view
    const urlParams = new URLSearchParams(window.location.search);
    const savedView = urlParams.get('view');
    if (savedView === 'list') {
        switchToView('list');
    }

    // Helper function to populate form and show/hide delete button
    function populateForm(date, dayHours, nightHours, hasEntry) {
        dateInput.value = date;

        // Convert decimal hours to hours and minutes
        const dayH = Math.floor(dayHours);
        const dayM = Math.round((dayHours - dayH) * 60);
        const nightH = Math.floor(nightHours);
        const nightM = Math.round((nightHours - nightH) * 60);

        if (dayHoursInput) dayHoursInput.value = dayH;
        if (dayMinutesInput) dayMinutesInput.value = dayM;
        if (nightHoursInput) nightHoursInput.value = nightH;
        if (nightMinutesInput) nightMinutesInput.value = nightM;

        if (deleteBtn) {
            deleteBtn.style.display = hasEntry ? 'block' : 'none';
        }

        form.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }

    // Calendar day clicks
    calendarDays.forEach(day => {
        day.addEventListener('click', function() {
            const date = this.dataset.date;
            if (!date) return;

            if (this.classList.contains('has-entry') && this.title) {
                const match = this.title.match(/Day: ([\d.]+)h, Night: ([\d.]+)h/);
                if (match) {
                    populateForm(date, parseFloat(match[1]), parseFloat(match[2]), true);
                }
            } else {
                populateForm(date, 0, 0, false);
            }
        });
    });

    // List view row clicks (edit button or row)
    const entryRows = document.querySelectorAll('.entry-row');
    const editButtons = document.querySelectorAll('.edit-entry-btn');

    function handleEditEntry(row) {
        const date = row.dataset.date;
        const dayHours = parseFloat(row.dataset.dayHours) || 0;
        const nightHours = parseFloat(row.dataset.nightHours) || 0;

        // Highlight selected row
        entryRows.forEach(r => r.classList.remove('selected'));
        row.classList.add('selected');

        populateForm(date, dayHours, nightHours, true);
    }

    editButtons.forEach(btn => {
        btn.addEventListener('click', function(e) {
            e.stopPropagation();
            const row = this.closest('.entry-row');
            handleEditEntry(row);
        });
    });

    // Also allow clicking the row itself (but not on action buttons)
    entryRows.forEach(row => {
        row.addEventListener('click', function(e) {
            // Don't trigger if clicking on a button or form
            if (e.target.closest('button') || e.target.closest('form')) {
                return;
            }
            handleEditEntry(this);
        });
    });
});
