package analytics

import (
	"database/sql"
	"time"
)

type Click struct {
	ID          int       `json:"id"`
	ShortCode   string    `json:"short_code"`
	ClickedAt   time.Time `json:"clicked_at"`
	IPAddress   string    `json:"ip_address"`
	UserAgent   string    `json:"user_agent"`
	Referrer    string    `json:"referrer"`
	CountryCode string    `json:"country_code"`
}

func RecordClick(db *sql.DB, shortCode, ip, userAgent, referrer string) error {
	query := `INSERT INTO clicks (short_code, ip_address, user_agent, referrer)
              VALUES (?, ?, ?, ?)`

	_, err := db.Exec(query, shortCode, ip, userAgent, referrer)
	return err
}

func GetClickCount(db *sql.DB, shortCode string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM clicks WHERE short_code = ?`
	err := db.QueryRow(query, shortCode).Scan(&count)
	return count, err
}
