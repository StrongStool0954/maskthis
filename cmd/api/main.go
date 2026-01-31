package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bullwinkle/maskthis/internal/analytics"
	"github.com/bullwinkle/maskthis/internal/database"
	"github.com/bullwinkle/maskthis/internal/shortener"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type ShortenRequest struct {
	URL            string `json:"url"`
	ExpiresInHours int    `json:"expires_in_hours,omitempty"`
}

type ShortenResponse struct {
	Success   bool   `json:"success"`
	ShortCode string `json:"short_code,omitempty"`
	ShortURL  string `json:"short_url,omitempty"`
	LongURL   string `json:"long_url,omitempty"`
	Error     string `json:"error,omitempty"`
}

func main() {
	// Load environment
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	log.Println("✅ Connected to Bunny Database")

	// Setup router
	r := mux.NewRouter()

	// API routes
	r.HandleFunc("/api/shorten", shortenHandler).Methods("POST")
	r.HandleFunc("/api/stats/{shortCode}", statsHandler).Methods("GET")
	r.HandleFunc("/api/health", healthHandler).Methods("GET")

	// Redirect route (must be after API routes)
	r.HandleFunc("/{shortCode}", redirectHandler).Methods("GET")

	// Static files (homepage)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/static")))

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("📍 http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok","service":"maskthis"}`))
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Get client IP
	clientIP := getClientIP(r)

	// Create short URL
	url, err := shortener.CreateShortURL(database.DB, req.URL, clientIP, req.ExpiresInHours)
	if err != nil {
		respondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Build response
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	response := ShortenResponse{
		Success:   true,
		ShortCode: url.ShortCode,
		ShortURL:  baseURL + "/" + url.ShortCode,
		LongURL:   url.LongURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	log.Printf("✅ Created short URL: %s -> %s", url.ShortCode, url.LongURL)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["shortCode"]

	// Get long URL
	longURL, err := shortener.GetLongURL(database.DB, shortCode)
	if err != nil {
		// Return 404 page instead of plain text
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`
<!DOCTYPE html>
<html>
<head><title>404 - Link Not Found</title></head>
<body style="font-family: sans-serif; text-align: center; padding: 100px;">
	<h1>🔗 Link Not Found</h1>
	<p>This short URL doesn't exist or has expired.</p>
	<a href="/" style="color: #3b82f6;">← Create a new short link</a>
</body>
</html>`))
		return
	}

	// Record analytics (async)
	go analytics.RecordClick(
		database.DB,
		shortCode,
		getClientIP(r),
		r.UserAgent(),
		r.Referer(),
	)

	log.Printf("➡️  Redirect: %s -> %s", shortCode, longURL)

	// Redirect
	http.Redirect(w, r, longURL, http.StatusFound)
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["shortCode"]

	stats, err := shortener.GetStats(database.DB, shortCode)
	if err != nil {
		respondError(w, "Short URL not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (from proxy/CDN)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Take first IP from the list
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	// Remove port if present
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}

func respondError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ShortenResponse{
		Success: false,
		Error:   message,
	})
}
