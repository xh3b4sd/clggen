package main

import (
	"os"
	"path/filepath"
	"strings"
)

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func (c *clggen) depthOf(inputPath string) int {
	return len(strings.Split(inputPath, string(filepath.Separator))) - 1
}

func (c *clggen) shouldBeLoaded(inputPath string, info os.FileInfo) bool {
	if inputPath == c.Flags.FileNamePrefix {
		return false
	}
	if info.IsDir() {
		return false
	}
	// Note that strings.Split always returns a slice containing at least one
	// item, that is the input string itself if the input could not be split.
	// Thus we need to subtract 1 for a proper depth check.
	if c.Flags.Depth != -1 && c.depthOf(inputPath)-c.InputDepth >= c.Flags.Depth {
		return false
	}

	return true
}
