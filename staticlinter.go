package main

import (
	"encoding/json"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
	"log"
	"os"
	"staticlinter/linters"
)

const Config = "config.json"

type ConfigData struct {
	Staticcheck []string
}

func main() {
	buf, err := os.ReadFile(Config)
	if err != nil {
		log.Fatalln(err)
	}

	var cfg *ConfigData
	err = json.Unmarshal(buf, &cfg)
	if err != nil {
		log.Fatalln(err)
	}

	// standard checks from "golang.org/x/tools/go/analysis/passes/..."
	checks := []*analysis.Analyzer{
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,

		// analyzer from example
		linters.ErrCheckAnalyzer,
		linters.OSExitCheckAnalyzer,
	}

	// checks from "honnef.co/go/tools/staticcheck"
	statics := make(map[string]bool)
	for _, v := range cfg.Staticcheck {
		statics[v] = true
	}
	//stChecks := make([]*analysis.Analyzer, len(statics))
	for _, v := range staticcheck.Analyzers {
		if statics[v.Analyzer.Name] {
			checks = append(checks, v.Analyzer)
		}
	}
	//log.Println(checks)

	multichecker.Main(
		checks...,
	)

}
