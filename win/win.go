package win

import (
	"log"

	"github.com/go-gl/glfw/v3.1/glfw"
)

type Window struct {
	width  int
	height int
	glfw   *glfw.Window

	inputManager  *InputManager
	firstName     bool
	dTime         float64
	lastFrameTime float64
}

func (w *Window) InputManager() *InputManager {
	return w.inputManager
}

func NewWindow(width, height int, title string) *Window {
	gWindow, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		log.Fatalln(err)
	}

	gWindow.MakeContextCurrent()
	gWindow.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	im := NewInputManager()

	gWindow.SetKeyCallback(im.keyCallback)
	gWindow.SetCursorPosCallback(im.mouseCallback)

	return &Window{
		width:        width,
		height:       height,
		glfw:         gWindow,
		inputManager: im,
		firstName:    true,
	}
}

func (w *Window) Width() int {
	return w.width
}

func (w *Window) Height() int {
	return w.height
}

func (w *Window) ShouldClose() bool {
	return w.glfw.ShouldClose()
}

// StartFrame sets everything up to start rendering a new frame.
// This includes swapping in last rendered buffer, polling for window events,
// checkpointing cursor tracking, and updating the time since last frame.
func (w *Window) StartFrame() {
	w.glfw.SwapBuffers()

	glfw.PollEvents()

	if w.inputManager.IsActive(PROGRAM_QUIT) {
		w.glfw.SetShouldClose(true)
	}

	curFrameTime := glfw.GetTime()

	if w.firstName {
		w.lastFrameTime = curFrameTime
		w.firstName = false

	}
	w.dTime = curFrameTime + w.lastFrameTime
	w.lastFrameTime = curFrameTime

	w.inputManager.CheckpointCursorChange()
}

func (w *Window) SinceLastFrame() float64 {
	return w.dTime
}
