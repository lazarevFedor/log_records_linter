package pkg

import (
	"go/ast"
	"go/token"
	"go/types"
	"strconv"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var LogsAnalyzer = &analysis.Analyzer{
	Name:     "logs",
	Doc:      "Checks log records correct formatting",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}
	insp.Preorder(nodeFilter, func(n ast.Node) {
		callExpr := n.(*ast.CallExpr)
		if !isLogCall(pass, callExpr) {
			return
		}

		msg, msgPos := extractMessage(callExpr)
		if msg == "" {
			return
		}

		checkLowercaseStart(pass, msgPos, msg)
		checkEnglishOnly(pass, msgPos, msg)
		checkNoSpecialChars(pass, msgPos, msg)
		checkNoSensitiveData(pass, msgPos, msg)
	})

	return nil, nil
}

func extractMessage(call *ast.CallExpr) (string, token.Pos) {
	if len(call.Args) == 0 {
		return "", token.NoPos
	}

	lit, ok := call.Args[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", token.NoPos
	}

	msg, err := strconv.Unquote(lit.Value)
	if err != nil {
		return "", token.NoPos
	}

	return msg, lit.Pos()
}

func isLogCall(pass *analysis.Pass, callExpr *ast.CallExpr) bool {
	selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	if !isLogLevel(selectorExpr) {
		return false
	}

	if isStdLoggerSelector(pass, selectorExpr) {
		return true
	}

	if isLoggerType(pass, selectorExpr.X) {
		return true
	}

	return false
}

func isLogLevel(selectorExpr *ast.SelectorExpr) bool {
	switch selectorExpr.Sel.Name {
	case "Debug", "Info", "Warn", "Error", "Fatal", "Panic":
		return true
	}
	return false
}

func isStdLoggerSelector(pass *analysis.Pass, selectorExpr *ast.SelectorExpr) bool {
	ident, ok := selectorExpr.X.(*ast.Ident)
	if !ok {
		return false
	}

	pkgName, ok := pass.TypesInfo.Uses[ident].(*types.PkgName)
	if !ok || pkgName.Imported() == nil {
		return false
	}

	switch pkgName.Imported().Path() {
	case "log", "log/slog":
		return true
	default:
		return false
	}
}

func isLoggerType(pass *analysis.Pass, expr ast.Expr) bool {
	loggerType := pass.TypesInfo.TypeOf(expr)
	if loggerType == nil {
		return false
	}

	if ptr, ok := loggerType.(*types.Pointer); ok {
		loggerType = ptr.Elem()
	}

	named, ok := loggerType.(*types.Named)
	if !ok || named.Obj() == nil || named.Obj().Pkg() == nil {
		return false
	}

	pkgPath := named.Obj().Pkg().Path()
	name := named.Obj().Name()

	switch pkgPath {
	case "go.uber.org/zap":
		return name == "Logger" || name == "SugaredLogger"
	case "log/slog":
		return name == "Logger"
	default:
		return false
	}
}
