package lambdatk

import (
	"context"
	"encoding/json"
	"log"
	"time"
)

func MakeDispatcher(handlers map[string]OpHandler) Dispatcher {
	return dispatcher{
		time.Now(),
		handlers,
	}
}

type dispatcher struct {
	startupTime time.Time
	handlers    map[string]OpHandler
}

func (d dispatcher) HandleHttpRequest(ctx context.Context, evt HttpEvent) (EventResult, error) {
	result, err := d.handleRequestHelper(ctx, evt)
	if err != nil {
		return EventResult{
			Error: err.Error(),
		}, nil
	} else {
		return EventResult{
			Result: result,
		}, nil
	}
}

func (d dispatcher) handleRequestHelper(ctx context.Context, evt HttpEvent) (interface{}, error) {
	startEventTime := time.Now()
	defer func() {
		endTime := time.Now()
		startupDeltaMs := endTime.Sub(d.startupTime) / time.Millisecond
		eventDeltaMs := endTime.Sub(startEventTime) / time.Millisecond

		log.Printf("Startup: %dms, Event: %dms", startupDeltaMs, eventDeltaMs)
	}()

	log.Printf("source body is %q", evt.Body)

	// Parse the HTTP body into an (op, args) tuple for dynamodbtk.
	var call HandlerCall
	err := json.Unmarshal([]byte(evt.Body), &call)
	if err != nil {
		return "", newErrf("unable to parse event: %s", err.Error())
	}

	// Find the handler method for the specified op.
	handler, hasHandler := d.handlers[call.Op]
	if !hasHandler {
		return "", newErrf("unknown op %q", call.Op)
	}

	// Invoke the handler with the args, plus some other HTTP metadata.
	handlerEvent := HandlerEvent{
		Args: call.Args,
		HttpMetadata: HttpMetadata{
			Header:     evt.Header,
			Host:       evt.Host,
			RemoteAddr: evt.RemoteAddr,
		},
	}

	result, err := handler.Handle(ctx, handlerEvent)
	if err != nil {
		return "", newErrf("error in %q handler: %s", call.Op, err.Error())
	}

	return result, nil
}

func (d dispatcher) HandleTestHttpRequest(eventStr string) {
	var evt HttpEvent
	err := json.Unmarshal([]byte(eventStr), &evt)
	if err != nil {
		log.Fatalf("error in test JSON: %s", err.Error())
	}

	result, err := d.HandleHttpRequest(context.Background(), evt)
	if err != nil {
		log.Fatalf("error in test call: %s\n", err.Error())
	}
	log.Printf("%+v", result)
}

func (d dispatcher) HandleTestOpHttpRequest(op string, payloadStr string) {
	evtBody := HandlerCall{
		op,
		json.RawMessage(payloadStr),
	}

	bodyJson, err := json.Marshal(evtBody)
	if err != nil {
		log.Fatalf("error in test call: %s\n", err.Error())
	}
	evt := HttpEvent{
		Body: string(bodyJson),
	}

	result, err := d.HandleHttpRequest(context.Background(), evt)
	if err != nil {
		log.Fatalf("error in test call: %s\n", err.Error())
	}
	log.Printf("%+v", result)
}
