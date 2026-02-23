package pkg

import (
	"go/ast"
	"go/token"
	"go/types"
	"strconv"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// LogsAnalyzer is the main analyzer for checking log message formatting.
// It can be configured via command-line flags or by providing a JSON config file.
var LogsAnalyzer = &analysis.Analyzer{
	Name:     "logs",
	Doc:      "Checks log records correct formatting",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// configPath is the path to the JSON configuration file, set via command-line flag.
var configPath string

// init registers the plugin and defines the command-line flag for the config file.
func init() {
	register.Plugin("logs", New)
	LogsAnalyzer.Flags.StringVar(&configPath, "config", "", "path to config JSON file")
}

// MySettings defines the configuration settings for the plugin, which can be loaded from a JSON file.
type MySettings struct {
	Config string `json:"config"`
}

// PluginExample implements the register.LinterPlugin interface, allowing it to be registered as a plugin for golangci-lint.
type PluginExample struct {
	settings MySettings
}

// BuildAnalyzers builds the analyzers for the plugin based on the provided settings.
func (f *PluginExample) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	settings := f.settings
	return []*analysis.Analyzer{
		{
			Name: "logs",
			Doc:  "Checks log records correct formatting",
			Run: func(pass *analysis.Pass) (interface{}, error) {
				return runWithSettings(pass, settings)
			},
			Requires: []*analysis.Analyzer{inspect.Analyzer},
		},
	}, nil
}

// GetLoadMode specifies that the plugin should be loaded in "info" mode,
// which allows it to access type information and other details necessary for analysis.
func (f *PluginExample) GetLoadMode() string {
	return register.LoadModeTypesInfo
}

// New is the constructor function for the plugin, which decodes the provided settings and returns an instance of the plugin.
func New(settings any) (register.LinterPlugin, error) {
	s, err := register.DecodeSettings[MySettings](settings)
	if err != nil {
		return nil, err
	}

	return &PluginExample{settings: s}, nil
}

// run is the main function that performs the analysis.
// It resolves the configuration and checks log messages based on the specified rules.
func run(pass *analysis.Pass) (interface{}, error) {
	return runWithSettings(pass, MySettings{Config: configPath})
}

// runWithSettings performs the analysis using the provided settings,
// allowing for configuration via a JSON file or command-line flags.
func runWithSettings(pass *analysis.Pass, settings MySettings) (interface{}, error) {
	cfg, err := ResolveConfig(settings.Config)
	if err != nil {
		return nil, err
	}

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

		if cfg.EnableLowercaseStart {
			checkLowercaseStart(pass, msgPos, msg)
		}
		if cfg.EnableNoSpecialChars {
			checkNoSpecialChars(pass, msgPos, msg)
		}
		if cfg.EnableSensitivePatterns {
			checkNoSensitiveData(pass, msgPos, msg)
		}
		if cfg.EnableEnglishOnly {
			checkEnglishOnly(pass, msgPos, msg)
		}
	})

	return nil, nil
}

// extractMessage attempts to extract the log message from the first argument of the log call.
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

// isLogCall checks if the given call expression is a log call by examining the function being called and its type information.
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

// isLogLevel checks if the selector expression corresponds
// to a common log level method (e.g., Debug, Info, Warn, Error, Fatal, Panic).
func isLogLevel(selectorExpr *ast.SelectorExpr) bool {
	switch selectorExpr.Sel.Name {
	case "Debug", "Info", "Warn", "Error", "Fatal", "Panic":
		return true
	}
	return false
}

// isStdLoggerSelector checks if the selector expression is a standard library logger (e.g., log or log/slog).
func isStdLoggerSelector(pass *analysis.Pass, selectorExpr *ast.SelectorExpr) bool {
	if pass == nil || pass.TypesInfo == nil {
		return false
	}

	ident, ok := selectorExpr.X.(*ast.Ident)
	if !ok {
		return false
	}

	obj, exists := pass.TypesInfo.Uses[ident]
	if !exists || obj == nil {
		return false
	}

	pkgName, ok := obj.(*types.PkgName)
	if !ok || pkgName == nil {
		return false
	}

	imported := pkgName.Imported()
	if imported == nil {
		return false
	}

	switch imported.Path() {
	case "log", "log/slog":
		return true
	default:
		return false
	}
}

// isLoggerType checks if the given expression corresponds to
// a logger type from supported logging libraries (e.g., zap.Logger, zap.SugaredLogger, slog.Logger).
func isLoggerType(pass *analysis.Pass, expr ast.Expr) bool {
	if pass == nil || pass.TypesInfo == nil || expr == nil {
		return false
	}

	loggerType := pass.TypesInfo.TypeOf(expr)
	if loggerType == nil {
		return false
	}

	if ptr, ok := loggerType.(*types.Pointer); ok {
		loggerType = ptr.Elem()
	}

	named, ok := loggerType.(*types.Named)
	if !ok || named == nil || named.Obj() == nil || named.Obj().Pkg() == nil {
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
