package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/xh3b4sd/clggen/spec"
)

var (
	// version is the project version. It is given via buildflags that inject the
	// commit hash.
	version string
)

// Config represents the configuration used to create a new command line
// object.
type Config struct {
	// Dependencies.
	Cmd *cobra.Command

	// Settings.
	Flags   Flags
	Version string
}

// DefaultConfig provides a default configuration to create a new command line
// object by best effort.
func DefaultConfig() Config {
	newConfig := Config{
		Version: version,
	}

	return newConfig
}

// NewCLGGen creates a new configured command line object.
func NewCLGGen(config Config) (spec.CLGGen, error) {
	// clggen
	newCLGGen := &clggen{
		Config: config,
	}

	// command
	newCLGGen.Cmd = &cobra.Command{
		Use:   "clggen",
		Short: "Asset management and code generation. For more information see https://github.com/xh3b4sd/clggen.",
		Long:  "Asset management and code generation. For more information see https://github.com/xh3b4sd/clggen.",

		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},

		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if newCLGGen.Flags.Depth < -1 {
				log.Fatalf("%#v\n", maskAnyf(invalidConfigError, "depth must be greater than -1"))
			}
			if newCLGGen.Flags.CLGExpression == "" {
				log.Fatalf("%#v\n", maskAnyf(invalidConfigError, "clg expression must not be empty"))
			}
			if newCLGGen.Flags.FileNamePrefix == "" {
				log.Fatalf("%#v\n", maskAnyf(invalidConfigError, "output file name must not be empty"))
			}
			if newCLGGen.Flags.InputPath == "" {
				log.Fatalf("%#v\n", maskAnyf(invalidConfigError, "input path must not be empty"))
			}
			if newCLGGen.Flags.Package == "" {
				log.Fatalf("%#v\n", maskAnyf(invalidConfigError, "package must not be empty"))
			}

			// Calculate input depth.
			newCLGGen.InputDepth = newCLGGen.depthOf(newCLGGen.Flags.InputPath)
			fileInfo, err := os.Stat(newCLGGen.Flags.InputPath)
			if err != nil {
				log.Fatalf("%#v\n", maskAny(err))
			}
			if fileInfo.IsDir() && newCLGGen.Flags.InputPath[len(newCLGGen.Flags.InputPath)-1] != filepath.Separator {
				// In case the given input path represents a directory, but the given
				// input path does not contain a slash, we need to fix the input depth
				// explicitely by increasing the input depth by 1.
				newCLGGen.InputDepth++
			}
		},
	}

	// flags
	newCLGGen.Cmd.PersistentFlags().IntVarP(&newCLGGen.Flags.Depth, "depth", "d", 0, "depth of traversed directories")
	newCLGGen.Cmd.PersistentFlags().StringVarP(&newCLGGen.Flags.CLGExpression, "clg-expression", "c", "func (c *clg) calculate", "regular expression identifying CLG packages")
	newCLGGen.Cmd.PersistentFlags().StringVarP(&newCLGGen.Flags.FileNamePrefix, "file-name-prefix", "f", "generated", "prefx of the generated output file")
	newCLGGen.Cmd.PersistentFlags().StringVarP(&newCLGGen.Flags.InputPath, "input-path", "i", ".", "input path to load CLGs from")
	newCLGGen.Cmd.PersistentFlags().StringVarP(&newCLGGen.Flags.Package, "package", "p", "clg", "package name of the generated source code file")

	return newCLGGen, nil
}

func (c *clggen) Boot() {
	// init
	c.Cmd.AddCommand(c.InitGenerateCmd())
	c.Cmd.AddCommand(c.InitVersionCmd())

	// execute
	c.Cmd.Execute()
}

type clggen struct {
	Config

	// InputDepth describes the initial depth of the given input path. See the
	// following examples.
	//
	//     bar.ext          0
	//     foo/bar          1
	//     foo/bar.ext      1
	//     foo/bar/baz.ext  2
	//
	InputDepth int
}

func main() {
	newCLGGen, err := NewCLGGen(DefaultConfig())
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}

	newCLGGen.Boot()
}
