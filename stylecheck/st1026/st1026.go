package st1026

import (
	"go/ast"
	"go/token"

	"honnef.co/go/tools/analysis/lint"
	"honnef.co/go/tools/analysis/report"

	"github.com/xckomorebi/collections"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var SCAnalyzer = lint.InitializeAnalyzer(&lint.Analyzer{
	Analyzer: &analysis.Analyzer{
		Name:     "ST1026",
		Run:      run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	},
	Doc: &lint.RawDocumentation{
		Title:    "check assigned too early",
		Text:     `check assigned too early`,
		Since:    "Unreleased",
		Severity: lint.SeverityWarning,
	},
})

var Analyzer = SCAnalyzer.Analyzer

type blockInfo struct {
	identStates map[string]*identStates

	parent *blockInfo
}

func newBlockInfo(parent *blockInfo) *blockInfo {
	return &blockInfo{
		identStates: make(map[string]*identStates),
		parent:      parent,
	}
}

type usageFsm uint8

const (
	justAssigned usageFsm = iota
	inFirstSubBlock
	afterFirstSubBlock
)

func (u *usageFsm) next(push bool) {
	if *u == justAssigned && push {
		*u = inFirstSubBlock
		return
	}

	if *u == inFirstSubBlock && !push {
		*u = afterFirstSubBlock
	}
}

type identStates struct {
	node                  ast.Node
	fsm                   usageFsm
	usedImmidiately       bool
	usedInSubBlock        bool
	usedAfterNextSubBlock bool
}

func (b *blockInfo) Check(pass *analysis.Pass) {
	for _, state := range b.identStates {
		if state.usedImmidiately {
			continue
		}
		if state.usedInSubBlock && !state.usedAfterNextSubBlock {
			report.Report(pass, state.node, "assigned outside of usage scope")
		}

		if !state.usedInSubBlock && state.usedAfterNextSubBlock {
			report.Report(pass, state.node, "assigned too early")
		}
	}
}

func run(pass *analysis.Pass) (any, error) {
	var blkInfoStack collections.Stack[*blockInfo]

	fn := func(node ast.Node, push bool, stack []ast.Node) bool {
		if !push {
			switch node.(type) {
			case *ast.BlockStmt:
				pop, _ := blkInfoStack.Pop()
				pop.Check(pass)

				curInfo, _ := blkInfoStack.Peek()
				if curInfo == nil {
					return false
				}
				for _, state := range curInfo.identStates {
					state.fsm.next(push)
				}
				return false
			}
		}

		top, _ := blkInfoStack.Peek()

		switch stmt := node.(type) {
		case *ast.AssignStmt:
			for _, lhs := range stmt.Lhs {
				if ident, ok := lhs.(*ast.Ident); ok {
					top.identStates[ident.Name] = &identStates{node: ident}
				}
			}
		case *ast.BlockStmt:
			blkInfoStack.Push(newBlockInfo(top))
			if top == nil {
				return true
			}
			for _, state := range top.identStates {
				state.fsm.next(push)
			}
		case *ast.Ident:
			if len(stack) < 2 {
				return false
			}
			parent := stack[len(stack)-2]
			if parent, ok := parent.(*ast.AssignStmt); ok && parent.Tok == token.DEFINE {
				for _, lhs := range parent.Lhs {
					if ident, ok := lhs.(*ast.Ident); ok && ident.Name == stmt.Name {
						return false
					}
				}
			}

			curBlkInfo := top
			var state *identStates
			var ok bool
			for curBlkInfo != nil {
				if state, ok = curBlkInfo.identStates[stmt.Name]; ok {
					break
				}
				curBlkInfo = curBlkInfo.parent
			}
			if curBlkInfo == nil {
				return false
			}

			switch state.fsm {
			case justAssigned:
				state.usedImmidiately = true
			case inFirstSubBlock:
				state.usedInSubBlock = true
			case afterFirstSubBlock:
				state.usedAfterNextSubBlock = true
			}
		}

		return true
	}

	pass.ResultOf[inspect.Analyzer].(*inspector.Inspector).WithStack(
		nil,
		fn,
	)

	return nil, nil
}
