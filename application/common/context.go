package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/purplior/edi-adam/application/config"
	"github.com/purplior/edi-adam/application/response"
	"github.com/purplior/edi-adam/domain/shared/constant"
	"github.com/purplior/edi-adam/domain/shared/exception"
	"github.com/purplior/edi-adam/domain/shared/inner"
	"github.com/purplior/edi-adam/domain/shared/logger"
)

type (
	Context struct {
		echo.Context
		Identity *inner.Identity
		session  inner.Session
	}
)

// @METHOD {세션 정보를 가져옴}
func (ctx *Context) Session() inner.Session {
	return ctx.session
}

// @METHOD {세션 정보를 삭제함}
func (ctx *Context) ClearSession() {
	if ctx.session != nil {
		ctx.session.ReleaseTransaction()
		ctx.session = nil
	}
}

// @METHOD {Param, Query, Body(JSON)을 구조체 기반으로 바로 파싱해오는 함수}
// @see Streaming Body와 같이 Body가 매우 큰 경우 사용하지 마세요
// [태그 규칙] 첫번째에 키 이름, 두번째부터 옵션
// - Parameter: `param:"parameter_name,required"`
// - QueryString: `query:"query_name,required"`
// - JSON Body: `body:"json_key_name,required"`
func (ctx *Context) BindRequestData(i interface{}) error {
	if err := bindRequestData(ctx, i); err != nil {
		if config.Phase() == constant.Phase_Production {
			return exception.ErrBadRequest
		}
		return err
	}

	return nil
}

// @METHOD {앱 응답구조에 맞춘 JSON 응답}
func (ctx Context) SendJSON(jsonResponse response.JSONResponse) error {
	if jsonResponse.Status == 0 {
		jsonResponse.Status = response.Status_Ok
	}
	return ctx.JSON(
		jsonResponse.Status,
		jsonResponse,
	)
}

// @METHOD {앱 응답구조에 맞춘 OK 응답}
func (ctx *Context) SendOK(data interface{}) error {
	return ctx.SendJSON(response.JSONResponse{
		Data: data,
	})
}

// @METHOD {앱 응답구조에 맞춘 Created 응답}
func (ctx *Context) SendCreated(data interface{}) error {
	return ctx.SendJSON(response.JSONResponse{
		Status: response.Status_Created,
		Data:   data,
	})
}

// @METHOD {앱 응답구조에 맞춘 오류 응답}
func (ctx *Context) SendError(err error) error {
	statusCode, message := getResponseOfError(err)
	if statusCode == response.Status_InternalServerError && config.Phase() != constant.Phase_Production {
		logger.Info(fmt.Sprintf("오류응답: %s", err.Error()))
		logger.Info(fmt.Sprintf("Stack Trace:\n%s", debug.Stack()))
	}

	return ctx.JSON(statusCode, response.ErrorResponse{
		Status:  statusCode,
		Message: message,
	})
}

// @METHOD {앱 응답구조에 맞춘 커스텀 오류 응답}
func (ctx Context) SendCustomError(
	statusCode int,
	errMessage string,
) error {
	return ctx.JSON(statusCode, response.ErrorResponse{
		Status:  statusCode,
		Message: errMessage,
	})
}

// @METHOD {앱 응답구조에 맞춘 인증오류 응답}
func (ctx *Context) SendUnauthorized() error {
	return ctx.SendError(exception.ErrUnauthorized)
}

// @FUNC {인터페이스 구조체의 reflect 정보를 읽어서, 요청 데이터를 찾아서 매핑}
func bindRequestData(c *Context, i interface{}) error {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return errors.New("BindAll requires a non-nil pointer to a struct")
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return errors.New("BindAll requires a pointer to a struct")
	}

	hasBodyField := containsBodyField(v)
	var bodyData map[string]interface{}
	if hasBodyField {
		bodyBytes, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return fmt.Errorf("failed to read request body: %w", err)
		}
		if err := json.Unmarshal(bodyBytes, &bodyData); err != nil {
			return fmt.Errorf("failed to unmarshal request body: %w", err)
		}
	}

	return bindStructFields(v, c, bodyData)
}

func containsBodyField(v reflect.Value) bool {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			if containsBodyField(v.Field(i)) {
				return true
			}
		} else if bodyTag := field.Tag.Get("body"); bodyTag != "" {
			return true
		}
	}
	return false
}

