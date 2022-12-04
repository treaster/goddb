package goddb_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/require"
	"github.com/treaster/goddb"
)

type MyStruct struct {
	IntField int    `ddb:"N"`
	StrField string `ddb:"S"`
}

func TestStructToAttributeMap(t *testing.T) {
	input := MyStruct{10, "abc"}
	am := goddb.StructToAttributeMap(input)

	expectedAm := map[string]types.AttributeValue{
		"IntField": &types.AttributeValueMemberN{Value: "10"},
		"StrField": &types.AttributeValueMemberS{Value: "abc"},
	}
	require.Equal(t, expectedAm, am)

	var output MyStruct

	goddb.AttributeMapToStruct(am, &output)
	require.Equal(t, input, output)
}
