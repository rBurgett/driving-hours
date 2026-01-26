# CLAUDE.md - Project Context for Claude

## Project Overview

Driving Hours Tracker - A Go web application for tracking student driving hours toward getting a driver's license.

## Tech Stack

- **Language**: Go 1.22
- **Router**: Chi (go-chi/chi/v5)
- **Templates**: html/template with layouts and partials
- **Storage**: JSON files (no database)
- **Auth**: Argon2id password hashing, cookie-based sessions
- **CSRF**: gorilla/csrf

## Architecture

### Directory Structure

```
cmd/server/main.go     - Entry point, router setup
internal/
  auth/                - Password hashing, session management, middleware
  config/              - Environment configuration
  handlers/            - HTTP handlers for auth, admin, driver
  middleware/          - CSRF middleware
  models/              - User, Session, DrivingLog structs
  storage/             - JSON file storage with mutex protection
  templates/           - Template renderer with FuncMap
  utils/               - Time (greeting), validation helpers
web/
  templates/           - HTML templates (layouts, partials, pages)
  static/              - CSS and JavaScript
data/                  - JSON storage (gitignored)
```

### Key Patterns

1. **Handler structure**: Each handler struct holds storage and renderer references
2. **Middleware chain**: Logger -> Recoverer -> RealIP -> CSRF -> Auth (per route group)
3. **Template rendering**: Base layout with content blocks, partials for reusable components
4. **Storage**: Interface-based with mutex protection for concurrent access
5. **Atomic writes**: Temp file + rename for data integrity

### Data Models

- **User**: ID, email, name, password_hash, role, required hours, driving_log
- **Session**: Token, user_id, expires_at
- **DrivingLog**: Map of date strings to DayEntry (day_hours, night_hours)

### Routes

- `/` - Redirect to role-appropriate dashboard
- `/login`, `/logout` - Authentication
- `/driver/*` - Driver dashboard, log hours, profile (RequireDriver)
- `/admin/*` - Admin dashboard, user management (RequireAdmin)

## Common Tasks

### Adding a new handler

1. Create handler struct in `internal/handlers/`
2. Add constructor with storage and renderer
3. Register routes in `cmd/server/main.go`
4. Create template in `web/templates/`

### Adding a new template

1. Create `.html` file in appropriate `web/templates/` subdirectory
2. Use `{{define "content"}}` block for page content
3. Access CSRF token via `{{.CSRFField}}`
4. Use template functions from `internal/templates/renderer.go`

### Modifying storage

1. Update interface in `internal/storage/storage.go`
2. Implement in `internal/storage/json_storage.go`
3. Use mutex for all read/write operations
4. Use atomic writes for file persistence

## Testing

```bash
make test           # Run all tests
make test-coverage  # Run with coverage report
```

## Development

```bash
make dev    # Hot reload with air
make run    # Build and run
make lint   # Run linter
```

## Security Notes

- Never log passwords or session tokens
- Always use Argon2id for password hashing
- CSRF tokens required on all POST forms
- Session cookies are HttpOnly, Secure (prod), SameSite=Lax
- File permissions: 0600 for data files, 0700 for directories
