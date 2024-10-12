package log

import (
	"bytes"
	"context"

	"github.com/cloudokyo/cast"
	"github.com/labstack/echo/v4"
)

const (
	// The key attachs the logger data to the context
	ContextKey = Key("logger")

	// The auth account key, refer to pkg/http/filter/auth.ContextKey
	AccountKey = "account"

	// The request id key attaches to the context
	RequestKey = "id"

	// The user id key attaches to the context
	UserKey = "uid"

	// The device id key attaches to the context
	DeviceKey = "did"

	// The client id key attaches to the context
	ClientKey = "cid"

	// The request user agent key attaches to the context
	UserAgentKey = "agent"
)

// Define the context key to inject data in to context
type Key string

// Define the injection data struct to context
type ContextData map[string]any

func (c ContextData) IsEmpty() bool {
	return len(c) == 0
}

func (c ContextData) Value() map[string]any {
	return c
}

func (c ContextData) String() string {
	// Ignore empty data
	if c.IsEmpty() {
		return ""
	}

	// Build the output: k1=v1 k2=v2 k3=v3
	output := bytes.NewBufferString("")
	for k, v := range c {
		if output.Len() > 1 {
			output.WriteString(" ")
		}
		output.WriteString(k)
		output.WriteString("=")
		output.WriteString(cast.ToString(v))
	}
	output.WriteString("")

	return output.String()
}

// Attach the data to the context.Context
//
// Example
//
//	ctx := log.WithValue(ctx, "key", "value")
func WithValue(ctx context.Context, key string, value any) context.Context {
	// Get or init data
	data := Detach(ctx)
	if data == nil {
		data = ContextData{}
	}

	// Inject new data
	data[key] = value

	// Inject data to context
	return context.WithValue(ctx, ContextKey, data)
}

// An alias function of log.Attach.
// We get and inject the data into echo request context.
//
// Example
//
//	ctx := log.GetContext(c)
var GetContext = Attach

// An alias function of log.Attach.
// We get and inject the data into echo request context.
//
// Example
//
//	ctx := log.WithContext(c)
var WithContext = Attach

// Attach the logger data into echo.Context
//
// Injection
//
//   - id: The request id
//   - uid: The user id (TBD)
//   - did: The user device id (TBD)
//   - cid: The user active channel id (TBD)
//   - agent: The user browser UserAgent (TBD)
//
// Example
//
//	ctx := log.Attach(ctx)
func Attach(c echo.Context) context.Context {
	// Init the logger injection data
	data := ContextData{}

	// Get the echo context
	ctx := c.Request().Context()

	// Get the request id from context
	if value := c.Request().Header.Get(echo.HeaderXRequestID); value != "" {
		data[RequestKey] = value
	}

	return context.WithValue(ctx, ContextKey, data)
}

// Detach the logger data injected to the context
//
// Example
//
//	log.Detach(ctx) --> {"id": "xxx"}
func Detach(ctx context.Context) ContextData {
	if value, ok := ctx.Value(ContextKey).(ContextData); ok {
		return value
	}
	return nil
}

// Get the requestId has injected to the context
//
// Example
//
//	log.RequestId(ctx)
func RequestId(ctx context.Context) string {
	if data := Detach(ctx); len(data) == 0 {
		return ""
	} else if value, ok := data[RequestKey]; ok {
		return cast.ToString(value)
	} else {
		return ""
	}
}
