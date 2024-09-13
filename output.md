# Output

Running `go run .` gives the following output currently:

```bash
Value Type: []interface {}
Output: [map[cpu:512 essential:true image:nginx:1.23.1 memory:2048 name:foo-task portMappings:[map[containerPort:80 hostPort:80]]]]

Value Type: map[string]interface {}
Output: {image: nginx:1.23.1, cpu: 512, essential: true}

Value Type: []interface {}
Output: [map[cpu:512 essential:true image:nginx:1.23.1]]

Value Type: map[string]interface {}
Output: {task: {cpu: 512, memory: 2048}}

Value Type: []interface {}
Output: [map[name:task1], map[name:task2]]

Value Type: string
Output: Hello, world!

Value Type: float64
Output: 42

Value Type: bool
Output: true

Value Type: map[string]interface {}
Output: {}

Value Type: []interface {}
Output: []

Value Type: []interface {}
Output: [42, text, map[key:value], true]

Value Type: <nil>
Output: 0

Value Type: string
Output: Hello, "world"!

Value Type: map[string]interface {}
Output: {key: value!@#$%^&*()_+-=<>?/}

Value Type: map[string]interface {}
Output: {largeArray: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10], nested: {field: value, anotherField: [100, 200, 300]}}

Value Type: string
Output: true

Value Type: map[string]interface {}
Output: {cpu: 512, essential: true}

```