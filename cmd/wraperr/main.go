package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/srvc/wraperr"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	flag.Parse()

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	err = wraperr.NewDetector().CheckPackages(flag.Args())
	if err == nil {
		return nil
	}

	errs, ok := wraperr.UnwrapUnwrappedErrorsError(err)
	if !ok {
		return errors.WithStack(err)
	}

	for _, e := range errs.Errors() {
		pos, err := filepath.Rel(wd, e.Position.String())
		if err != nil {
			pos = e.Position.String()
		}
		fmt.Printf("%s:\t%s.%s\n", pos, e.Pkgname, e.Funcname)
	}

	return nil
}
