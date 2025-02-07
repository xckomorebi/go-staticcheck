package cs1001

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
		Name:     "ST1025",
		Run:      run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	},
	Doc: &lint.RawDocumentation{
		Title:    "check cyclomatic complexity",
		Text:     `check cyclomatic complexity",`,
		Since:    "Unreleased",
		Severity: lint.SeverityInfo,
	},
})

var Analyzer = SCAnalyzer.Analyzer

func run(pass *analysis.Pass) (any, error) {
	var (
		complexity int
	)

	fn := func(node ast.Node, push bool) bool {
		if !push {
			switch node.(type) {
			case *ast.FuncDecl:
				if complexity > 4 {
					pass.Reportf(node.Pos(), "function has cyclomatic complexity of %d", complexity)
				}
				complexity = 0
			}
			return false
		}

		switch n := node.(type) {
		case *ast.FuncDecl:
			complexity = 1
		case *ast.IfStmt, *ast.RangeStmt, *ast.ForStmt:
			complexity++
		case *ast.CaseClause:
			if n.List != nil {
				complexity++
			}
		case *ast.CommClause:
			if n.Comm != nil {
				complexity++
			}
		case *ast.BinaryExpr:
			if n.Op == token.LAND || n.Op == token.LOR {
				complexity++
			}
		}
		return true
	}

	pass.ResultOf[inspect.Analyzer].(*inspector.Inspector).Nodes(
		[]ast.Node{
			(*ast.FuncDecl)(nil),
			(*ast.IfStmt)(nil),
			(*ast.ForStmt)(nil),
			(*ast.RangeStmt)(nil),
			(*ast.CommClause)(nil),
			(*ast.BinaryExpr)(nil),
			(*ast.CaseClause)(nil),
		},
		fn,
	)
	return nil, nil
}
