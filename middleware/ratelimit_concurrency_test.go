package middleware

import (
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"github.com/gofiber/fiber/v2"
)

// Stress test: many concurrent requests from same IP should be rate-limited
func TestRateLimitConcurrencyStress(t *testing.T) {
	os.Setenv("RATE_LIMIT_REQUESTS", "10")
	os.Setenv("RATE_LIMIT_WINDOW_SECONDS", "1")

	app := fiber.New()
	app.Get("/stress", SetRateLimit(), func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	total := 100
	var wg sync.WaitGroup
	wg.Add(total)
	success := 0
	mu := sync.Mutex{}

	for i := 0; i < total; i++ {
		go func() {
			defer wg.Done()
			req := httptest.NewRequest("GET", "/stress", nil)
			// ensure same IP
			req.Header.Set("X-Real-IP", "10.0.0.1")
			resp, _ := app.Test(req, -1)
			if resp.StatusCode == 200 {
				mu.Lock()
				success++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	if success > 12 { // allow small race margin
		t.Fatalf("too many successes: %d", success)
	}
}
