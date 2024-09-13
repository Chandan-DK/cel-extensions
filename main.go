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
		cel.Function("json_parse", // Name by which the function will be called in the CEL expression Example: In this case it would be json_parse(`{}`)
			// Overload gives us a way to define a function that can be used like this json_parse('{"key": "value"}') instead of <string>.json_parse() which would need to use MemberOverload()
			cel.Overload("json_parse_string", // ID of the Overload (https://pkg.go.dev/github.com/google/cel-go@v0.21.0/common/decls#OverloadDecl.ID)
				[]*cel.Type{cel.StringType}, // The argument type our custom function will accept
				cel.DynType,                 // Return type of our custom function
				// UnaryBinding takes a function of type UnaryOp which is basically a function that takes a `single`` value and produces an output. (https://pkg.go.dev/github.com/google/cel-go@v0.21.0/common/functions#UnaryOp)
				// Likewise you also have BinaryBinding and FunctionBinding (https://pkg.go.dev/github.com/google/cel-go@v0.21.0/common/functions#pkg-types)
				cel.UnaryBinding(jsonParseString),
			),
		),
	)
	if err != nil {
		log.Fatalf("Environment creation error: %v", err)
	}

	testCases := []string{
		`json_parse('[{\"cpu\":512,\"essential\":true,\"image\":\"nginx:1.23.1\",\"memory\":2048,\"name\":\"foo-task\",\"portMappings\":[{\"containerPort\":80,\"hostPort\":80}]}]')`,
		`json_parse('{"cpu":512,"essential":true,"image":"nginx:1.23.1"}')`,
		`json_parse('[{"cpu":512,"essential":true,"image":"nginx:1.23.1"}]')`,
		`json_parse('{"task":{"cpu":512,"memory":2048}}')`,
		`json_parse('[{"name":"task1"},{"name":"task2"}]')`,
		`json_parse('"Hello, world!"')`,
		`json_parse('42')`,
		`json_parse('true')`,
		`json_parse('{}')`,
		`json_parse('[]')`,
		`json_parse('[42, "text", {"key": "value"}, true]')`,
		`json_parse('null')`,
		`json_parse('"Hello, \\"world\\"!"')`,
		`json_parse('{"key":"value!@#$%^&*()_+-=<>?/"}')`,
		`json_parse('{"largeArray":[1,2,3,4,5,6,7,8,9,10],"nested":{"field":"value","anotherField":[100,200,300]}}')`,
		`json_parse('"true"')`,
		`json_parse('   { "cpu" : 512 , "essential" : true }   ')`,
	}

	for _, expression := range testCases {
		helperFuncToEvaluateExpression(env, expression)
	}
}

func jsonParseString(val ref.Val) ref.Val {
	jsonString, ok := val.Value().(string)
	if !ok {
		return types.NewErr("expected a string, got %T", val.Type())
	}

	var parsedJSON any
	if err := json.Unmarshal([]byte(jsonString), &parsedJSON); err != nil {
		return types.NewErr("error while parsing JSON: %v", err)
	}

	return convertToCelValue(parsedJSON)
}

func convertToCelValue(value any) ref.Val {
	fmt.Printf("Value Type: %T", value)
	switch v := value.(type) {
	case map[string]any:
		return types.NewStringInterfaceMap(types.DefaultTypeAdapter, v)
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
		log.Fatalf("Compilation error: %v", iss.Err())
	}

	prg, err := env.Program(ast)
	if err != nil {
		log.Fatalf("Program creation error: %v", err)
	}

	out, _, err := prg.Eval(cel.NoVars()) // We pass values for variables that are declared in the environment, using Eval() Example: https://github.com/google/cel-go/tree/master/examples#simple-example-using-builtin-operators
	if err != nil {
		log.Fatalf("Evaluation error: %v", err)
	}

	fmt.Printf("\nOutput: %v\n\n", out)
}
