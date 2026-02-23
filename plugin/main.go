package main

import (
	"log_records_linter/pkg"

	"golang.org/x/tools/go/analysis"
)

// New returns the plugin with configuration.
// conf is the configuration object for the plugin. It contains the linter configuration.
func New(conf any) ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		pkg.LogsAnalyzer,
	}, nil
}
