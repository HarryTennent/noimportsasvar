Archived: This linter already exists https://revive.run/r#import-shadowing

# noimportsasvar

## Overview

This is a custom linter that detects variable names, constant names, and function parameters that share names with imported packages. The linter analyzes source code and reports occurrences where these names clash with names used in import statements.

### Installation

1. Run the following command to install the `noimportsasvar` linter package:

```bash
go install github.com/HarryTennent/noimportsasvar/cmd/noimportsasvar@latest
```

### Usage

1. Once the installation process is complete, you can use the `noimportsasvar` linter by running the following command:

```bash
noimportsasvar /path/to/your/package
```

2. Replace `/path/to/your/package` with the actual path to the directory containing your Go source files.

3. The `noimportsasvar` linter will analyze your Go code, detect any variable names, constant names, and function parameter names that clash with import package names, and report any occurrences found. It will respect dot imports (imports with `.`) and underscores (anonymous variables) and avoid reporting clashes for them.

## Example

Consider the following Go code:

```go
package main

import (
	"fmt"
	"go/ast"
)

func main() {
	ast := 10 
	fmt.Println(ast)
}

func doSomething(ast int) {
	fmt.Println(ast)
}
```

The linter will detect the variable `ast` in the `main` function and the function parameter `ast` in the `doSomething` function, both of which clash with the import of the package `go/ast`. The linter will generate reports for these occurrences, helping you identify potential naming conflicts.

## Contributing

Contributions to this linter are welcome! If you find any issues or have suggestions for improvements, feel free to open a new issue or submit a pull request.

## License

This linter is licensed under the [MIT License](LICENSE). Feel free to use, modify, and distribute it according to the terms of the license.
