package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

	"learn-opengl/color"
	"learn-opengl/object"
	"learn-opengl/shader"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		return fmt.Errorf("initializing glwf: %w", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(800, 600, "LearnOpenGL", nil, nil)
	if err != nil {
		return fmt.Errorf("creating window: %w", err)
	}

	window.MakeContextCurrent()
	window.SetFramebufferSizeCallback(func(_ *glfw.Window, width int, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	if err := gl.Init(); err != nil {
		return fmt.Errorf("initializing gl: %w", err)
	}

	vertex, err := shader.NewFromFile(gl.VERTEX_SHADER, "vertex.glsl")
	if err != nil {
		return fmt.Errorf("new vertex shader: %w", err)
	}

	fragment, err := shader.NewFromFile(gl.FRAGMENT_SHADER, "fragment.glsl")
	if err != nil {
		return fmt.Errorf("new fragment shader: %w", err)
	}

	program, err := shader.NewProgram(vertex, fragment)
	if err != nil {
		return fmt.Errorf("new shader program: %w", err)
	}
	defer program.Delete()

	program.Use()

	obj := object.New(
		[]object.Vertex{
			// positions(3) + colors(3)
			{0.5, -0.5, 0, 1, 0, 0},
			{-0.5, -0.5, 0, 0, 1, 0},
			{0, 0.5, 0, 0, 0, 1},
		},
		object.WithAttribute(3, false), // +color
	)

	handleInput := inputHandler(window)

	for !window.ShouldClose() {
		shift := handleInput()
		color.Clear(color.DarkGreen)

		if err := program.SetUniform("shift", shift[0], shift[1]); err != nil {
			return fmt.Errorf("setting uniform: %w", err)
		}

		obj.Draw()
		glfw.PollEvents()
		window.SwapBuffers()
	}

	return nil
}

func inputHandler(w *glfw.Window) func() [2]float32 {
	var shift [2]float32 // [x, y]

	return func() [2]float32 {
		const shiftStep = 0.005

		switch {
		case w.GetKey(glfw.KeyEscape) == glfw.Press:
			w.SetShouldClose(true)

		case w.GetKey(glfw.Key1) == glfw.Press:
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		case w.GetKey(glfw.Key2) == glfw.Press:
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

		case w.GetKey(glfw.KeyUp) == glfw.Press:
			shift[1] += shiftStep
		case w.GetKey(glfw.KeyDown) == glfw.Press:
			shift[1] -= shiftStep
		case w.GetKey(glfw.KeyRight) == glfw.Press:
			shift[0] += shiftStep
		case w.GetKey(glfw.KeyLeft) == glfw.Press:
			shift[0] -= shiftStep
		}

		return shift
	}
}