func bindStructFields(v reflect.Value, c *Context, bodyData map[string]interface{}) error {
	t := v.Type()
	// 태그 값을 "키,옵션" 형식으로 파싱하는 헬퍼
	parseTag := func(tag string) (key string, required bool) {
		parts := strings.Split(tag, ",")
		key = parts[0]
		required = false
		for _, opt := range parts[1:] {
			if strings.TrimSpace(opt) == "required" {
				required = true
				break
			}
		}
		return
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldVal := v.Field(i)

		// 익명 필드(임베디드 구조체)의 경우 재귀 호출
		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			if err := bindStructFields(fieldVal, c, bodyData); err != nil {
				return err
			}
			continue
		}

		// 컨텍스트(c)가 nil이 아니면 param, query 태그를 처리 (중첩 구조체에서는 nil을 전달함)
		if c != nil {
			// URL 파라미터 처리
			if paramTag := field.Tag.Get("param"); paramTag != "" {
				key, required := parseTag(paramTag)
				value := c.Param(key)
				if required && value == "" {
					return exception.ErrBadRequest
				}
				if err := setFieldValue(fieldVal, value, key, "param"); err != nil {
					return err
				}
			}

			// Query 파라미터 처리
			if queryTag := field.Tag.Get("query"); queryTag != "" {
				key, required := parseTag(queryTag)
				value := c.QueryParam(key)
				if required && value == "" {
					return exception.ErrBadRequest
				}
				if err := setFieldValue(fieldVal, value, key, "query"); err != nil {
					return err
				}
			}
		}

		// JSON Body 처리
		if bodyTag := field.Tag.Get("body"); bodyTag != "" {
			key, required := parseTag(bodyTag)
			if bodyData == nil {
				if required {
					return exception.ErrBadRequest
				}
				continue
			}
			raw, exists := bodyData[key]
			if !exists {
				if required {
					return exception.ErrBadRequest
				}
				continue
			}
			if err := setFieldValueFromInterface(fieldVal, raw, key, "body"); err != nil {
				return err
			}
		}
	}
	return nil
}

func setFieldValue(fieldVal reflect.Value, value string, key, source string) error {
	switch fieldVal.Kind() {
	case reflect.String:
		fieldVal.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value == "" {
			fieldVal.SetInt(0)
		} else {
			n, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return fmt.Errorf("failed to convert %s field %s to int: %w", source, key, err)
			}
			fieldVal.SetInt(n)
		}
	case reflect.Float32, reflect.Float64:
		if value == "" {
			fieldVal.SetFloat(0)
		} else {
			f, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("failed to convert %s field %s to float: %w", source, key, err)
			}
			fieldVal.SetFloat(f)
		}
	case reflect.Bool:
		if value == "" {
			fieldVal.SetBool(false)
		} else {
			b, err := strconv.ParseBool(value)
			if err != nil {
				return fmt.Errorf("failed to convert %s field %s to bool: %w", source, key, err)
			}
			fieldVal.SetBool(b)
		}
	default:
		return fmt.Errorf("unsupported type for %s field %s", source, key)
	}
	return nil
}

func setFieldValueFromInterface(fieldVal reflect.Value, raw interface{}, key, source string) error {
	switch fieldVal.Kind() {
	case reflect.String:
		switch v := raw.(type) {
		case string:
			fieldVal.SetString(v)
		default:
			fieldVal.SetString(fmt.Sprintf("%v", v))
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch v := raw.(type) {
		case float64:
			fieldVal.SetInt(int64(v))
		case string:
			n, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return fmt.Errorf("failed to convert %s field %s to int: %w", source, key, err)
			}
			fieldVal.SetInt(n)
		default:
			return fmt.Errorf("unsupported type for %s field %s", source, key)
		}
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch v := raw.(type) {
		case float64:
			fieldVal.SetUint(uint64(v))
		case string:
			n, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return fmt.Errorf("failed to convert %s field %s to uint: %w", source, key, err)
			}
			fieldVal.SetUint(n)
		default:
			return fmt.Errorf("unsupported type for %s field %s", source, key)
		}
	case reflect.Float32, reflect.Float64:
		switch v := raw.(type) {
		case float64:
			fieldVal.SetFloat(v)
		case string:
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return fmt.Errorf("failed to convert %s field %s to float: %w", source, key, err)
			}
			fieldVal.SetFloat(f)
		default:
			return fmt.Errorf("unsupported type for %s field %s", source, key)
		}
	case reflect.Bool:
		switch v := raw.(type) {
		case bool:
			fieldVal.SetBool(v)
		case string:
			b, err := strconv.ParseBool(v)
			if err != nil {
				return fmt.Errorf("failed to convert %s field %s to bool: %w", source, key, err)
			}
			fieldVal.SetBool(b)
		default:
			return fmt.Errorf("unsupported type for %s field %s", source, key)
		}
	case reflect.Struct:
		// 중첩 구조체의 경우, raw가 map[string]interface{}여야 합니다.
		nestedMap, ok := raw.(map[string]interface{})
		if !ok {
			return fmt.Errorf("failed to convert %s field %s: expected object", source, key)
		}
		// 중첩 구조체는 body 데이터이므로, 컨텍스트는 nil로 전달합니다.
		return bindStructFields(fieldVal, nil, nestedMap)
	default:
		return fmt.Errorf("unsupported field type for %s field %s", source, key)
	}
	return nil
}
