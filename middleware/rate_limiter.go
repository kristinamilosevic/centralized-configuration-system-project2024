package middleware

import (
	"log"
	"net/http"
	"sync"
	"time"
)

// RateLimiter predstavlja middleware za ograničavanje brzine zahteva.
type RateLimiter struct {
	Interval  time.Duration // Vremenski interval
	MaxEvents int           // Maksimalni broj događaja
	Lock      sync.Mutex    // Koristi se za sinhronizaciju
	Count     int           // Broj zahteva
}

// NewRateLimiter kreira novi RateLimiter sa zadatim intervalom i maksimalnim brojem događaja.
func NewRateLimiter(maxEvents int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		MaxEvents: maxEvents,
		Interval:  interval,
		Count:     0,
	}
}

// ServeHTTP obrađuje HTTP zahtev i primenjuje rate limiting.
func (rl *RateLimiter) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rl.Lock.Lock()
	defer rl.Lock.Unlock()

	log.Printf("Request received")
	log.Printf("Current count: %d", rl.Count)

	// Provera broja zahteva u poslednjem intervalu
	if rl.Count >= rl.MaxEvents {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		log.Printf("Rate limit exceeded")
		return
	}

	// Ažuriranje broja zahteva
	rl.Count++

	// Resetovanje broja zahteva nakon intervala
	time.AfterFunc(rl.Interval, func() {
		rl.Lock.Lock()
		defer rl.Lock.Unlock()
		rl.Count = 0
		log.Printf("Counter reset")
	})

	// Pozivanje sledećeg handlera
	next(w, r)
}
