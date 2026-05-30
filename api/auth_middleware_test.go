package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/logger"

	"proxy-hub/api/h"
	"proxy-hub/model"
	"proxy-hub/utils"
)

func TestAuthProtectedRouteRejectsNoToken(t *testing.T) {
	t.Parallel()

	app := fiber.New()
	_, v1 := h.NewAPI(app, &utils.AppConfig{
		APITitle:   "Proxy Hub API",
		APIVersion: "test",
		DocsPath:   "/docs",
	})
	t.Cleanup(func() {
		_ = app.Shutdown()
	})
	h.HumaRegister(v1, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/private",
		OperationID: "test-private",
	}, func(context.Context, *struct{}) (*h.MessageResponse, error) {
		return h.NewMessageResponse("ok"), nil
	})

	resp := mustAuthTestRequest(t, app, "", http.MethodGet, "/api/v1/private")
	if got := resp.StatusCode; got != http.StatusBadRequest {
		t.Fatalf("protected route without token status = %d, want %d", got, http.StatusBadRequest)
	}
}

func TestAuthProtectedRouteAcceptsValidToken(t *testing.T) {
	if err := model.InitWithDSN(":memory:", int(logger.Silent), true); err != nil {
		t.Fatalf("InitWithDSN(:memory:) failed: %v", err)
	}
	t.Cleanup(model.DBClose)

	app := fiber.New()
	_, v1 := h.NewAPI(app, &utils.AppConfig{
		APITitle:   "Proxy Hub API",
		APIVersion: "test",
		DocsPath:   "/docs",
	})
	t.Cleanup(func() {
		_ = app.Shutdown()
	})
	h.HumaRegister(v1, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/private",
		OperationID: "test-private-auth",
	}, func(context.Context, *struct{}) (*h.MessageResponse, error) {
		return h.NewMessageResponse("ok"), nil
	})

	token, _ := createTestUserAndToken(t)
	resp := mustAuthTestRequest(t, app, token, http.MethodGet, "/api/v1/private")
	if got := resp.StatusCode; got != http.StatusOK {
		t.Fatalf("protected route with token status = %d, want %d", got, http.StatusOK)
	}
}

func TestAuthSigninAllowsNoToken(t *testing.T) {
	if err := model.InitWithDSN(":memory:", int(logger.Silent), true); err != nil {
		t.Fatalf("InitWithDSN(:memory:) failed: %v", err)
	}
	t.Cleanup(model.DBClose)

	app := fiber.New()
	_, v1 := h.NewAPI(app, &utils.AppConfig{
		APITitle:   "Proxy Hub API",
		APIVersion: "test",
		DocsPath:   "/docs",
	})
	t.Cleanup(func() {
		_ = app.Shutdown()
	})
	h.HumaRegister(v1, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/user/signin",
		OperationID: "test-signin-auth-check",
	}, func(ctx context.Context, _ *struct{}) (*struct{}, error) {
		return &struct{}{}, nil
	})

	resp := mustAuthTestRequest(t, app, "", http.MethodPost, "/api/v1/user/signin")
	if got := resp.StatusCode; got < 200 || got >= 300 {
		t.Fatalf("signin without token status = %d, want 2xx (auth middleware should bypass signin)", got)
	}
}
