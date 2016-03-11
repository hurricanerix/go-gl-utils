// Borrwoed from: https://github.com/go-gl/examples/blob/master/glfw31-gl41core-cube/cube.go#L321
// This isn't really OpenGL specific, so maybe should go elsewhere.

// package path
package path

import (
	"fmt"
	"go/build"
	"os"
)

// SetWorkingDir to path provided. (in the context of gopath)
func SetWorkingDir(path string) error {
	dir, err := importPathToDir(path)
	if err != nil {
		return fmt.Errorf("Unable to find Go package in your GOPATH, it's needed to load assets:", err)
	}
	err = os.Chdir(dir)
	if err != nil {
		return fmt.Errorf("os.Chdir:", err)
	}
	return nil
}

// importPathToDir resolves the absolute path from importPath.
// There doesn't need to be a valid Go package inside that import path,
// but the directory must exist.
func importPathToDir(importPath string) (string, error) {
	p, err := build.Import(importPath, "", build.FindOnly)
	if err != nil {
		return "", err
	}
	return p.Dir, nil
}
