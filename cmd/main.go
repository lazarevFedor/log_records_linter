package main

import (
	"log_records_linter/pkg"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(pkg.LogsAnalyzer)
}
