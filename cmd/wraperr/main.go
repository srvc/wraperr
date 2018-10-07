package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/srvc/wraperr"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	flag.Parse()

	err := wraperr.NewDetector().CheckPackages(flag.Args())
	if err == nil {
		return nil
	}

	errs, ok := wraperr.UnwrapUnwrappedErrorsError(err)
	if !ok {
		return errors.WithStack(err)
	}

	fprintUnwrappedErrors(os.Stdout, errs)

	return nil
}

func fprintUnwrappedErrors(w io.Writer, errs wraperr.UnwrappedErrors) {
	wd, _ := os.Getwd()

	for _, e := range errs.Errors() {
		var occPos, retPos string
		if wd != "" {
			occPos, _ = filepath.Rel(wd, e.OccurredAt.String())
			retPos, _ = filepath.Rel(wd, e.ReturnedAt.String())
		}
		if occPos == "" {
			occPos = e.OccurredAt.String()
		}
		if retPos == "" {
			retPos = e.ReturnedAt.String()
		}
		fmt.Fprintf(w, "%s:\t%s:\t%s\n", retPos, occPos, sprintInlineCode(e.Line))
	}
}

func sprintInlineCode(s string) string {
	cc := 1
	c := cc
	for _, r := range s {
		if r == '`' {
			cc++
			if cc > c {
				c = cc
			}
		} else {
			cc = 1
		}
	}
	q := strings.Repeat("`", c)
	return q + s + q
}
