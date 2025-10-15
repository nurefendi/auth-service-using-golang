package middleware

import (
    "net/http/httptest"
    "os"
    "testing"
    "time"

    "github.com/gofiber/fiber/v2"
)

func TestLockoutRegisterAndReset(t *testing.T) {
    username := "test@example.com"
    // ensure clean state
    ResetFailedAttempts(username)

    max := 3
    os.Setenv("LOGIN_MAX_ATTEMPTS", "3")
    os.Setenv("LOGIN_LOCKOUT_SECONDS", "2")

    // register failed attempts less than max
    for i := 0; i < max-1; i++ {
        locked, remaining := RegisterFailedAttempt(username)
        if locked {
            t.Fatalf("unexpected lock at attempt %d", i+1)
        }
        if remaining != (max-(i+1)) {
            t.Fatalf("unexpected remaining attempts: got %d expected %d", remaining, max-(i+1))
        }
    }

    // next attempt should lock
    locked, _ := RegisterFailedAttempt(username)
    if !locked {
        t.Fatalf("expected account to be locked after max attempts")
    }

    // should report locked
    ok, until := AccountIsLocked(username)
    if !ok {
        t.Fatalf("expected account locked")
    }
    if time.Now().After(until) {
        t.Fatalf("locked until should be in future")
    }

    // wait lockout expiry
    time.Sleep(3 * time.Second)
    ok2, _ := AccountIsLocked(username)
    if ok2 {
        t.Fatalf("expected lock to expire")
    }

    // reset explicitly
    RegisterFailedAttempt(username)
    ResetFailedAttempts(username)
    ok3, _ := AccountIsLocked(username)
    if ok3 {
        t.Fatalf("expected not locked after reset")
    }
}

func TestRateLimitMiddleware(t *testing.T) {
    os.Setenv("RATE_LIMIT_REQUESTS", "2")
    os.Setenv("RATE_LIMIT_WINDOW_SECONDS", "1")

    app := fiber.New()
    app.Get("/test", SetRateLimit(), func(c *fiber.Ctx) error {
        return c.SendString("ok")
    })

    // first two requests should be ok
    for i := 0; i < 2; i++ {
        req := httptest.NewRequest("GET", "/test", nil)
        resp, _ := app.Test(req, -1)
        if resp.StatusCode != 200 {
            t.Fatalf("expected 200, got %d on attempt %d", resp.StatusCode, i+1)
        }
    }

    // third request should be rate limited
    req := httptest.NewRequest("GET", "/test", nil)
    resp, _ := app.Test(req, -1)
    if resp.StatusCode != 429 {
        t.Fatalf("expected 429, got %d on attempt 3", resp.StatusCode)
    }
}
