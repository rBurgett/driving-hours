# Driving Hours Tracker

A web application for tracking student driving hours toward getting a driver's license. Features role-based access (admin/driver), progress tracking, and a calendar interface.

## Features

- **Role-based access**: Admin and driver accounts with different capabilities
- **Driver dashboard**: Track day and night driving hours with progress bars
- **Calendar view**: Visual representation of logged driving sessions
- **Admin management**: Create/edit drivers, set required hours, view statistics
- **Secure**: Argon2id password hashing, CSRF protection, HTTP-only cookies
- **Simple storage**: JSON file-based storage (no database required)

## Quick Start

### Prerequisites

- Go 1.22 or later

### Running Locally

```bash
# Clone the repository
git clone <repository-url>
cd driving-hours

# Build and run
make run
```

On first run, admin credentials will be printed to the console:

```
========================================
  FIRST RUN - Admin Account Created
========================================
  Email:    admin@localhost
  Password: <random-16-char-password>
========================================
  Please save these credentials!
========================================
```

Open http://localhost:8080 in your browser.

### Development

For hot-reload during development, install [air](https://github.com/cosmtrek/air):

```bash
go install github.com/cosmtrek/air@latest
make dev
```

## Configuration

Configuration is done via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `DATA_DIR` | `data` | Directory for JSON storage |
| `CSRF_KEY` | (random) | Base64-encoded 32-byte key for CSRF protection |
| `ENV` | (empty) | Set to `production` for secure cookies |

## Docker

```bash
# Build image
make docker-build

# Run container
make docker-run
```

Or manually:

```bash
docker build -t driving-hours .
docker run -p 8080:8080 -v ./data:/app/data driving-hours
```

## Project Structure

```
driving-hours/
├── cmd/server/          # Application entry point
├── internal/
│   ├── auth/            # Authentication (Argon2id, sessions, middleware)
│   ├── config/          # Configuration loading
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # CSRF protection
│   ├── models/          # Data models
│   ├── storage/         # JSON file storage
│   ├── templates/       # Template rendering
│   └── utils/           # Utilities (time, validation)
├── web/
│   ├── templates/       # HTML templates
│   └── static/          # CSS and JavaScript
└── data/                # JSON storage (gitignored)
```

## Usage

### Admin Functions

1. **Create drivers**: Add new driver accounts with required day/night hours
2. **View statistics**: See progress for each driver
3. **Edit hours**: Manually adjust logged hours if needed
4. **Manage profiles**: Update driver names, emails, and passwords

### Driver Functions

1. **Log hours**: Enter day and night driving hours with date
2. **View progress**: See progress bars for day and night hour requirements
3. **Calendar**: Click any day to log or edit hours for that date
4. **Celebration**: Fireworks animation when hours are logged

## Security

- Passwords are hashed using Argon2id with BitWarden-recommended parameters
- CSRF protection on all forms
- HTTP-only, secure (in production) session cookies
- Role-based middleware prevents unauthorized access
- Input validation on all user inputs
- Atomic file writes prevent data corruption

## License

MIT
