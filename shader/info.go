// Copyright 2016 Richard Hawkins
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shader

import (
	"bytes"
	"fmt"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// Info representing a shader.
type Info struct {
	// Type of shader: gl.VERTEX_SHADER, gl.FRAGMENT_SHADER, ...
	Type uint32
	// Filename of shader source file.
	Filename string
	// shader ID.
	shader uint32
}

// Compile the shader using the info provided.
func (i *Info) Compile(program uint32) error {
	i.shader = gl.CreateShader(i.Type)
	source, err := readShader(i.Filename)
	if err != nil {
		return err
	}

	csrc, free := gl.Strs(source)
	gl.ShaderSource(i.shader, 1, csrc, nil)
	free()
	gl.CompileShader(i.shader)

	var compiled int32
	if gl.GetShaderiv(i.shader, gl.COMPILE_STATUS, &compiled); compiled == gl.FALSE {
		return fmt.Errorf("failed to compile %s: %s", i.Filename, getErrorMsg(true, i.shader))
	}
	gl.AttachShader(program, i.shader)
	return nil
}

func readShader(filename string) (string, error) {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(f)
	source := fmt.Sprintf("%s\x00", buf.String())

	return source, nil
}

func (i *Info) Delete() {
	gl.DeleteShader(i.shader)
}
