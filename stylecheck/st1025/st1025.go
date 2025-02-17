package st1025

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/types"

	"honnef.co/go/tools/analysis/code"
	"honnef.co/go/tools/analysis/lint"
	"honnef.co/go/tools/analysis/report"
	"honnef.co/go/tools/internal/passes/buildir"
	"honnef.co/go/tools/pattern"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var SCAnalyzer = lint.InitializeAnalyzer(&lint.Analyzer{
	Analyzer: &analysis.Analyzer{
		Name:     "ST1025",
		Run:      run,
		Requires: []*analysis.Analyzer{inspect.Analyzer, buildir.Analyzer},
	},
	Doc: &lint.RawDocumentation{
		Title:    "check test",
		Text:     `check test`,
		Since:    "Unreleased",
		Severity: lint.SeverityInfo,
	},
})

var Analyzer = SCAnalyzer.Analyzer

var (
	ConstNameQ = pattern.MustParse(`(GenDecl "CONST" [(ValueSpec [(Ident "mName")] nil [(BasicLit "STRING" mName)])])`)
	RecvTypeQ  = pattern.MustParse(`(Or (StarExpr (Ident typeName)) (Ident typeName))`)
	UncErrQ    = pattern.MustParse(
		`(CallExpr 
			(SelectorExpr (Ident pkgName@(Or "xcerr" "udcerr")) (Ident methodName@(Or "New" "Wrap")))
			args
		)`)
)

func run(pass *analysis.Pass) (any, error) {
	var (
		constsMap   = make(map[string]*types.Const)
		typeNameMap = make(map[string]struct{})
	)
	irpkg := pass.ResultOf[buildir.Analyzer].(*buildir.IR).Pkg
	for _, m := range irpkg.Members {
		switch member := m.Object().(type) {
		case *types.Const:
			if member.Val().Kind() == constant.String && member.Name()[0] == 'n' {
				constsMap[member.Name()] = member
			}
		case *types.TypeName:
			typeNameMap[member.Name()] = struct{}{}
		}
	}

	for typeName := range typeNameMap {
		if consts, ok := constsMap["n"+typeName]; ok {
			if fmt.Sprintf(`"%s"`, typeName) != consts.Val().ExactString() {
				report.Report(pass, consts, fmt.Sprintf("const name should be %s", typeName))
			}
		} else {
			delete(constsMap, "n"+typeName)
		}
	}

	var (
		funcName string
		typeName string
		hasMName bool
	)
	_, _ = typeName, hasMName

	fn := func(node ast.Node, push bool) bool {
		if push {
			switch n := node.(type) {
			case *ast.GenDecl:
				if funcName == "" {
					return false
				}
				if m, ok := code.Match(pass, ConstNameQ, node); ok {
					hasMName = true
					if m.State["mName"] != funcName {
						report.Report(pass, node, fmt.Sprintf("mName should be %s", funcName))
					}
				}
			case *ast.FuncDecl:
				funcName = n.Name.Name
				if n.Recv != nil {
					m, _ := code.Match(pass, RecvTypeQ, n.Recv.List[0].Type)
					typeName = m.State["typeName"].(string)
				}
				return true
			case *ast.CallExpr:
				m, ok := code.Match(pass, UncErrQ, node)
				if !ok {
					return false
				}
				pkgName := m.State["pkgName"].(string)

				var nObj, mName ast.Expr
				switch m.State["methodName"].(string) {
				case "New":
					nObj = m.State["args"].([]ast.Expr)[0]
					mName = m.State["args"].([]ast.Expr)[1]
				case "Wrap":
					nObj = m.State["args"].([]ast.Expr)[1]
					mName = m.State["args"].([]ast.Expr)[2]
				}

				if typeName == "" {
					if sel, ok := nObj.(*ast.SelectorExpr); !ok ||
						code.SelectorName(pass, sel) != fmt.Sprintf("(%s).%s", typeName, funcName) {
						report.Report(pass, nObj, fmt.Sprintf("method name should be %s.%s", typeName, funcName))
					}
				} else {
					// if ident, ok := nObj.(*ast.Ident); ok {
					// 	if constsMap["n" + ident.Name]
					// }

				}
				_, _ = pkgName, mName

			}
			return false
		}

		if _, ok := node.(*ast.FuncDecl); ok {
			funcName = ""
			typeName = ""
			hasMName = false
		}
		return false
	}

	pass.ResultOf[inspect.Analyzer].(*inspector.Inspector).Nodes([]ast.Node{(*ast.GenDecl)(nil), (*ast.CallExpr)(nil), (*ast.FuncDecl)(nil)}, fn)

	return nil, nil
}
