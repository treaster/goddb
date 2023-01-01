package dynamodbtk_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/require"
	"github.com/treaster/goawstk/dynamodbtk"
)

type MyStruct struct {
	Field1  int    `ddb:"field_1,N"`
	Field2  int    `ddb:"-,N"`
	Field3  int    `ddb:",N"`
	Field4  string `ddb:",S"`
	Field5  string `ddb:"field_5,S"`
	Field6  string `ddb:"field_6,"`
	Field7  string
	Field8  int    `ddb:" field_8 , N "`
	Field9  string `ddb:"-"`          // - means ignore
	Field10 string `ddb:"f,S,"`       // too many tag fields
	Field11 string `ddb:"f,S,xyz"`    // too many tag fields
	Field12 string `ddb:",Q"`         // illegal type code
	Field13 bool   `ddb:"field_13,N"` // bool type
	Field14 bool   `ddb:"field_14,N"` // another bool type
}

func TestStructToAttributeMap(t *testing.T) {
	input := MyStruct{1, 2, 3, "s4", "s5", "s6", "s7", 8, "s9", "s10", "s11", "s12", true, false}
	am := dynamodbtk.StructToAttributeMap(input)

	expectedAm := map[string]types.AttributeValue{
		"field_1": &types.AttributeValueMemberN{Value: "1"},
		"Field2":  &types.AttributeValueMemberN{Value: "2"},
		"Field3":  &types.AttributeValueMemberN{Value: "3"},
		"Field4":  &types.AttributeValueMemberS{Value: "s4"},
		"field_5": &types.AttributeValueMemberS{Value: "s5"},
		"field_6": &types.AttributeValueMemberS{Value: "s6"},
		"Field7":  &types.AttributeValueMemberS{Value: "s7"},
		"field_8": &types.AttributeValueMemberN{Value: "8"},
		// Field9 and onwards are omitted
		"field_13": &types.AttributeValueMemberN{Value: "1"},
		"field_14": &types.AttributeValueMemberN{Value: "0"},
	}

	require.Equal(t, expectedAm, am)

	var output MyStruct

	err := dynamodbtk.AttributeMapToStruct(am, &output)
	require.NoError(t, err)

	input.Field9 = ""
	input.Field10 = ""
	input.Field11 = ""
	input.Field12 = ""
	require.Equal(t, input, output)
}
