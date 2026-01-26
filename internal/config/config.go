package config

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type Config struct {
	Port      int
	DataDir   string
	CSRFKey   []byte
	IsProd    bool
}

func Load() (*Config, error) {
	port := 8080
	if p := os.Getenv("PORT"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			port = parsed
		}
	}

	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "data"
	}

	csrfKey, err := getCSRFKey(dataDir)
	if err != nil {
		return nil, err
	}

	isProd := os.Getenv("ENV") == "production"

	return &Config{
		Port:    port,
		DataDir: dataDir,
		CSRFKey: csrfKey,
		IsProd:  isProd,
	}, nil
}

func getCSRFKey(dataDir string) ([]byte, error) {
	if key := os.Getenv("CSRF_KEY"); key != "" {
		return base64.StdEncoding.DecodeString(key)
	}

	// Try to load existing key from file
	keyFile := filepath.Join(dataDir, ".csrf_key")
	if data, err := os.ReadFile(keyFile); err == nil {
		key, err := base64.StdEncoding.DecodeString(string(data))
		if err == nil && len(key) == 32 {
			return key, nil
		}
	}

	// Ensure data directory exists
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// Generate a new random key
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}

	// Persist the key for future restarts
	encoded := base64.StdEncoding.EncodeToString(key)
	if err := os.WriteFile(keyFile, []byte(encoded), 0600); err != nil {
		return nil, fmt.Errorf("failed to save CSRF key: %w", err)
	}

	return key, nil
}
