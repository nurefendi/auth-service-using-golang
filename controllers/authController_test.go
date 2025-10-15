package controllers

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/nurefendi/auth-service-using-golang/dto"
	"github.com/nurefendi/auth-service-using-golang/middleware"
	"github.com/nurefendi/auth-service-using-golang/usecase"
)

type fakeAuthUsecase struct{}

func (f *fakeAuthUsecase) Register(c *fiber.Ctx) *fiber.Error                             { return nil }
func (f *fakeAuthUsecase) CheckEmailExist(c *fiber.Ctx, email *string) *fiber.Error       { return nil }
func (f *fakeAuthUsecase) CheckUserNameExist(c *fiber.Ctx, userName *string) *fiber.Error { return nil }
func (f *fakeAuthUsecase) Login(c *fiber.Ctx) *fiber.Error                                { return nil }
func (f *fakeAuthUsecase) Logout(c *fiber.Ctx) *fiber.Error                               { return nil }
func (f *fakeAuthUsecase) RefreshToken(c *fiber.Ctx) *fiber.Error                         { return nil }
func (f *fakeAuthUsecase) Me(c *fiber.Ctx) (dto.AuthUserResponse, *fiber.Error) {
	return dto.AuthUserResponse{}, nil
}
func (f *fakeAuthUsecase) CheckAccess(c *fiber.Ctx, r dto.AuthCheckAccessRequest) *fiber.Error {
	return nil
}
func (f *fakeAuthUsecase) MyAcl(c *fiber.Ctx) ([]dto.AuthUserFunction, *fiber.Error) {
	return []dto.AuthUserFunction{{FunctionName: "f1", PortalName: "p1"}}, nil
}

func TestAuthLoginHandler(t *testing.T) {
	// inject fake usecase
	usecase.SetAuthInstance(&fakeAuthUsecase{})
	app := fiber.New()
	app.Post("/login", AuthLogin)

	body := `{"email":"a@b.com","password":"secret"}`
	req := httptest.NewRequest("POST", "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusAccepted {
		t.Fatalf("expected 202 accepted, got %d", resp.StatusCode)
	}
}

func TestGetMyAclHandler(t *testing.T) {
	usecase.SetAuthInstance(&fakeAuthUsecase{})
	app := fiber.New()
	app.Get("/acl", middleware.SetMiddlewareJSON(), GetMyAcl)

	req := httptest.NewRequest("GET", "/acl", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}
	var data map[string]any
	json.NewDecoder(resp.Body).Decode(&data)
	if _, ok := data["data"]; !ok {
		t.Fatalf("expected data field in response")
	}
}

func TestAuthLoginValidationError(t *testing.T) {
	usecase.SetAuthInstance(&fakeAuthUsecase{})
	app := fiber.New()
	app.Post("/login", middleware.SetMiddlewareJSON(), AuthLogin)

	// missing password
	body := `{"email":"invalid-email"}`
	req := httptest.NewRequest("POST", "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusUnprocessableEntity {
		t.Fatalf("expected 422 unprocessable entity, got %d", resp.StatusCode)
	}
}

type fakeAuthUsecaseErr struct{}

func (f *fakeAuthUsecaseErr) Register(c *fiber.Ctx) *fiber.Error                       { return nil }
func (f *fakeAuthUsecaseErr) CheckEmailExist(c *fiber.Ctx, email *string) *fiber.Error { return nil }
func (f *fakeAuthUsecaseErr) CheckUserNameExist(c *fiber.Ctx, userName *string) *fiber.Error {
	return nil
}
func (f *fakeAuthUsecaseErr) Login(c *fiber.Ctx) *fiber.Error {
	return &fiber.Error{Code: fiber.StatusBadRequest, Message: "bad"}
}
func (f *fakeAuthUsecaseErr) Logout(c *fiber.Ctx) *fiber.Error       { return nil }
func (f *fakeAuthUsecaseErr) RefreshToken(c *fiber.Ctx) *fiber.Error { return nil }
func (f *fakeAuthUsecaseErr) Me(c *fiber.Ctx) (dto.AuthUserResponse, *fiber.Error) {
	return dto.AuthUserResponse{}, nil
}
func (f *fakeAuthUsecaseErr) CheckAccess(c *fiber.Ctx, r dto.AuthCheckAccessRequest) *fiber.Error {
	return nil
}
func (f *fakeAuthUsecaseErr) MyAcl(c *fiber.Ctx) ([]dto.AuthUserFunction, *fiber.Error) {
	return nil, nil
}

func TestAuthLoginUsecaseError(t *testing.T) {
	usecase.SetAuthInstance(&fakeAuthUsecaseErr{})
	app := fiber.New()
	app.Post("/login", middleware.SetMiddlewareJSON(), AuthLogin)

	body := `{"email":"a@b.com","password":"secret"}`
	req := httptest.NewRequest("POST", "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected 400 bad request, got %d", resp.StatusCode)
	}
}
