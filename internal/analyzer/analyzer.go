package analyzer

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "noimportsasvar",
	Doc:  "Checks that a file's imports are not used as variable names.",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		ast.Inspect(f, inspectNode(pass))
	}

	return nil, nil
}

func inspectNode(pass *analysis.Pass) func(node ast.Node) bool {
	imports := make(map[string]string)

	return func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.ImportSpec:
			// don't add dot imports or udnerscore imports to list of imports to be checked
			if n.Name != nil && n.Name.Name != "_" && n.Name.Name != "." {
				imports[n.Name.Name] = n.Name.Name
				break
			}
			if n.Path != nil {
				if n.Name != nil && (n.Name.Name == "_" || n.Name.Name == ".") {
					break
				}
				splitPath := strings.SplitN(n.Path.Value, "/", -1)
				p := strings.ReplaceAll(splitPath[len(splitPath)-1], "\"", "")
				imports[p] = n.Path.Value
				break
			}

		case *ast.AssignStmt:
			for _, lhs := range n.Lhs {
				ident, ok := lhs.(*ast.Ident)
				if !ok || ident == nil {
					continue
				}

				validateIdent(ident, pass, imports)
			}

		case *ast.RangeStmt:
			outerIdent, ok := n.Key.(*ast.Ident)
			if !ok || outerIdent == nil || outerIdent.Obj == nil {
				break
			}

			stmt, ok := outerIdent.Obj.Decl.(*ast.AssignStmt)
			if !ok {
				break
			}

			for _, lhs := range stmt.Lhs {
				ident, ok := lhs.(*ast.Ident)
				if !ok {
					continue
				}

				validateIdent(ident, pass, imports)
			}

		case *ast.GenDecl:
			if n.Tok != token.VAR && n.Tok != token.CONST {
				break
			}

			for _, spec := range n.Specs {
				v, ok := spec.(*ast.ValueSpec)
				if !ok || v == nil {
					continue
				}

				for _, name := range v.Names {
					validateIdent(name, pass, imports)
				}
			}

		case *ast.FuncDecl:
			if n.Type == nil {
				break
			}

			if n.Type.Params != nil && len(n.Type.Params.List) > 0 {
				for _, field := range n.Type.Params.List {
					if field == nil {
						continue
					}
					for _, name := range field.Names {
						validateIdent(name, pass, imports)
					}
				}
			}

			if n.Type.Results != nil && len(n.Type.Results.List) > 0 {
				for _, field := range n.Type.Results.List {
					if field == nil {
						continue
					}
					for _, name := range field.Names {
						validateIdent(name, pass, imports)
					}
				}
			}

		default:
			return true
		}
		return true
	}
}

func validateIdent(ident *ast.Ident, pass *analysis.Pass, imports map[string]string) {
	if ident.Obj != nil && (ident.Obj.Kind == ast.Var || ident.Obj.Kind == ast.Con) {
		if impName, found := imports[ident.Name]; found {
			impName = strings.ReplaceAll(impName, "\"", "")
			pass.Reportf(ident.NamePos, "%s name '%s' shared with import '%s'", ident.Obj.Kind, ident.Name, impName)
		}
	}
}
