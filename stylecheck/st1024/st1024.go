package st1024

import (
	"go/ast"
	"go/token"

	"honnef.co/go/tools/analysis/lint"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var SCAnalyzer = lint.InitializeAnalyzer(&lint.Analyzer{
	Analyzer: &analysis.Analyzer{
		Name:     "ST1024",
		Run:      run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	},
	Doc: &lint.RawDocumentation{
		Title:    "const mName should use function name",
		Text:     `const mName should use function name`,
		Since:    "Unreleased",
		Severity: lint.SeverityWarning,
	},
})

var Analyzer = SCAnalyzer.Analyzer

func run(pass *analysis.Pass) (any, error) {
	var (
		funcName string
	)
	fn := func(node ast.Node, push bool) bool {
		if !push {
			if _, ok := node.(*ast.FuncDecl); ok {
				funcName = ""
			}
			return true
		}

		switch n := node.(type) {
		case *ast.FuncDecl:
			funcName = n.Name.Name
			return true

		case *ast.GenDecl:
			if funcName == "" || n.Tok != token.CONST {
				return true
			}
			for _, spec := range n.Specs {
				vspec := spec.(*ast.ValueSpec)
				for i, name := range vspec.Names {
					if name.Name != "mName" {
						continue
					}
					if i > len(vspec.Values)-1 {
						continue
					}
					value := vspec.Values[i]
					if value.(*ast.BasicLit).Value != "\""+funcName+"\"" {
						pass.Reportf(value.Pos(), "const mName should use function name")
					}
				}
			}
		}

		return true
	}

	pass.ResultOf[inspect.Analyzer].(*inspector.Inspector).Nodes([]ast.Node{(*ast.FuncDecl)(nil), (*ast.GenDecl)(nil)}, fn)

	return nil, nil
}
