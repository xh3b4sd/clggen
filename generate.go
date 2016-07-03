package main

import (
	"bufio"
	"bytes"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

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

type TmplCtx struct {
	CLGName     string
	PackageName string
}

func (c *clggen) ExecGenerateCmd(cmd *cobra.Command, args []string) {
	err := filepath.Walk(c.Flags.InputPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		var isCLGPackage bool
		{
			raw, err := ioutil.ReadFile(path)
			if err != nil {
				return maskAny(err)
			}
			scanner := bufio.NewScanner(bytes.NewBuffer(raw))
			for scanner.Scan() {
				if strings.Contains(scanner.Text(), c.Flags.CLGExpression) {
					isCLGPackage = true
					break
				}
			}
			err = scanner.Err()
			if err != nil {
				return maskAny(err)
			}
		}

		if isCLGPackage {
			dirName := filepath.Base(filepath.Dir(path))
			newTmplCtx := TmplCtx{
				CLGName:     dirName,
				PackageName: strings.Replace(dirName, "-", "", -1),
			}

			for fileName, sourceCode := range Templates {
				// tmpl
				filePath := filepath.Join(filepath.Dir(path), fileName)
				tmpl, err := template.New(filePath).Parse(sourceCode)
				if err != nil {
					return maskAny(err)
				}
				var b bytes.Buffer
				err = tmpl.Execute(&b, newTmplCtx)
				if err != nil {
					return maskAny(err)
				}
				// format
				raw, err := format.Source(b.Bytes())
				if err != nil {
					return maskAny(err)
				}
				// write
				err = ioutil.WriteFile(filePath, raw, os.FileMode(0644))
				if err != nil {
					return maskAny(err)
				}
			}
		}

		return nil
	})

	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}
}
