package log_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudokyo/log"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestLogContextDetach(t *testing.T) {
	ctx := context.Background()

	assert.Empty(t, log.RequestId(ctx))

	fmt.Println("--> OUT:", log.RequestId(ctx))

	id := "test-request-id"
	ctx = log.WithValue(ctx, log.RequestKey, id)

	data := log.Detach(ctx)

	fmt.Println("--> OUT:", data)
	fmt.Println("--> OUT:", log.RequestId(ctx))

	assert.Contains(t, data, "id")
	assert.Equal(t, id, data["id"])
	assert.Equal(t, id, log.RequestId(ctx))
}

func TestLogContextAttach(t *testing.T) {
	// Init the echo server
	engine := echo.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// Init the echo context
	id := "test-request-id"
	c := engine.NewContext(req, rec)
	c.Request().Header.Set(echo.HeaderXRequestID, id)

	// Init the context
	ctx := log.WithContext(c)

	// Detach the data from context
	data := log.Detach(ctx)

	fmt.Println("--> OUT:", data)

	assert.Contains(t, data, "id")
	assert.Equal(t, id, data["id"])
}
