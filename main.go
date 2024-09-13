package main

import (
	"fmt"
	"jsonparsemod/library"
	"log"

	"github.com/google/cel-go/cel"
)

func main() {
	env, err := cel.NewEnv(
		library.JsonParseLib(),
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
