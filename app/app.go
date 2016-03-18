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

// Package app manages the main game loop.
package app

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	screenWidth  int
	screenHeight int
	screen       int
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()

	flag.IntVar(&screenWidth, "width", 0, "Set screen width in pixels.")
	flag.IntVar(&screenHeight, "height", 0, "Set screen height in pixels.")
	flag.IntVar(&screen, "screen", 0, "Set screen to display on. If set to 0, will run in windowed mode, otherwise will run in fullscreen mode.")
}

var (
	sceneKeyCallback func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey)
)

// Scene to be rendered by the app.
type Scene interface {
	// Setup any OpenGL requirements for the scene.
	Setup(ctx *Context) error
	// Update the state of the scene.
	Update(dt float32)
	// Display the scene.
	Display()
	// Cleanup any resources the scene allocated in Setup.
	Cleanup()
}

// Config of the app.
type Config struct {
	// Name of app, to be placed in the title bar.
	Name string
	// DefaultScreenWidth in pixels.
	DefaultScreenWidth int
	// DefaultScreenHeight in pixels.
	DefaultScreenHeight int
	// EscapeToQuit gives an easy way to quit the app.
	EscapeToQuit bool
	// SupportedGLVers will be used when attempting to create a window
	// with the provided version of OpenGL starting at the beginning
	// of the slice.  Typically you would want your versions in decending
	// order in the slice.
	// For example if you want to support OpenGL 4.3/4.1/2.1 this would be set to
	// []mgl32.Vec2{
	//     mgl32.Vec2{4, 3},
	//     mgl32.Vec2{4, 1},
	//     mgl32.Vec2{2, 1},
	// }
	// NOTE: When providing more one version of OpenGL, it is the
	// responsibility of the Scene to ensure that valid calls for the version
	// of OpenGL are provided.
	SupportedGLVers    []mgl32.Vec2
	SkipFlushAfterLoop bool
	KeyCallback        func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey)
}

// Context of the app.
type Context struct {
	Config       Config
	ScreenWidth  int
	ScreenHeight int
	Scene        Scene
}

// New app context is returned.
func New(c Config, s Scene) Context {
	return Context{
		Config: c,
		Scene:  s,
	}
}

// Run the app.
func (ctx *Context) Run() error {
	flag.Parse()

	// Set Defaults if needed
	if ctx.Config.Name == "" {
		ctx.Config.Name = os.Args[0]
	}
	if screenWidth <= 0 {
		screenWidth = ctx.Config.DefaultScreenWidth
	}
	ctx.ScreenWidth = screenWidth
	if screenHeight <= 0 {
		screenHeight = ctx.Config.DefaultScreenHeight
	}
	ctx.ScreenHeight = screenHeight

	// NOTE: Using GLFW instead of GLUT
	if err := glfw.Init(); err != nil {
		return fmt.Errorf("failed to initialize glfw: %s", err)
	}
	defer glfw.Terminate()

	var window *glfw.Window
	var monitor *glfw.Monitor
	if screen != 0 {
		m := glfw.GetMonitors()
		if screen > len(m) {
			msg := "0 - windowed mode\n"
			for i := range m {
				msg += fmt.Sprintf("%d - %s\n", i+1, m[i].GetName())
			}
			return fmt.Errorf("invalid monitor %d, please select from the following:\n%s", screen, msg)
		}
		monitor = m[screen-1]
	}

	msg := ""
	for i := range ctx.Config.SupportedGLVers {
		var err error
		maj := int(ctx.Config.SupportedGLVers[i][0])
		min := int(ctx.Config.SupportedGLVers[i][1])

		glfw.WindowHint(glfw.Resizable, glfw.False)
		glfw.WindowHint(glfw.ContextVersionMajor, maj)
		glfw.WindowHint(glfw.ContextVersionMinor, min)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

		window, err = glfw.CreateWindow(screenWidth, screenHeight, ctx.Config.Name, monitor, nil)
		if err != nil {
			msg += fmt.Sprintf("trying to set GL version %d.%d: %s\n", maj, min, err)
			if i == len(ctx.Config.SupportedGLVers)-1 {
				return fmt.Errorf("failed to create window: %s", msg)
			}
			continue
		}
		break
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		return fmt.Errorf("unable to initialize Glow ... exiting: %s", err)
	}

	fmt.Println("OpenGL vendor", gl.GoStr(gl.GetString(gl.VENDOR)))
	fmt.Println("OpenGL renderer", gl.GoStr(gl.GetString(gl.RENDERER)))
	fmt.Println("OpenGL version", gl.GoStr(gl.GetString(gl.VERSION)))
	fmt.Println("GLSL version", gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION)))

	if err := ctx.Scene.Setup(ctx); err != nil {
		return err
	}

	sceneKeyCallback = ctx.Config.KeyCallback

	if ctx.Config.EscapeToQuit {
		window.SetKeyCallback(keyCallback)
	}

	// TODO: Change to reflect the time between ticks
	dt := float32(0.0001)
	for !window.ShouldClose() {
		ctx.Scene.Update(dt)
		ctx.Scene.Display()

		if !ctx.Config.SkipFlushAfterLoop {
			gl.Flush()
		}
		window.SwapBuffers()
		glfw.PollEvents()
	}

	ctx.Scene.Cleanup()
	return nil
}

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Release && key == glfw.KeyEscape {
		w.SetShouldClose(true)
	}
	if sceneKeyCallback != nil {
		sceneKeyCallback(w, key, scancode, action, mods)
	}
}
