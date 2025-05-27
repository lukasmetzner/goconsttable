# goconsttable

Creates a Markdown table with documentation for your Go constants.

## Usage

1. Create a Go file:

```go
// example.go
const (
	// This is a comment
	TestValue string = "test-value"

	// This is the start of a comment
	// This is the end of a comment
	AnotherValue string = "another-value"
)
```

2. Run `goconsttable`:

```bash
goconsttable -path ./example.go
```

3. Final result:

```md
| Constant | Description |
| --- | --- |
| TestValue | This is a comment |
| AnotherValue | This is the start of a comment This is the end of a comment |
```
