package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"

	"proxy-hub/model"
	"proxy-hub/model/tables"
	userService "proxy-hub/service/user"
	"proxy-hub/utils"
)

func createTestUserAndToken(t *testing.T) (string, *tables.UserTable) {
	t.Helper()
	ctx := context.Background()
	var token string
	var u *tables.UserTable
	err := model.Transaction(ctx, func(tx model.DBTx) error {
		var createErr error
		username := fmt.Sprintf("test-%s", utils.NewID())
		password := fmt.Sprintf("pwd-%s", utils.NewID())
		u, createErr = userService.UserCreate(ctx, tx, username, password, "test", "", "", nil)
		if createErr != nil {
			return createErr
		}
		token, createErr = userService.AccessTokenGenerate(ctx, tx, u.ID)
		return createErr
	})
	if err != nil {
		t.Fatalf("createTestUserAndToken: %v", err)
	}
	return token, u
}

func mustAuthTestRequest(t *testing.T, app *fiber.App, token, method, target string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(method, target, nil)
	if err != nil {
		t.Fatalf("new request %s %s: %v", method, target, err)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test %s %s: %v", method, target, err)
	}
	t.Cleanup(func() {
		_ = resp.Body.Close()
	})
	return resp
}

func mustAuthJSONTestRequest(t *testing.T, app *fiber.App, token, method, target, body string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(method, target, strings.NewReader(body))
	if err != nil {
		t.Fatalf("new request %s %s: %v", method, target, err)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test %s %s: %v", method, target, err)
	}
	t.Cleanup(func() {
		_ = resp.Body.Close()
	})
	return resp
}
