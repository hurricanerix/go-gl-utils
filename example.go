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

// package main demonstrates how to use go-gl-utils
package main

import (
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/go-gl-utils/app"
	"github.com/hurricanerix/go-gl-utils/path"
	"github.com/hurricanerix/go-gl-utils/shader"
)

func init() {
	if err := path.SetWorkingDir("github.com/hurricanerix/go-gl-utils"); err != nil {
		panic(err)
	}
}

func main() {
	// Create a config.  See app.Config for details on supported values.
	c := app.Config{
		Name:                "Example App",
		DefaultScreenWidth:  320,
		DefaultScreenHeight: 200,
		EscapeToQuit:        true,
		SupportedGLVers: []mgl32.Vec2{
			mgl32.Vec2{4, 3}, // Try to load a OpenGL 4.3 context.
			mgl32.Vec2{4, 1}, // If that fails, try to load a 4.1 contex.
			// If all fail, a.Run() will return an error.
		},
	}

	// Create an instance of your scene.
	// See app.Scene for details on this interface.
	s := &scene{}

	// Create a new app, providing a config and scene.
	a := app.New(c, s)
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Everything below this line is for the Scene implementation.
// this would probably go in it's own package, but creating a
// package for an example seems like overkill.

const ( // Program IDs
	simpleProgID = iota
	numPrograms  = iota
)

const ( // VAO Names
	triangleName = iota // VAO
	numVAOs      = iota
)

const ( // Buffer Names
	aBufferName = iota // Array Buffer
	numBuffers  = iota
)

const ( // Attrib Locations
	mcVertexLoc = 0
)

type scene struct {
	Programs    [numPrograms]uint32
	VAOs        [numVAOs]uint32
	NumVertices [numVAOs]int32
	Buffers     [numBuffers]uint32
}

// Setup resources required to update/display the scene.
func (s *scene) Setup() error {
	shaders := []shader.Info{
		shader.Info{Type: gl.VERTEX_SHADER, Filename: "simple.vert"},
		shader.Info{Type: gl.FRAGMENT_SHADER, Filename: "simple.frag"},
	}

	program, err := shader.Load(&shaders)
	if err != nil {
		return err
	}
	s.Programs[simpleProgID] = program

	gl.UseProgram(s.Programs[simpleProgID])

	vertices := []float32{
		-0.90, -0.90,
		0.90, -0.90,
		0.0, 0.90,
	}
	s.NumVertices[triangleName] = int32(len(vertices))

	gl.GenVertexArrays(numVAOs, &s.VAOs[0])
	gl.BindVertexArray(s.VAOs[triangleName])

	sizeVertices := len(vertices) * int(unsafe.Sizeof(vertices[0]))
	gl.GenBuffers(numBuffers, &s.Buffers[0])
	gl.BindBuffer(gl.ARRAY_BUFFER, s.Buffers[aBufferName])
	gl.BufferData(gl.ARRAY_BUFFER, sizeVertices, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(mcVertexLoc, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(mcVertexLoc)

	return nil
}

// Update the state of your scene.
func (s *scene) Update(dt float32) {
	// This scene does not change, so there is nothing here.
}

// Display the scene.
func (s *scene) Display() {
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.BindVertexArray(s.VAOs[triangleName])
	gl.DrawArrays(gl.TRIANGLES, 0, s.NumVertices[triangleName])
}

// Cleanup any resources allocated in Setup.
func (s *scene) Cleanup() {
	var id uint32
	for i := 0; i < numPrograms; i++ {
		id = s.Programs[i]
		gl.UseProgram(id)
		gl.DeleteProgram(id)
	}
}
