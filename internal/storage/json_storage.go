package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"driving-hours/internal/models"
)

type JSONStorage struct {
	dataDir string
	mu      sync.RWMutex
}

func NewJSONStorage(dataDir string) (*JSONStorage, error) {
	// Create data directories if they don't exist
	usersDir := filepath.Join(dataDir, "users")
	if err := os.MkdirAll(usersDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create users directory: %w", err)
	}

	return &JSONStorage{
		dataDir: dataDir,
	}, nil
}

// writeFile writes data atomically using a temp file and rename
func (s *JSONStorage) writeFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	tmpFile, err := os.CreateTemp(dir, ".tmp-*")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()

	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return err
	}

	if err := tmpFile.Chmod(0600); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return err
	}

	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpPath)
		return err
	}

	return os.Rename(tmpPath, path)
}

func (s *JSONStorage) readFile(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// User operations

func (s *JSONStorage) GetUser(id string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	path := filepath.Join(s.dataDir, "users", id+".json")
	var user models.User
	if err := s.readFile(path, &user); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *JSONStorage) GetUserByEmail(email string) (*models.User, error) {
	// Check admin first
	admin, err := s.GetAdmin()
	if err != nil {
		return nil, err
	}
	if admin != nil && admin.Email == email {
		return admin, nil
	}

	// Search through users
	users, err := s.GetAllUsers()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, nil
}

func (s *JSONStorage) GetAllUsers() ([]*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	usersDir := filepath.Join(s.dataDir, "users")
	entries, err := os.ReadDir(usersDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var users []*models.User
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		var user models.User
		path := filepath.Join(usersDir, entry.Name())
		if err := s.readFile(path, &user); err != nil {
			continue
		}
		users = append(users, &user)
	}

	return users, nil
}

func (s *JSONStorage) GetDrivers() ([]*models.User, error) {
	users, err := s.GetAllUsers()
	if err != nil {
		return nil, err
	}

	var drivers []*models.User
	for _, user := range users {
		if user.Role == models.RoleDriver {
			drivers = append(drivers, user)
		}
	}

	return drivers, nil
}

func (s *JSONStorage) SaveUser(user *models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	user.UpdatedAt = time.Now()

	data, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		return err
	}

	path := filepath.Join(s.dataDir, "users", user.ID+".json")
	return s.writeFile(path, data)
}

func (s *JSONStorage) DeleteUser(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := filepath.Join(s.dataDir, "users", id+".json")
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// Session operations

type sessionsFile struct {
	Sessions map[string]*models.Session `json:"sessions"`
}

func (s *JSONStorage) loadSessions() (*sessionsFile, error) {
	path := filepath.Join(s.dataDir, "sessions.json")
	var sf sessionsFile
	if err := s.readFile(path, &sf); err != nil {
		if os.IsNotExist(err) {
			return &sessionsFile{Sessions: make(map[string]*models.Session)}, nil
		}
		return nil, err
	}
	if sf.Sessions == nil {
		sf.Sessions = make(map[string]*models.Session)
	}
	return &sf, nil
}

func (s *JSONStorage) saveSessions(sf *sessionsFile) error {
	data, err := json.MarshalIndent(sf, "", "  ")
	if err != nil {
		return err
	}
	path := filepath.Join(s.dataDir, "sessions.json")
	return s.writeFile(path, data)
}

func (s *JSONStorage) GetSession(token string) (*models.Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sf, err := s.loadSessions()
	if err != nil {
		return nil, err
	}

	session, exists := sf.Sessions[token]
	if !exists {
		return nil, nil
	}

	if session.IsExpired() {
		return nil, nil
	}

	return session, nil
}

func (s *JSONStorage) SaveSession(session *models.Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sf, err := s.loadSessions()
	if err != nil {
		return err
	}

	sf.Sessions[session.Token] = session
	return s.saveSessions(sf)
}

func (s *JSONStorage) DeleteSession(token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sf, err := s.loadSessions()
	if err != nil {
		return err
	}

	delete(sf.Sessions, token)
	return s.saveSessions(sf)
}

func (s *JSONStorage) CleanExpiredSessions() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sf, err := s.loadSessions()
	if err != nil {
		return err
	}

	for token, session := range sf.Sessions {
		if session.IsExpired() {
			delete(sf.Sessions, token)
		}
	}

	return s.saveSessions(sf)
}

// Admin operations

func (s *JSONStorage) GetAdmin() (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	path := filepath.Join(s.dataDir, "admin.json")
	var admin models.User
	if err := s.readFile(path, &admin); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return &admin, nil
}

func (s *JSONStorage) SaveAdmin(admin *models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	admin.UpdatedAt = time.Now()

	data, err := json.MarshalIndent(admin, "", "  ")
	if err != nil {
		return err
	}

	path := filepath.Join(s.dataDir, "admin.json")
	return s.writeFile(path, data)
}
