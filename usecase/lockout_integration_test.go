package usecase

import (
    "net/http/httptest"
    "os"
    "testing"
    "time"

    "auth-service/middleware"
    "github.com/gofiber/fiber/v2"
)

// This test simulates repeated failed login attempts through the middleware lockout helpers.
func TestLoginFlowLockoutIntegration(t *testing.T) {
    username := "int-test@example.com"
    // ensure clean
    middleware.ResetFailedAttempts(username)

    os.Setenv("LOGIN_MAX_ATTEMPTS", "3")
    os.Setenv("LOGIN_LOCKOUT_SECONDS", "2")

    // simulate 3 failed attempts
    for i := 0; i < 3; i++ {
        locked, _ := middleware.RegisterFailedAttempt(username)
        if i < 2 && locked {
            t.Fatalf("unexpected lock on attempt %d", i+1)
        }
    }

    // now should be locked
    ok, until := middleware.AccountIsLocked(username)
    if !ok {
        t.Fatalf("expected locked after attempts")
    }
    if time.Now().After(until) {
        t.Fatalf("locked until is in past")
    }

    // simple HTTP app that checks lock
    app := fiber.New()
    app.Post("/login", func(c *fiber.Ctx) error {
        if locked, _ := middleware.AccountIsLocked(username); locked {
            return c.Status(429).SendString("locked")
        }
        return c.SendString("ok")
    })

    req := httptest.NewRequest("POST", "/login", nil)
    resp, _ := app.Test(req, -1)
    if resp.StatusCode != 429 {
        t.Fatalf("expected 429 from locked account, got %d", resp.StatusCode)
    }
}
