package lambdatk_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/treaster/goawstk/lambdatk"
)

type AdderInput struct {
	Input int `json:"input"`
}
type AdderOutput struct {
	Output int `json:"output"`
}

type AdderHandler struct {
	AdderValue int
}

func (h AdderHandler) Handle(ctx context.Context, args json.RawMessage) (interface{}, error) {
	var input AdderInput
	err := json.Unmarshal(args, &input)
	if err != nil {
		return nil, err
	}

	output := AdderOutput{
		h.AdderValue + input.Input,
	}
	return output, nil
}

func TestDispatcher(t *testing.T) {
	require.True(t, true)

	handlers := map[string]lambdatk.OpHandler{
		"add1": AdderHandler{1},
		"add5": AdderHandler{5},
	}

	dispatcher := lambdatk.MakeDispatcher(handlers)

	evt := lambdatk.Event{
		Body: `{
			"op": "add1",
			"args": {"input": 2}
		}`,
	}
	result, err := dispatcher.HandleRequest(context.TODO(), evt)
	require.NoError(t, err)

	expected := lambdatk.EventResult{
		Result: AdderOutput{
			Output: 3,
		},
		Error: "",
	}
	require.Equal(t, expected, result)
}
