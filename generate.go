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

// TmplCtx represents the context provided to the template parsing. Using the
// template context the CLG name and package name is available inside the
// templates.
type TmplCtx struct {
	CLGName          string
	IsErrorInterface bool
	PackageName      string
}

func (c *clggen) ExecGenerateCmd(cmd *cobra.Command, args []string) {
	newTemplates, err := loadTemplates(c.Flags.TemplateDir)
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}

	err = filepath.Walk(c.Flags.CLGDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		var isCLGPackage bool
		var isErrorInterface bool
		{
			raw, err := ioutil.ReadFile(path)
			if err != nil {
				return maskAny(err)
			}
			scanner := bufio.NewScanner(bytes.NewBuffer(raw))
			for scanner.Scan() {
				if strings.HasPrefix(scanner.Text(), c.Flags.CLGExp) {
					isCLGPackage = true

					if strings.HasSuffix(scanner.Text(), "error {") || strings.HasSuffix(scanner.Text(), ", error) {") {
						isErrorInterface = true
					}

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
				CLGName:          dirName,
				IsErrorInterface: isErrorInterface,
				PackageName:      strings.Replace(dirName, "-", "", -1),
			}

			for fileName, sourceCode := range newTemplates {
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
