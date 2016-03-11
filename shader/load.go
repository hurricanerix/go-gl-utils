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
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// Load the shaders, returning the ID of the resulting program.  Any problems
// compiling or linking will result in an error.
func Load(shaders *[]Info) (uint32, error) {
	return load(shaders, false)
}

// LoadSeparable is the same as Load with the exception that before the link stage
// GL_PROGRAM_SEPARABLE is set to GL_TRUE.
func LoadSeparable(shaders *[]Info) (uint32, error) {
	return load(shaders, true)
}

// load the shaders
func load(shaders *[]Info, separable bool) (uint32, error) {
	program := gl.CreateProgram()

	for _, s := range *shaders {
		if err := s.Compile(program); err != nil {
			cleanup(shaders)
			gl.DeleteProgram(program)
			return 0, err
		}
	}

	if separable {
		gl.ProgramParameteri(program, gl.PROGRAM_SEPARABLE, gl.TRUE)
	}

	gl.LinkProgram(program)
	cleanup(shaders)
	var linked int32
	if gl.GetProgramiv(program, gl.LINK_STATUS, &linked); linked == gl.FALSE {
		msg := getErrorMsg(false, program)
		gl.DeleteProgram(program)
		return 0, fmt.Errorf("failed to link program: %s", msg)
	}

	return program, nil
}

// cleanup all shaders by calling Delete on any non-zero shader in the slice.
func cleanup(shaders *[]Info) {
	for _, s := range *shaders {
		if s.shader != 0 {
			s.Delete()
		}
	}
}
