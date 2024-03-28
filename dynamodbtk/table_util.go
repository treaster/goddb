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
) (InputT, error) {
	itemAttributes := StructToAttributeMap(item)
	var oldItem InputT

	putResult, err := client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      itemAttributes,
	})
	if err != nil {
		fmt.Print("PutItem error %s", err.Error())
		return oldItem, err
	}

	AttributeMapToStruct(putResult.Attributes, &oldItem)
	return oldItem, nil
}

func QueryItemsByInt[OutputT any](
	ctx context.Context,
	client DynamoDBClient,
	tableName string,
	fieldName string,
	queryValue int) ([]OutputT, error) {
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
			fmt.Println("error:", err)
			continue
		}

		result = append(result, rowValues)
	}

	return result, nil
}
