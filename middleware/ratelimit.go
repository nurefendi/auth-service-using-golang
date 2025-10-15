package middleware

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

// simple in-memory rate limiter and account lockout

type ipCounter struct {
	Count     int
	ExpiresAt time.Time
}

var (
	ipStore = sync.Map{} // map[string]*ipCounter
	ipMutex = sync.Mutex{}
	// account lockout store: map[username]*lockEntry
	lockStore = sync.Map{}
)

type lockEntry struct {
	Attempts    int
	LockedUntil time.Time
}

// defaults from env
func getEnvInt(name string, def int) int {
	v := os.Getenv(name)
	if v == "" {
		return def
	}
	var parsed int
	_, err := fmt.Sscanf(v, "%d", &parsed)
	if err != nil {
		return def
	}
	return parsed
}

// SetRateLimit returns a middleware that limits requests per IP.
// Controlled by env RATE_LIMIT_REQUESTS and RATE_LIMIT_WINDOW_SECONDS.
func SetRateLimit() fiber.Handler {
	max := getEnvInt("RATE_LIMIT_REQUESTS", 20)
	windowSec := getEnvInt("RATE_LIMIT_WINDOW_SECONDS", 60)
	window := time.Duration(windowSec) * time.Second

	return func(c *fiber.Ctx) error {
		ip := c.IP()
		now := time.Now()

		val, _ := ipStore.LoadOrStore(ip, &ipCounter{Count: 0, ExpiresAt: now.Add(window)})
		counter := val.(*ipCounter)

		ipMutex.Lock()
		if now.After(counter.ExpiresAt) {
			counter.Count = 0
			counter.ExpiresAt = now.Add(window)
		}
		counter.Count++
		cnt := counter.Count
		ipMutex.Unlock()

		if cnt > max {
			return c.Status(fiber.StatusTooManyRequests).SendString("Too many requests")
		}
		return c.Next()
	}
}

// Account lockout helpers
func AccountIsLocked(username string) (bool, time.Time) {
	val, ok := lockStore.Load(username)
	if !ok {
		return false, time.Time{}
	}
	le := val.(*lockEntry)
	if time.Now().Before(le.LockedUntil) {
		return true, le.LockedUntil
	}
	return false, time.Time{}
}

// RegisterFailedAttempt increments attempt counter and returns whether locked and remaining attempts
func RegisterFailedAttempt(username string) (locked bool, attemptsLeft int) {
	maxAttempts := getEnvInt("LOGIN_MAX_ATTEMPTS", 5)
	lockSeconds := getEnvInt("LOGIN_LOCKOUT_SECONDS", 900)

	val, _ := lockStore.LoadOrStore(username, &lockEntry{Attempts: 0, LockedUntil: time.Time{}})
	le := val.(*lockEntry)

	// if currently locked and not expired, keep locked
	if time.Now().Before(le.LockedUntil) {
		return true, 0
	}

	le.Attempts++
	if le.Attempts >= maxAttempts {
		le.LockedUntil = time.Now().Add(time.Duration(lockSeconds) * time.Second)
		le.Attempts = 0
		lockStore.Store(username, le)
		return true, 0
	}
	lockStore.Store(username, le)
	return false, maxAttempts - le.Attempts
}

// ResetFailedAttempts clears counters for username
func ResetFailedAttempts(username string) {
	lockStore.Delete(username)
}
