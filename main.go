package main

import (
	"log"

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
			if newCLGGen.Flags.CLGExpression == "" {
				log.Fatalf("%#v\n", maskAnyf(invalidConfigError, "clg expression must not be empty"))
			}
			if newCLGGen.Flags.InputPath == "" {
				log.Fatalf("%#v\n", maskAnyf(invalidConfigError, "input path must not be empty"))
			}
		},
	}

	// flags
	newCLGGen.Cmd.PersistentFlags().StringVarP(&newCLGGen.Flags.CLGExpression, "clg-expression", "c", "func (c *clg) calculate", "regular expression identifying CLG packages")
	newCLGGen.Cmd.PersistentFlags().StringVarP(&newCLGGen.Flags.InputPath, "input-path", "i", ".", "input path to load CLGs from")

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
