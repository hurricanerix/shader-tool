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
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/shader-tool/model"
	"github.com/hurricanerix/shader-tool/shader"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
	AmbientColor = mgl32.Vec4{0.2, 0.2, 0.2, 1.0}
	LightPos = mgl32.Vec3{0.0, 0.0, 10.0}
	LightColor = mgl32.Vec4{0.7, 0.7, 0.7}
	LightPower = 1000
}

type Config struct {
	ModelName    string
	TextureGroup string
	ShaderGroup  string
	MajorVer     int
	MinorVer     int
	UseCore      bool
	WinWidth     int
	WinHeight    int
	Fullscreen   bool
}

var (
	MouseX       float32
	MouseY       float32
	MouseLeft    bool
	AmbientColor mgl32.Vec4
	LightPos     mgl32.Vec3
	LightColor   mgl32.Vec4
	LightPower   float32
)

func Run(config Config) {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, config.MajorVer)
	glfw.WindowHint(glfw.ContextVersionMinor, config.MinorVer)
	if config.UseCore {
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	}

	var monitor *glfw.Monitor
	if config.Fullscreen {
		// TODO: Maybe choose monitor based on config?
		// http://www.glfw.org/docs/latest/monitor.html#monitor_monitors
		monitor = glfw.GetPrimaryMonitor()
	}
	window, err := glfw.CreateWindow(config.WinWidth, config.WinHeight, "Shader-Tool", monitor, nil)
	if err != nil {
		log.Fatalln("failed to create window:", err)
	}
	window.MakeContextCurrent()
	window.SetKeyCallback(keyCallback)
	window.SetCursorPosCallback(cursorPositionCallback)
	window.SetMouseButtonCallback(mouseButtonCallback)

	// Initialize Glow
	if err := gl.Init(); err != nil {
		log.Fatalln("failed to init Glow:", err)
	}

	fmt.Println("OpenGL vendor", gl.GoStr(gl.GetString(gl.VENDOR)))
	fmt.Println("OpenGL renderer", gl.GoStr(gl.GetString(gl.RENDERER)))
	fmt.Println("OpenGL version", gl.GoStr(gl.GetString(gl.VERSION)))
	fmt.Println("GLSL version", gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION)))

	// Configure the vertex and fragment shaders
	vert, err := os.Open(fmt.Sprintf("%s.vert", config.ShaderGroup))
	defer vert.Close()
	if err != nil {
		log.Fatalln("failed to open vert shader src file:", err)
	}
	frag, err := os.Open(fmt.Sprintf("%s.frag", config.ShaderGroup))
	defer frag.Close()
	if err != nil {
		log.Fatalln("failed to open frag shader src file:", err)
	}
	shader := shader.New()
	if err := shader.Compile(vert, frag); err != nil {
		log.Fatalln("failed to compile shader:", err)
	}

	gl.UseProgram(shader.Prog)

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(config.WinWidth)/float32(config.WinHeight), 0.1, 10.0)
	projectionUniform := gl.GetUniformLocation(shader.Prog, gl.Str("ProjMatrix\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GetUniformLocation(shader.Prog, gl.Str("ViewMatrix\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	modelMatrix := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(shader.Prog, gl.Str("ModelMatrix\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &modelMatrix[0])

	textureUniform := gl.GetUniformLocation(shader.Prog, gl.Str("ColorMap\x00"))
	gl.Uniform1i(textureUniform, 0)

	gl.BindFragDataLocation(shader.Prog, 0, gl.Str("FragColor\x00"))

	// Load the texture
	colorMap, err := os.Open(fmt.Sprintf("%s.png", config.TextureGroup))
	defer colorMap.Close()
	if err != nil {
		log.Fatalln("failed to open tex:", err)
	}
	if err := shader.LoadTex(colorMap, gl.TEXTURE0); err != nil {
		log.Fatalln(err)
	}

	normalMap, err := os.Open(fmt.Sprintf("%s.normal.png", config.TextureGroup))
	defer normalMap.Close()
	if err != nil {
		log.Fatalln("failed to open tex:", err)
	}
	if err := shader.LoadTex(normalMap, gl.TEXTURE1); err != nil {
		log.Fatalln(err)
	}

	mdlReader, err := os.Open(config.ModelName)
	defer mdlReader.Close()
	if err != nil {
		log.Fatalln("could not open model:", err)
	}
	mdl := model.New()
	if err := mdl.Load(mdlReader); err != nil {
		log.Fatalln("could not load model:", err)
	}

	// Configure the vertex data
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(mdl.VertexData)*8, gl.Ptr(mdl.VertexData), gl.STATIC_DRAW)

	mcVertexLoc := uint32(gl.GetAttribLocation(shader.Prog, gl.Str("MCVertex\x00")))
	gl.EnableVertexAttribArray(mcVertexLoc)
	gl.VertexAttribPointer(mcVertexLoc, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(0))

	mcNormalLoc := uint32(gl.GetAttribLocation(shader.Prog, gl.Str("MCNormal\x00")))
	gl.EnableVertexAttribArray(mcNormalLoc)
	gl.VertexAttribPointer(mcNormalLoc, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(3*4))

	texCoordLoc := uint32(gl.GetAttribLocation(shader.Prog, gl.Str("TexCoord0\x00")))
	gl.EnableVertexAttribArray(texCoordLoc)
	gl.VertexAttribPointer(texCoordLoc, 2, gl.FLOAT, false, 8*4, gl.PtrOffset(6*4))

	var fbo uint32
	gl.GenBuffers(1, &fbo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, fbo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(mdl.FaceData)*4, gl.Ptr(mdl.FaceData), gl.STATIC_DRAW)

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	angle := 0.0
	previousTime := glfw.GetTime()

	for !window.ShouldClose() {
		gl.ClearColor(AmbientColor[0], AmbientColor[1], AmbientColor[2], 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		lightPosLoc := gl.GetUniformLocation(shader.Prog, gl.Str("LightPos\x00"))
		gl.Uniform3f(lightPosLoc, LightPos[0], LightPos[1], LightPos[2])

		ambientColorLoc := gl.GetUniformLocation(shader.Prog, gl.Str("AmbientColor\x00"))
		gl.Uniform4f(ambientColorLoc, AmbientColor[0], AmbientColor[1], AmbientColor[2], AmbientColor[3])

		lightColorLoc := gl.GetUniformLocation(shader.Prog, gl.Str("LightColor\x00"))
		gl.Uniform4f(lightColorLoc, LightColor[0], LightColor[1], LightColor[2], LightColor[3])

		lightPowerLoc := gl.GetUniformLocation(shader.Prog, gl.Str("LightPower\x00"))
		gl.Uniform1f(lightPowerLoc, LightPower)

		//for i := glfw.Joystick(0); i < 16; i++ {
		//i := glfw.Joystick(0)
		//fmt.Println("joystick:", i, glfw.JoystickPresent(i))
		//fmt.Println("Name:", glfw.GetJoystickName(i))
		//fmt.Println("Axis:", glfw.GetJoystickAxes(i))
		//fmt.Println("Buttons:", glfw.GetJoystickButtons(i))
		//}

		// Update
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time

		angle += elapsed
		modelMatrix = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})

		// Render
		gl.UseProgram(shader.Prog)
		gl.UniformMatrix4fv(modelUniform, 1, false, &modelMatrix[0])

		gl.BindVertexArray(vao)

		if _, ok := shader.Tex[gl.TEXTURE0]; ok {
			gl.ActiveTexture(gl.TEXTURE0)
			gl.BindTexture(gl.TEXTURE_2D, shader.Tex[gl.TEXTURE0])
		}

		if _, ok := shader.Tex[gl.TEXTURE1]; ok {
			gl.ActiveTexture(gl.TEXTURE1)
			gl.BindTexture(gl.TEXTURE_2D, shader.Tex[gl.TEXTURE1])
		}

		gl.DrawElements(gl.TRIANGLES, int32(mdl.FaceCount)*3, gl.UNSIGNED_INT, nil)

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func cursorPositionCallback(w *glfw.Window, x, y float64) {
	MouseX = float32(x)
	_, h := w.GetSize()
	MouseY = float32(h) - float32(y)
	LightPos[0] = MouseX
	LightPos[1] = MouseY
}

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Release && key == glfw.KeyEscape {
		w.SetShouldClose(true)
	}

	if action == glfw.Release && key == glfw.KeyEqual {
		LightPos[2] += 1
	}
	if action == glfw.Release && key == glfw.KeyMinus {
		LightPos[2] -= 1
	}

	if action == glfw.Release && key == glfw.KeyQ {
		AmbientColor[0] += 0.1
		if AmbientColor[0] > 1.0 {
			AmbientColor[0] = 1.0
		}
	}
	if action == glfw.Release && key == glfw.KeyA {
		AmbientColor[0] -= 0.1
		if AmbientColor[0] < 0 {
			AmbientColor[0] = 0
		}
	}
	if action == glfw.Release && key == glfw.KeyW {
		AmbientColor[1] += 0.1
		if AmbientColor[1] > 1.0 {
			AmbientColor[1] = 1.0
		}
	}
	if action == glfw.Release && key == glfw.KeyS {
		AmbientColor[1] -= 0.1
		if AmbientColor[1] < 0 {
			AmbientColor[1] = 0
		}
	}
	if action == glfw.Release && key == glfw.KeyE {
		AmbientColor[2] += 0.1
		if AmbientColor[2] > 1.0 {
			AmbientColor[2] = 1.0
		}
	}
	if action == glfw.Release && key == glfw.KeyD {
		AmbientColor[2] -= 0.1
		if AmbientColor[2] < 0 {
			AmbientColor[2] = 0
		}
	}

	if action == glfw.Release && key == glfw.KeyR {
		LightColor[0] += 0.1
		if LightColor[0] > 1.0 {
			LightColor[0] = 1.0
		}
	}
	if action == glfw.Release && key == glfw.KeyF {
		LightColor[0] -= 0.1
		if LightColor[0] < 0 {
			LightColor[0] = 0
		}
	}
	if action == glfw.Release && key == glfw.KeyT {
		LightColor[1] += 0.1
		if LightColor[1] > 1.0 {
			LightColor[1] = 1.0
		}
	}
	if action == glfw.Release && key == glfw.KeyG {
		LightColor[1] -= 0.1
		if LightColor[1] < 0 {
			LightColor[1] = 0
		}
	}
	if action == glfw.Release && key == glfw.KeyY {
		LightColor[2] += 0.1
		if LightColor[2] > 1.0 {
			LightColor[2] = 1.0
		}
	}
	if action == glfw.Release && key == glfw.KeyH {
		LightColor[2] -= 0.1
		if LightColor[2] < 0 {
			LightColor[2] = 0
		}
	}

	if action == glfw.Release && key == glfw.KeyU {
		LightPower += 10
	}
	if action == glfw.Release && key == glfw.KeyJ {
		LightPower -= 10
	}
}

func mouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if (action == glfw.Press || action == glfw.Repeat) && button == glfw.MouseButtonLeft {
		MouseLeft = true
	}
	if action == glfw.Release && button == glfw.MouseButtonLeft {
		MouseLeft = false
	}
}
