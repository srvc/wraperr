package main

import (
	"github.com/srvc/wraperr"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(wraperr.Analyzer)
}
