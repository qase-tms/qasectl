package xcresult

import (
	"fmt"
	"github.com/qase-tms/qasectl/pkg"
	"strconv"
	"strings"
)

type Decoder interface {
	TypeName() string
	Decode(map[string]any)
}

type DecoderFactory func(typ string) Decoder

type pt[T any] interface {
	Decoder
	*T
}

func DecodeVarArray(m map[string]any, factory DecoderFactory) []Decoder {
	var result []Decoder

	if typeName(m) != "Array" {
		panic("type is not array")
	}
	for _, v := range m["_values"].([]any) {
		obj := DecodeVarObject(v.(map[string]any), factory)
		result = append(result, obj)
	}

	return result
}

func DecodeVarObject(m map[string]any, factory DecoderFactory) Decoder {
	actualTypes := types(m)
	firstType := actualTypes[0]

	obj := factory(firstType)
	obj.Decode(m)

	return obj
}

func DecodeObject[T any, PT pt[T]](m map[string]any) T {
	var obj = PT(new(T))

	actualTypes := types(m)
	expectedType := obj.TypeName()
	if !pkg.Contains(actualTypes, expectedType) {
		actualTypes := strings.Join(actualTypes, ",")
		panic(fmt.Errorf("incorrectObjectType (actual=%s, expected=%q)", actualTypes, expectedType))
	}
	obj.Decode(m)

	return *obj
}

func DecodeArray[T any, PT pt[T]](m map[string]any) []T {
	var result []T

	if typeName(m) != "Array" {
		panic("type is not array")
	}
	for _, v := range m["_values"].([]any) {
		obj := DecodeObject[T, PT](v.(map[string]any))
		result = append(result, obj)
	}

	return result
}

func DecodeDouble(m map[string]any) float64 {
	if typeName(m) != "Double" {
		panic(fmt.Errorf("type is not double, but %v", typeName(m)))
	}

	strValue := m["_value"].(string)
	value, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		panic(err)
	}

	return value
}

func DecodeString(m map[string]any) string {
	if typeName(m) != "String" {
		panic(fmt.Errorf("type is not string, but %v", typeName(m)))
	}

	return m["_value"].(string)
}

func typeName(m map[string]any) string {
	return m["_type"].(map[string]any)["_name"].(string)
}

func types(m map[string]any) []string {
	var result []string
	_type := m["_type"].(map[string]any)

	for _type != nil {
		name := _type["_name"].(string)
		result = append(result, name)

		if v, ok := _type["_supertype"].(map[string]any); ok {
			_type = v
		} else {
			_type = nil
		}
	}

	return result
}
