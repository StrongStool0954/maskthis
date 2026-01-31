package shortener

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"time"
)

const charset = "abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"

type URL struct {
	ID        int        `json:"id"`
	ShortCode string     `json:"short_code"`
	LongURL   string     `json:"long_url"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	Clicks    int        `json:"clicks"`
}

func GenerateShortCode(length int) string {
	b := make([]byte, length)
	rand.Read(b)

	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}

	return string(b)
}

func CreateShortURL(db *sql.DB, longURL, clientIP string, expiresInHours int) (*URL, error) {
	// Validate URL
	if !isValidURL(longURL) {
		return nil, errors.New("invalid URL")
	}

	// Generate unique short code
	shortCode, err := generateUniqueCode(db, 6)
	if err != nil {
		return nil, err
	}

	// Calculate expiration
	var expiresAt *time.Time
	if expiresInHours > 0 {
		expiry := time.Now().Add(time.Duration(expiresInHours) * time.Hour)
		expiresAt = &expiry
	}

	// Insert into database
	query := `INSERT INTO urls (short_code, long_url, created_by_ip, expires_at)
              VALUES (?, ?, ?, ?)`

	result, err := db.Exec(query, shortCode, longURL, clientIP, expiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert URL: %w", err)
	}

	id, _ := result.LastInsertId()

	return &URL{
		ID:        int(id),
		ShortCode: shortCode,
		LongURL:   longURL,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
		Clicks:    0,
	}, nil
}

func GetLongURL(db *sql.DB, shortCode string) (string, error) {
	var longURL string
	var expiresAt sql.NullTime

	query := `SELECT long_url, expires_at FROM urls
              WHERE short_code = ?
              AND (expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP)
              LIMIT 1`

	err := db.QueryRow(query, shortCode).Scan(&longURL, &expiresAt)
	if err == sql.ErrNoRows {
		return "", errors.New("short URL not found or expired")
	}
	if err != nil {
		return "", err
	}

	// Increment click counter (async, don't wait)
	go incrementClicks(db, shortCode)

	return longURL, nil
}

func incrementClicks(db *sql.DB, shortCode string) {
	query := `UPDATE urls SET clicks = clicks + 1 WHERE short_code = ?`
	db.Exec(query)
}

func generateUniqueCode(db *sql.DB, length int) (string, error) {
	maxAttempts := 10

	for i := 0; i < maxAttempts; i++ {
		code := GenerateShortCode(length)

		// Check if exists
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM urls WHERE short_code = ?)", code).Scan(&exists)

		if err != nil {
			return "", err
		}

		if !exists {
			return code, nil
		}
	}

	return "", errors.New("failed to generate unique short code")
}

func isValidURL(rawURL string) bool {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	// Must be http or https
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return false
	}

	// Must have a host
	if parsed.Host == "" {
		return false
	}

	return true
}

func GetStats(db *sql.DB, shortCode string) (*URL, error) {
	var u URL

	query := `SELECT id, short_code, long_url, created_at, expires_at, clicks
              FROM urls WHERE short_code = ? LIMIT 1`

	var expiresAt sql.NullTime
	err := db.QueryRow(query, shortCode).Scan(
		&u.ID, &u.ShortCode, &u.LongURL, &u.CreatedAt, &expiresAt, &u.Clicks,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("short URL not found")
	}
	if err != nil {
		return nil, err
	}

	if expiresAt.Valid {
		u.ExpiresAt = &expiresAt.Time
	}

	return &u, nil
}
