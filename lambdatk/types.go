package lambdatk

import (
	"context"
	"encoding/json"
	"net/http"
)

type Dispatcher interface {
	HandleHttpRequest(context.Context, HttpEvent) (EventResult, error)
	HandleTestHttpRequest(eventStr string)
	HandleTestOpHttpRequest(op string, payloadStr string)
}

type OpHandler interface {
	Handle(ctx context.Context, evt HandlerEvent) (interface{}, error)
}

// When receiving an HTTP-driven lambda request, the payload is the full HTTP
// request. But the interesting part for us is the Body. If we JSON-decode into
// this, it'll dump all of the extraneous HTTP parts that we don't need.
type HttpEvent struct {
	Header     http.Header `json:"header"`
	Host       string      `json:"host"`
	RemoteAddr string      `json:"remoteAddr"`
	Body       string      `json:"body"`
}

type HttpMetadata struct {
	Header     http.Header
	Host       string
	RemoteAddr string
}

type HandlerEvent struct {
	Args         json.RawMessage
	HttpMetadata HttpMetadata
}

type HandlerCall struct {
	Op   string          `json:"op"`
	Args json.RawMessage `json:"args"`
}

type EventResult struct {
	Result interface{}
	Error  string
}
