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
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// getErrorsMsg helps to return an error message when compiling/linking goes
// wrong.  If shader is true, check for logs relating to a shader failing to
// compile.  In this context id should be the ID of the shader that failed to
// compile.  If shader is false, check for logs relating to linking shaters
// to a program.  in this context, id is the ID of the program that was
// attempting to link the shaders.
func getErrorMsg(shader bool, id uint32) string {
	var l int32
	if shader {
		gl.GetShaderiv(id, gl.INFO_LOG_LENGTH, &l)
	} else {
		gl.GetProgramiv(id, gl.INFO_LOG_LENGTH, &l)
	}

	msg := strings.Repeat("\x00", int(l+1))

	if shader {
		gl.GetShaderInfoLog(id, l, nil, gl.Str(msg))
	} else {
		gl.GetProgramInfoLog(id, l, nil, gl.Str(msg))
	}

	return msg
}
