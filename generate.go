package main

import (
	"bytes"
	"compress/gzip"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/fatih/camelcase"
	"github.com/spf13/cobra"
)

func (c *clggen) InitGenerateCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate code for all CLGs found in the specified packages.",
		Long:  "Generate code for all CLGs found in the specified packages.",
		Run:   c.ExecGenerateCmd,
	}

	return newCmd
}

type tmplCtx struct {
	CLGName string
	Package string
}

// TODO
//
//     collect all packages implementing func (c *clg) calculate
//     create template context for each package: directory name (foo-bar), package name (foobar), clg name (foo-bar),
//     generate source code by writing template files into identified clg package

func (c *clggen) ExecGenerateCmd(cmd *cobra.Command, args []string) {
	newTmplCtx := tmplCtx{
		CLGName: map[string][]byte{},
		Package: c.Flags.Package,
	}

	err := filepath.Walk(c.Flags.InputPath, func(path string, info os.FileInfo, err error) error {
		if !c.shouldBeLoaded(path, info) {
			return nil
		}

		raw, err := ioutil.ReadFile(path)
		if err != nil {
			return maskAny(err)
		}
		var b bytes.Buffer
		w := gzip.NewWriter(&b)
		_, err = w.Write(raw)
		w.Close()
		if err != nil {
			return maskAny(err)
		}
		newTmplCtx.AssetMap[path] = b.Bytes()

		return nil
	})

	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}

	// tmpl
	tmpl, err := template.New(c.Flags.FileNamePrefix).Parse(Template)
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}
	var b bytes.Buffer
	err = tmpl.Execute(&b, newTmplCtx)
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}

	// format
	raw, err := format.Source(b.Bytes())
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}

	err = ioutil.WriteFile(c.Flags.FileNamePrefix, raw, os.FileMode(0644))
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}
}
