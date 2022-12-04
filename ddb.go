package goddb

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBClient interface {
	GetItem(context.Context, *dynamodb.GetItemInput, ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	ListTables(context.Context, *dynamodb.ListTablesInput, ...func(*dynamodb.Options)) (*dynamodb.ListTablesOutput, error)
	Scan(context.Context, *dynamodb.ScanInput, ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
	UpdateItem(context.Context, *dynamodb.UpdateItemInput, ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
}

func StructToAttributeMap(input interface{}) map[string]types.AttributeValue {
	inputValue := reflect.ValueOf(input)
	inputType := inputValue.Type()

	am := map[string]types.AttributeValue{}

	for i := 0; i < inputType.NumField(); i++ {
		f := inputType.Field(i)
		fieldValue := inputValue.Field(i)

		var strValue string
		switch {
		case fieldValue.CanInt():
			strValue = strconv.FormatInt(fieldValue.Int(), 10)
		case fieldValue.CanUint():
			strValue = strconv.FormatUint(fieldValue.Uint(), 10)
		case fieldValue.CanFloat():
			strValue = strconv.FormatFloat(fieldValue.Float(), 'E', -1, 64)
		case fieldValue.Kind() == reflect.String:
			strValue = fieldValue.String()
		default:
			panic(fmt.Sprintf("unsupported kind %v for field %q", fieldValue.Kind(), f.Name))
		}

		field := inputType.Field(i)
		ddbType := field.Tag.Get("ddb")
		var v types.AttributeValue
		switch ddbType {
		case "N":
			v = &types.AttributeValueMemberN{Value: strValue}

		case "S":
			v = &types.AttributeValueMemberS{Value: strValue}

		}

		am[f.Name] = v
	}

	return am
}

func AttributeMapToStruct(am map[string]types.AttributeValue, out interface{}) error {
	outPtrValue := reflect.ValueOf(out)
	if outPtrValue.Kind() != reflect.Pointer {
		panic(fmt.Sprintf("output value must be a pointer-to-struct, but got %v", outPtrValue.Kind()))
	}
	outValue := outPtrValue.Elem()
	if outValue.Kind() != reflect.Struct {
		panic(fmt.Sprintf("output value must be a pointer-to-struct, but got %v", outPtrValue.Kind()))
	}

	outType := outValue.Type()
	for colName, ddbValue := range am {
		fieldSpec, hasField := outType.FieldByName(colName)

		if !hasField {
			panic(fmt.Sprintf("no struct field matching ddb column name %s", colName))
		}

		ddbType := fieldSpec.Tag.Get("ddb")
		var strValue string
		switch ddbType {
		case "N":
			strValue = ddbValue.(*types.AttributeValueMemberN).Value

		case "S":
			strValue = ddbValue.(*types.AttributeValueMemberS).Value
		}

		fieldValue := outValue.FieldByName(colName)

		switch {
		case fieldValue.CanInt():
			finalValue, err := strconv.ParseInt(strValue, 10, 64)
			if err != nil {
				return fmt.Errorf("unable to parse int value from %q for field %q", strValue, colName)
			}
			if fieldValue.OverflowInt(finalValue) {
				return fmt.Errorf("int value %d is too big for field %q type %v", finalValue, colName, fieldValue.Kind())
			}
			fieldValue.SetInt(finalValue)
		case fieldValue.CanUint():
			finalValue, err := strconv.ParseUint(strValue, 10, 64)
			if err != nil {
				return fmt.Errorf("unable to parse uint value from %q for field %q", strValue, colName)
			}
			if fieldValue.OverflowUint(finalValue) {
				return fmt.Errorf("uint value %d is too big for field %q type %v", finalValue, colName, fieldValue.Kind())
			}
			fieldValue.SetUint(finalValue)
		case fieldValue.CanFloat():
			finalValue, err := strconv.ParseFloat(strValue, 64)
			if err != nil {
				return fmt.Errorf("unable to parse float value from %q for field %q", strValue, colName)
			}
			if fieldValue.OverflowFloat(finalValue) {
				return fmt.Errorf("float value %f is too big for field %q type %v", finalValue, colName, fieldValue.Kind())
			}
			fieldValue.SetFloat(finalValue)
		case fieldValue.Kind() == reflect.String:
			finalValue := strValue
			fieldValue.SetString(finalValue)
		default:
			panic(fmt.Sprintf("unsupported kind %v for field %q", fieldValue.Kind(), colName))
		}
	}

	return nil
}
