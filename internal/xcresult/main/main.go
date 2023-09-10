package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	buff, err := os.ReadFile(`/Users/home/go/src/github.com/qase-tms/qasectl/internal/xcresult/main/result.json`)
	if err != nil {
		panic(err)
	}

	var v map[string]any
	err = json.Unmarshal(buff, &v)
	if err != nil {
		panic(err)
	}

	decode(v)
}

type Decoder interface {
	TypeName() string
	Decode(map[string]any)
}

func decodeObject[T Decoder](m map[string]any) T {
	var t T

	if t.TypeName() != typeName(m) {
		panic("incorrectTypeName")
	}
	t.Decode(m)

	return t
}

func decodeArray[T Decoder](m map[string]any) []T {
	var result []T

	if typeName(m) != "Array" {
		panic("wrong type")
	}
	for _, v := range m["_values"].([]map[string]any) {
		obj := decodeObject[T](v)
		result = append(result, obj)
	}

	return result
}

func decode(m map[string]any) any {
	typeName := m["_type"].(map[string]any)["_name"].(string)
	fmt.Println("Type name is ", typeName)

	return nil
}

func typeName(m map[string]any) string {
	return m["_type"].(map[string]any)["_name"].(string)
}
