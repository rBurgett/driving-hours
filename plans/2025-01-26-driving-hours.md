# Driving Hours Application - Implementation Plan

## Overview
A Go web application for tracking student driving hours toward getting a driver's license. Features server-side rendering, JSON file storage, and role-based user management.

## Confirmed Decisions
- **Framework**: Chi router (lightweight, uses Go standard context)
- **Time Input**: Hours and minutes (separate fields)
- **Calendar**: Full navigation (prev/next months)
- **Admin Role**: Management only (no driving log for admins)

---

## Project Structure

```
driving-hours/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── auth/
│   │   ├── argon2.go              # Argon2id password hashing
│   │   ├── session.go             # Session management
│   │   └── middleware.go          # Auth middleware
│   ├── config/
│   │   └── config.go              # Configuration
│   ├── handlers/
│   │   ├── auth.go                # Login/logout
│   │   ├── admin.go               # Admin dashboard
│   │   ├── driver.go              # Driver dashboard
│   │   └── profile.go             # Profile management
│   ├── models/
│   │   ├── user.go                # User struct
│   │   ├── session.go             # Session struct
│   │   └── driving_log.go         # Driving log struct
│   ├── storage/
│   │   ├── storage.go             # Storage interface
│   │   ├── json_storage.go        # JSON implementation
│   │   └── init.go                # First-run initialization
│   ├── middleware/
│   │   └── csrf.go                # CSRF protection
│   ├── templates/
│   │   └── renderer.go            # Template rendering
│   └── utils/
│       ├── time.go                # Greeting helper
│       └── validation.go          # Input validation
├── web/
│   ├── templates/
│   │   ├── layouts/
│   │   │   └── base.html
│   │   ├── partials/
│   │   │   ├── nav.html
│   │   │   ├── flash.html
│   │   │   └── calendar.html
│   │   ├── auth/
│   │   │   └── login.html
│   │   ├── admin/
│   │   │   ├── dashboard.html
│   │   │   ├── users.html
│   │   │   ├── user_form.html
│   │   │   ├── driver_stats.html
│   │   │   ├── driver_hours.html
│   │   │   └── profile.html
│   │   └── driver/
│   │       ├── dashboard.html
│   │       └── profile.html
│   └── static/
│       ├── css/
│       │   └── styles.css
│       └── js/
│           ├── calendar.js
│           ├── fireworks.js
│           └── forms.js
├── data/                          # JSON storage (gitignored)
│   ├── admin.json
│   ├── sessions.json
│   └── users/
│       └── {user_id}.json
├── Dockerfile
├── Makefile
├── README.md
├── CLAUDE.md
├── go.mod
└── go.sum
```

---

## Data Models

### User JSON Structure
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "name": "Grace",
  "password_hash": "$argon2id$v=19$m=65536,t=3,p=4$salt$hash",
  "role": "driver",
  "required_day_hours": 40,
  "required_night_hours": 10,
  "created_at": "2025-01-25T10:00:00Z",
  "updated_at": "2025-01-25T10:00:00Z",
  "driving_log": {
    "2025-01-15": { "day_hours": 1.5, "night_hours": 0.5 }
  }
}
```

### Argon2id Parameters (BitWarden Defaults)
- Memory: 64 MB (65536 KB)
- Iterations: 3
- Parallelism: 4
- Salt: 16 bytes (random)
- Key: 32 bytes

---

## Routes

```
GET  /                  → Redirect to role-appropriate dashboard
GET  /login             → Login page
POST /login             → Process login
POST /logout            → Logout

# Driver Routes (RequireDriver middleware)
GET  /driver            → Driver dashboard (stats, calendar, log form)
POST /driver/log        → Log driving hours
GET  /driver/profile    → Profile page
POST /driver/profile    → Update name/password

# Admin Routes (RequireAdmin middleware)
GET  /admin             → Admin dashboard
GET  /admin/users       → List all users
GET  /admin/users/new   → New user form
POST /admin/users       → Create user
GET  /admin/users/{id}  → View user stats
GET  /admin/users/{id}/edit   → Edit user form
POST /admin/users/{id}        → Update user
GET  /admin/users/{id}/hours  → Edit hours form
POST /admin/users/{id}/hours  → Update hours
GET  /admin/profile     → Admin profile
POST /admin/profile     → Update admin name/password
```

---

## Implementation Status

All phases completed:

1. **Phase 1: Core Infrastructure** - go.mod, config, models, auth, storage
2. **Phase 2: Authentication** - Sessions, middleware, CSRF, login/logout handlers
3. **Phase 3: Templates** - Renderer, base layout, partials, login page
4. **Phase 4: Admin Features** - Admin handlers and templates
5. **Phase 5: Driver Features** - Driver handlers, calendar, greeting utility
6. **Phase 6: Profile & Validation** - Input validation utilities
7. **Phase 7: Static Assets** - CSS, JavaScript (calendar, fireworks, forms)
8. **Phase 8: Entry Point** - main.go with router setup
9. **Phase 9: Build & Deploy** - Makefile, Dockerfile, README, .gitignore, CLAUDE.md

---

## Security Measures
- CSRF tokens on all forms (gorilla/csrf)
- Argon2id password hashing with proper parameters
- HttpOnly, Secure cookies
- Role-based middleware
- Input validation server-side
- Template auto-escaping (html/template)
- File permissions (0600 for data files)
- Atomic file writes (temp file + rename)

---

## Dependencies
```
github.com/go-chi/chi/v5 v5.2.4
github.com/gorilla/csrf v1.7.3
github.com/google/uuid v1.6.0
golang.org/x/crypto v0.21.0
```
