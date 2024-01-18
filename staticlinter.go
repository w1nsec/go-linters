package main

import (
	"encoding/json"
	"github.com/alexkohler/nakedret"
	"github.com/gnieto/mulint/mulint"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
	"log"
	"os"
	"staticlinter/linters"
	"strings"
)

const Config = "config.json"

type ConfigData struct {
	Staticcheck []string
}

func ReadConfig() []*analysis.Analyzer {
	buf, err := os.ReadFile(Config)
	if err != nil {
		log.Fatalln(err)
	}

	var cfg *ConfigData
	err = json.Unmarshal(buf, &cfg)
	if err != nil {
		log.Fatalln(err)
	}
	// checks from "honnef.co/go/tools/staticcheck"
	statics := make(map[string]bool)
	for _, v := range cfg.Staticcheck {
		statics[v] = true
	}
	stChecks := make([]*analysis.Analyzer, len(statics))
	for _, v := range staticcheck.Analyzers {
		if statics[v.Analyzer.Name] {
			stChecks = append(stChecks, v.Analyzer)
		}
	}
	return stChecks
}

func AddStaticchecks(additionals []string) []*analysis.Analyzer {

	m := make(map[string]bool)
	for _, v := range additionals {
		m[v] = true
	}

	checks := make([]*analysis.Analyzer, 0)
	for _, v := range staticcheck.Analyzers {
		// ALL SA checks
		if strings.Contains(v.Analyzer.Name, "SA") {
			checks = append(checks)
		}
		// other checks
		if m[v.Analyzer.Name] {
			checks = append(checks)
		}
	}
	return checks
}

func main() {
	var (
		readConfig  = true
		additionals = []string{
			// SA	static checks
			// ALL

			// S	simple checks
			"S1011", // Use a single append to concatenate two slices
			"S1008", // Simplify returning boolean expression

			// ST	style checks
			"ST1001", // Dot imports are discouraged
			"ST1019", // Importing the same package multiple times

			// QF	quickfix
			"QF1001", // Apply De Morgan’s law
		}

		DefaultLines = uint(5)
	)

	// standard checks from "golang.org/x/tools/go/analysis/passes/..."
	checks := []*analysis.Analyzer{
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,

		// analyzer from example
		linters.ErrCheckAnalyzer,
		linters.OSExitCheckAnalyzer,

		// third-party analyzers
		nakedret.NakedReturnAnalyzer(DefaultLines),
		mulint.Mulint,
	}

	readConfig = false
	if readConfig {
		confChecks := ReadConfig()
		if confChecks != nil {
			checks = append(checks, confChecks...)
		}
	}

	checks = append(checks, AddStaticchecks(additionals)...)

	multichecker.Main(
		checks...,
	)

}
