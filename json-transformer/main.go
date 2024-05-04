package main

import (
	"encoding/json"
	"io/fs"
	"log"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"
)

// Represent a Json
type Json map[string]interface{}

// Represent a Json array
type JsonArray []Json

// Constant to tag value to be omitted
const OMIT string = "_OMIT"

/*
Constant listing the supported data types, as follows:
  - S: string
  - N: number
  - M: map
  - BOOL: boolean
  - NULL: null
  - L: list
*/
var DATA_TYPES = []string{"N", "S", "M", "BOOL", "NULL", "L"}

func main() {
	var filepath string = "./input.json"
	const outputFilename string = "output.json"
	const outputFilePermission fs.FileMode = 0666
	var inputJson Json
	output := make(JsonArray, 1)
	output[0] = make(Json)

	if len(os.Args) > 1 {
		filepath = os.Args[1]
	}

	content, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	err = json.Unmarshal(content, &inputJson)
	if err != nil {
		log.Panic(err)
	}

	it := reflect.ValueOf(inputJson).MapRange()
	for it.Next() {
		k := sanitizeString(it.Key().Interface().(string))
		v := it.Value().Interface()

		res := dfs(v, k)

		output[0][k] = res

		omit(res, k, output[0])
	}

	out, err := json.Marshal(output)
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(outputFilename, out, outputFilePermission); err != nil {
		log.Fatal(err)
	}
}

// Recursive Depth-First Search: traverse json
func dfs(valL1 any, keyL1 string) any {

	switch reflect.TypeOf(valL1).Kind() {
	case reflect.Map:
		var out any
		mapRes := make(Json)

		it := reflect.ValueOf(valL1).MapRange()
		for it.Next() {
			keyL2 := sanitizeString(it.Key().Interface().(string))
			valL2 := it.Value().Interface()

			res := dfs(valL2, keyL2)

			if slices.Contains(DATA_TYPES, keyL2) != true {
				mapRes[keyL2] = res
				omit(res, keyL2, mapRes)
				out = mapRes
			} else {
				out = res
			}
		}

		return out

	case reflect.Slice:
		return transformL(valL1.([]interface{}))
	case reflect.String:
		str := sanitizeString(valL1.(string))

		switch keyL1 {
		case "S":
			return transformS(str)
		case "N":
			return transformN(str)
		case "BOOL":
			return transformBOOL(str)
		case "NULL":
			return transformNULL(str)
		default:
			return OMIT
		}
	}

	return OMIT
}

/*
Description:
  - Remove trailing and leading space
*/
func sanitizeString(input string) string {
	return strings.Trim(input, " ")
}

/*
Description:
  - Remove leading "0"
  - Convert `string` to `number`
*/
func sanitizeNumber(input string) (float64, error) {
	return strconv.ParseFloat(strings.TrimLeft(input, "0"), 64)
}

/*
Description:
  - Convert `string` to `number`
  - Tag result to be omitted if invalid `number`
*/
func transformN(input string) any {
	res, err := sanitizeNumber(input)
	if err != nil {
		return OMIT
	}
	return res
}

/*
Description:
  - Tag result to be omitted if empty `string`
  - Convert `datetime` to Unix Epoch
  - By default return `string`
*/
func transformS(input string) any {
	if input == "" {
		return OMIT
	}

	res, err := time.Parse(time.RFC3339, input)
	if err != nil {
		return input
	}
	return res.Unix()
}

/*
Description:
  - Convert `string` to `boolean`
  - Tag result to be omitted if invalid `boolean`
*/
func transformBOOL(input string) any {
	res, err := strconv.ParseBool(input)
	if err != nil {
		return OMIT
	}
	return res
}

/*
Description:
  - Tag result to be omitted if invalid boolean
  - Return `null` if `boolean` is `true` else `nil`
*/
func transformNULL(input string) any {
	res, err := strconv.ParseBool(input)
	if err != nil {
		return OMIT
	}

	if res {

		return nil
	}
	return OMIT
}

/*
Description:
  - Tag result to be omitted if unsupported data types
  - Convert Json array to list
*/
func transformL(input []interface{}) any {
	out := make([]any, 0)

	for idx, item := range input {
		res := dfs(item, string(idx))

		if res != OMIT {
			out = append(out, res)
		}
	}

	if len(out) == 0 {
		return OMIT
	}

	return out
}

// Remove key from output if input tagged to be omitted
func omit(input any, key string, output Json) {
	if input == OMIT || key == "" {
		delete(output, key)
	}
}
