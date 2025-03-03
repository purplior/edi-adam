package repository

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/purplior/edi-adam/domain/shared/dto/pagination"
	"github.com/purplior/edi-adam/domain/shared/inner"
	"github.com/purplior/edi-adam/infra/database/dynamo"
)

type (
	dynamoRepository[T any, Q any] struct {
		client *dynamo.Client
		// 사용 대상 테이블 이름
		tableName string
		// 쿼리 옵션(Q)을 받아 DynamoDB QueryInput을 생성합니다.
		buildQueryInput func(opt Q) (*dynamodb.QueryInput, error)
		// 쿼리 옵션(Q)을 받아 단일 아이템 조작에 필요한 키를 생성합니다.
		buildKey func(opt Q) (map[string]types.AttributeValue, error)
		// 업데이트할 구조체(T)를 업데이트 표현식과 attribute values로 변환합니다.
		buildUpdateExpression func(m T) (string, map[string]types.AttributeValue, error)
	}
)

func (r *dynamoRepository[T, Q]) Read(
	session inner.Session,
	queryOption Q,
) (
	m T,
	err error,
) {
	key, err := r.buildKey(queryOption)
	if err != nil {
		return m, err
	}
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key:       key,
	}
	output, err := r.client.GetItem(session.Context(), input)
	if err != nil {
		return m, err
	}
	if output.Item == nil {
		return m, errors.New("item not found")
	}
	err = attributevalue.UnmarshalMap(output.Item, &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

func (r *dynamoRepository[T, Q]) ReadCount(
	session inner.Session,
	queryOption Q,
) (int, error) {
	qInput, err := r.buildQueryInput(queryOption)
	if err != nil {
		return 0, err
	}
	// COUNT 쿼리 설정
	qInput.Select = types.SelectCount
	output, err := r.client.Query(session.Context(), qInput)
	if err != nil {
		return 0, err
	}
	return int(output.Count), nil
}

func (r *dynamoRepository[T, Q]) ReadList(
	session inner.Session,
	queryOption Q,
) (
	mArr []T,
	err error,
) {
	qInput, err := r.buildQueryInput(queryOption)
	if err != nil {
		return nil, err
	}
	output, err := r.client.Query(session.Context(), qInput)
	if err != nil {
		return nil, err
	}
	err = attributevalue.UnmarshalListOfMaps(output.Items, &mArr)
	if err != nil {
		return nil, err
	}
	return mArr, nil
}

// 커서 기반의 페이지네이션을 사용함.
func (r *dynamoRepository[T, Q]) ReadPaginatedList(
	session inner.Session,
	option pagination.PaginationQuery[Q],
) (
	mArr []T,
	meta pagination.PaginationMeta,
	err error,
) {
	qInput, err := r.buildQueryInput(option.QueryOption)
	if err != nil {
		return nil, meta, err
	}

	qInput.Limit = aws.Int32(int32(option.PageRequest.Size))
	if option.PageRequest.Cursor != "" {
		lastKey, err := decodeCursor(option.PageRequest.Cursor)
		if err != nil {
			return nil, meta, err
		}
		qInput.ExclusiveStartKey = lastKey
	}

	output, err := r.client.Query(session.Context(), qInput)
	if err != nil {
		return nil, meta, err
	}

	if err = attributevalue.UnmarshalListOfMaps(output.Items, &mArr); err != nil {
		return nil, meta, err
	}

	nextCursor := ""
	if output != nil && output.LastEvaluatedKey != nil && len(output.LastEvaluatedKey) > 0 {
		nextCursor, err = encodeCursor(output.LastEvaluatedKey)
		if err != nil {
			return nil, meta, err
		}
	}

	meta = pagination.PaginationMeta{
		Size:       option.PageRequest.Size,
		NextCursor: nextCursor,
	}
	return mArr, meta, nil
}

func (r *dynamoRepository[T, Q]) Create(
	session inner.Session,
	m T,
) (
	mRet T,
	err error,
) {
	item, err := attributevalue.MarshalMap(m)
	if err != nil {
		return m, err
	}
	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	}
	_, err = r.client.PutItem(session.Context(), input)
	if err != nil {
		return m, err
	}
	mRet = m
	return mRet, nil
}

func (r *dynamoRepository[T, Q]) Updates(
	session inner.Session,
	queryOption Q,
	m T,
) (
	err error,
) {
	key, err := r.buildKey(queryOption)
	if err != nil {
		return err
	}
	updateExpr, exprAttrValues, err := r.buildUpdateExpression(m)
	if err != nil {
		return err
	}
	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(r.tableName),
		Key:                       key,
		UpdateExpression:          aws.String(updateExpr),
		ExpressionAttributeValues: exprAttrValues,
		ReturnValues:              types.ReturnValueUpdatedNew,
	}
	_, err = r.client.UpdateItem(session.Context(), input)
	if err != nil {
		return err
	}
	return nil
}

func (r *dynamoRepository[T, Q]) Delete(
	session inner.Session,
	queryOption Q,
) (
	err error,
) {
	key, err := r.buildKey(queryOption)
	if err != nil {
		return err
	}
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key:       key,
	}
	_, err = r.client.DeleteItem(session.Context(), input)
	if err != nil {
		return err
	}
	return nil
}

func encodeCursor(key map[string]types.AttributeValue) (string, error) {
	// JSON 직렬화 (AttributeValue 인터페이스가 바로 직렬화되지 않을 수 있으므로, 필요한 경우 별도 변환이 필요합니다)
	bytes, err := json.Marshal(key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

// decodeCursor는 base64로 인코딩된 커서를 디코딩하여 LastEvaluatedKey 형식으로 변환합니다.
func decodeCursor(cursor string) (map[string]types.AttributeValue, error) {
	bytes, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, err
	}
	var key map[string]types.AttributeValue
	err = json.Unmarshal(bytes, &key)
	if err != nil {
		return nil, err
	}
	return key, nil
}
