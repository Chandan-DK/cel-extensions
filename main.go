package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

func main() {
	env, err := cel.NewEnv(
		cel.Function("json_parse",
			cel.Overload("json_parse_string",
				[]*cel.Type{cel.StringType},
				cel.DynType,
				cel.UnaryBinding(jsonParseString),
			),
		),
	)
	if err != nil {
		log.Fatalf("environment creation error: %v\n", err)
	}

	helperFuncToEvaluateExpression(env, `json_parse('[{\"cpu\":512,\"essential\":true,\"image\":\"nginx:1.23.1\",\"memory\":2048,\"name\":\"foo-task\",\"portMappings\":[{\"containerPort\":80,\"hostPort\":80}]}]')`)
	helperFuncToEvaluateExpression(env, `json_parse('{"cpu":512,"essential":true,"image":"nginx:1.23.1"}')`)
	helperFuncToEvaluateExpression(env, `json_parse('[{"cpu":512,"essential":true,"image":"nginx:1.23.1"}]')`)
	helperFuncToEvaluateExpression(env, `json_parse('{"task":{"cpu":512,"memory":2048}}')`)
	helperFuncToEvaluateExpression(env, `json_parse('[{"name":"task1"},{"name":"task2"}]')`)
	helperFuncToEvaluateExpression(env, `json_parse('"Hello, world!"')`)
	helperFuncToEvaluateExpression(env, `json_parse('42')`)
	helperFuncToEvaluateExpression(env, `json_parse('true')`)
	helperFuncToEvaluateExpression(env, `json_parse('{}')`)
	helperFuncToEvaluateExpression(env, `json_parse('[]')`)
	// helperFuncToEvaluateExpression(env, `json_parse('{"cpu":512,')`)
	helperFuncToEvaluateExpression(env, `json_parse('[42, "text", {"key": "value"}, true]')`)
	helperFuncToEvaluateExpression(env, `json_parse('null')`)
	helperFuncToEvaluateExpression(env, `json_parse('"Hello, \\"world\\"!"')`)
	helperFuncToEvaluateExpression(env, `json_parse('{"key":"value!@#$%^&*()_+-=<>?/"}')`)
	helperFuncToEvaluateExpression(env, `json_parse('{"largeArray":[1,2,3,4,5,6,7,8,9,10],"nested":{"field":"value","anotherField":[100,200,300]}}')`)
	// helperFuncToEvaluateExpression(env, `json_parse('{"cpu":512,}')`)
	helperFuncToEvaluateExpression(env, `json_parse('"true"')`)
	helperFuncToEvaluateExpression(env, `json_parse('   { "cpu" : 512 , "essential" : true }   ')`)
}

func jsonParseString(val ref.Val) ref.Val {
	jsonString := val.Value().(string)

	var parsedJSON any
	if err := json.Unmarshal([]byte(jsonString), &parsedJSON); err != nil {
		return types.NewErr("error while parsing json %v", err)
	}

	return convertToCelValue(parsedJSON)
}

func convertToCelValue(value any) ref.Val {
	switch v := value.(type) {
	case map[string]any:
		return types.NewDynamicMap(types.DefaultTypeAdapter, v)
	case []any:
		return types.NewDynamicList(types.DefaultTypeAdapter, v)
	case string:
		return types.String(v)
	case float64:
		return types.Double(v)
	case bool:
		return types.Bool(v)
	case nil:
		return types.NullValue
	default:
		return types.NewErr("unsupported JSON type %T", v)
	}
}

func helperFuncToEvaluateExpression(env *cel.Env, expression string) {
	ast, iss := env.Compile(expression)
	if iss.Err() != nil {
		log.Fatalf("error during compilation: %v\n", iss.Err())
	}

	prg, err := env.Program(ast)
	if err != nil {
		log.Fatalf("program creation error: %v\n", err)
	}

	out, details, err := prg.Eval(cel.NoVars())
	if err != nil {
		log.Fatalf("evaluation error: %v\n", err)
	}

	fmt.Printf("out: %v\ndetails: %v\n", out, details)
}
