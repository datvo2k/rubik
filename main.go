package main

import (
	_ "image/png"
	"log"
	"rubik/gfx"
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	windowWidth  = 800
	windowHeight = 600
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Rubik 3D", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	window.SetKeyCallback(keyCallback)

	err = programLoop(window)
	if err != nil {
		log.Fatal(err)
	}
}

/*
* Create the Vertex Array Object for a cube
 */
func createVAO(vertices []float32, indices []uint32) uint32 {
	var VAO uint32
	// glGenVertexArrays — generate vertex array object names
	gl.GenVertexArrays(1, &VAO)

	var VBO uint32
	// glGenBuffers — generate buffer object names
	gl.GenBuffers(1, &VBO)

	var EBO uint32
	gl.GenBuffers(1, &EBO)

	// Bind the Vertex Array Object first, then bind and set vertex buffer(s) and attribute pointers()
	gl.BindVertexArray(VAO)

	// copy vertices data into VBO (it needs to be bound first)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	// glBufferData, glNamedBufferData — creates and initializes a buffer object's data store
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// size of one whole vertex (sum of attrib sizes)
	var stride int32 = 3*4 + 2*4
	var offset int = 0

	// position
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, stride, gl.PtrOffset(offset))
	gl.EnableVertexAttribArray(0)
	offset += 3 * 4

	// texture position
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, stride, gl.PtrOffset(offset))
	gl.EnableVertexAttribArray(1)
	offset += 2 * 4

	// unbind the VAO (safe practice so we don't accidentally (mis)configure it later)
	gl.BindVertexArray(0)

	return VAO
}

func keyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {

	// When a user presses the escape key, we set the WindowShouldClose property to true,
	// which closes the application
	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}
}

func programLoop(window *glfw.Window) error {

	// the linked shader program determines how the data will be rendered
	vertShader, err := gfx.NewShaderFromFile("shaders/basic.vert", gl.VERTEX_SHADER)
	if err != nil {
		return err
	}

	fragShader, err := gfx.NewShaderFromFile("shaders/colour.frag", gl.FRAGMENT_SHADER)
	if err != nil {
		return err
	}

	program, err := gfx.NewProgram(vertShader, fragShader)
	if err != nil {
		return err
	}
	defer program.Delete()

	vertices := []float32{
		// position        // texture position
		-0.5, -0.5, -0.5, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0,

		-0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,

		-0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, 0.5, 1.0, 0.0,

		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,

		-0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,

		-0.5, 0.5, -0.5, 0.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
	}

	indices := []uint32{}

	VAO := createVAO(vertices, indices)

	// matrix array
	cubePositions := [][]float32{
		{0.0, 0.0, 0.0},
	}

	gl.Enable(gl.DEPTH_TEST)

	for !window.ShouldClose() {
		// poll events and call their registered callbacks
		glfw.PollEvents()

		// background color
		/*
			glClearColor — specify clear values for the color buffers
			void glClearColor( 	GLfloat red,
								GLfloat green,
								GLfloat blue,
								GLfloat alpha);
		*/
		gl.ClearColor(0.1, 0.1, 0.1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// draw vertices
		program.Use()

		// Create transformation matrices
		rotateX := (mgl32.Rotate3DX(mgl32.DegToRad(-60 * float32(glfw.GetTime()))))
		rotateY := (mgl32.Rotate3DY(mgl32.DegToRad(-60 * float32(glfw.GetTime()))))
		rotateZ := (mgl32.Rotate3DZ(mgl32.DegToRad(-60 * float32(glfw.GetTime()))))

		viewTransform := mgl32.Translate3D(0, 0, -3)
		projectTransform := mgl32.Perspective(mgl32.DegToRad(60), windowWidth/windowHeight, 0.1, 100.0)

		gl.UniformMatrix4fv(program.GetUniformLocation("view"), 1, false,
			&viewTransform[0])
		gl.UniformMatrix4fv(program.GetUniformLocation("project"), 1, false,
			&projectTransform[0])

		gl.UniformMatrix4fv(program.GetUniformLocation("worldRotateX"), 1, false,
			&rotateX[0])
		gl.UniformMatrix4fv(program.GetUniformLocation("worldRotateY"), 1, false,
			&rotateY[0])
		gl.UniformMatrix4fv(program.GetUniformLocation("worldRotateZ"), 1, false,
			&rotateZ[0])

		gl.BindVertexArray(VAO)

		for _, pos := range cubePositions {

			worldTranslate := mgl32.Translate3D(pos[0], pos[1], pos[2])
			worldTransform := (worldTranslate.Mul4(rotateX.Mul3(rotateY).Mul3(rotateZ).Mat4()))

			gl.UniformMatrix4fv(program.GetUniformLocation("world"), 1, false,
				&worldTransform[0])

			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}
		// gl.DrawElements(gl.TRIANGLES, 36, gl.UNSIGNED_INT, unsafe.Pointer(nil))
		gl.BindVertexArray(0)

		// end of draw loop

		// swap in the rendered buffer
		window.SwapBuffers()
	}

	return nil
}

func colorCube() {
	cubeIndices := [24]byte{
		0, 3, 2, 1,
		2, 3, 7, 6,
		0, 4, 7, 3,
		1, 2, 6, 5,
		4, 5, 6, 7,
		0, 1, 5, 4,
	}
	p := unsafe.Pointer(&cubeIndices)
	gl.DrawElements(gl.QUADS, 24, gl.UNSIGNED_BYTE, p)
}
