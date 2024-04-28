package dynamodbtk

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func GetItem[InputT any, OutputT any](
	ctx context.Context,
	client DynamoDBClient,
	tableName string,
	rowKey InputT,
) (OutputT, error) {
	rowKeyMap := StructToAttributeMap(rowKey)
	var item OutputT

	getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       rowKeyMap,
	})
	if err != nil {
		return item, err
	}

	AttributeMapToStruct(getResult.Item, &item)
	return item, nil
}

func PutItem[InputT any](
	ctx context.Context,
	client DynamoDBClient,
	tableName string,
	item InputT,
	failIfExistsKey string,
) (InputT, error) {
	if tableName == "" {
		panic(fmt.Sprintf("error in PutItem: tableName is empty"))
	}
	itemAttributes := StructToAttributeMap(item)
	var oldItem InputT

	input := dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      itemAttributes,
	}

	if failIfExistsKey != "" {
		input.ConditionExpression = aws.String(fmt.Sprintf("attribute_not_exists(%s)", failIfExistsKey))
	}

	putResult, err := client.PutItem(ctx, &input)
	if err != nil {
		return oldItem, err
	}

	AttributeMapToStruct(putResult.Attributes, &oldItem)
	return oldItem, nil
}

func QueryItemsByIntField[OutputT any](
	ctx context.Context,
	client DynamoDBClient,
	tableName string,
	fieldName string,
	queryValue int) ([]OutputT, error) {
	if tableName == "" {
		panic(fmt.Sprintf("error in QueryItemsByIntField: tableName is empty"))
	}
	if fieldName == "" {
		panic(fmt.Sprintf("error in QueryItemsByIntField: fieldName is empty"))
	}

	queryValues := struct {
		Query int `ddb:":queryValue,N"`
	}{
		queryValue,
	}
	out, err := client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		KeyConditionExpression:    aws.String(fmt.Sprintf("%s = :queryValue", fieldName)),
		ExpressionAttributeValues: StructToAttributeMap(queryValues),
	})
	if err != nil {
		return nil, err
	}

	var result []OutputT
	for _, item := range out.Items {
		var rowValues OutputT
		err := AttributeMapToStruct(item, &rowValues)
		if err != nil {
			continue
		}

		result = append(result, rowValues)
	}

	return result, nil
}
