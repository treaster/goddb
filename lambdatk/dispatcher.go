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

func (d dispatcher) HandleRequest(ctx context.Context, evt Event) (EventResult, error) {
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

func (d dispatcher) handleRequestHelper(ctx context.Context, evt Event) (interface{}, error) {
	startEventTime := time.Now()
	defer func() {
		endTime := time.Now()
		startupDeltaMs := endTime.Sub(d.startupTime) / time.Millisecond
		eventDeltaMs := endTime.Sub(startEventTime) / time.Millisecond

		log.Printf("Startup: %dms, Event: %dms", startupDeltaMs, eventDeltaMs)
	}()

	log.Printf("source body is %q", evt.Body)

	var evtBody EventBody
	err := json.Unmarshal([]byte(evt.Body), &evtBody)
	if err != nil {
		return "", newErrf("unable to parse event: %s", err.Error())
	}

	handler, hasHandler := d.handlers[evtBody.Op]
	if !hasHandler {
		return "", newErrf("unknown op %q", evtBody.Op)
	}

	result, err := handler.Handle(ctx, evtBody.Args)
	if err != nil {
		return "", newErrf("error in %q handler: %s", evtBody.Op, err.Error())
	}

	return result, nil
}

func (d dispatcher) HandleTestRequest(eventStr string) {
	var evt Event
	err := json.Unmarshal([]byte(eventStr), &evt)
	if err != nil {
		log.Fatalf("error in test JSON: %s", err.Error())
	}

	result, err := d.HandleRequest(context.Background(), evt)
	if err != nil {
		log.Fatalf("error in test call: %s\n", err.Error())
	}
	log.Printf("%+v", result)
}

func (d dispatcher) HandleTestOpRequest(op string, payloadStr string) {
	evtBody := EventBody{
		op,
		json.RawMessage(payloadStr),
	}

	bodyJson, err := json.Marshal(evtBody)
	if err != nil {
		log.Fatalf("error in test call: %s\n", err.Error())
	}
	evt := Event{
		Body: string(bodyJson),
	}

	result, err := d.HandleRequest(context.Background(), evt)
	if err != nil {
		log.Fatalf("error in test call: %s\n", err.Error())
	}
	log.Printf("%+v", result)
}
