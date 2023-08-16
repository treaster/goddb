package lambdatk

import (
	"context"
	"encoding/json"
)

type Dispatcher interface {
	HandleRequest(context.Context, Event) (EventResult, error)
	HandleTestRequest(eventStr string)
	HandleTestOpRequest(op string, payloadStr string)
}

type OpHandler interface {
	Handle(ctx context.Context, args json.RawMessage) (interface{}, error)
}

// When receiving an HTTP-driven lambda request, the payload is the full HTTP
// request. But the interesting part for us is the Body. If we JSON-decode into
// this, it'll dump all of the extraneous HTTP parts that we don't need.
type Event struct {
	Body string `json:"body"`
}

type EventBody struct {
	Op   string          `json:"op"`
	Args json.RawMessage `json:"args"`
}

type EventResult struct {
	Result interface{}
	Error  string
}
