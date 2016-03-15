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

// package secne renders an object.
package scene

import (
	"log"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/go-gl-utils/app"
	"github.com/hurricanerix/go-gl-utils/shader"
	"github.com/hurricanerix/shader-tool/model"
)

const ( // Program IDs
	progID      = iota
	numPrograms = iota
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

type Scene struct {
	// Config vars
	ModelFile  string
	ColorFile  string
	NormalFile string

	// Internal vars (maybe should be private)
	Programs    [numPrograms]uint32
	VAOs        [numVAOs]uint32
	NumVertices [numVAOs]int32
	Buffers     [numBuffers]uint32

	MouseX          float32
	MouseY          float32
	MouseLeft       bool
	AmbientColor    mgl32.Vec4
	AmbientColorLoc int32
	LightPos        mgl32.Vec3
	LightPosLoc     int32
	LightColor      mgl32.Vec4
	LightColorLoc   int32
	LightPower      float32
	LightPowerLoc   int32
	Model           model.Model
	Angle           float32
	ModelMatrix     mgl32.Mat4
	ModelMatrixLoc  int32
}

// Setup resources required to update/display the scene.
func (s *Scene) Setup(ctx *app.Context) error {
	s.AmbientColor = mgl32.Vec4{0.2, 0.2, 0.2, 1.0}
	s.LightPos = mgl32.Vec3{0.0, 0.0, 10.0}
	s.LightColor = mgl32.Vec4{0.7, 0.7, 0.7}
	s.LightPower = 1000

	shaders := []shader.Info{

		shader.Info{Type: gl.VERTEX_SHADER, Filename: "assets/shaders/normalmap.vert"},
		shader.Info{Type: gl.FRAGMENT_SHADER, Filename: "assets/shaders/normalmap.frag"},
	}

	program, err := shader.Load(&shaders)
	if err != nil {
		return err
	}
	s.Programs[progID] = program

	gl.UseProgram(s.Programs[progID])

	gl.Enable(gl.DEPTH_TEST)

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(ctx.ScreenWidth)/float32(ctx.ScreenHeight), 0.1, 10.0)
	projectionUniform := gl.GetUniformLocation(s.Programs[progID], gl.Str("ProjMatrix\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GetUniformLocation(s.Programs[progID], gl.Str("ViewMatrix\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	modelMatrix := mgl32.Ident4()
	s.ModelMatrixLoc = gl.GetUniformLocation(s.Programs[progID], gl.Str("ModelMatrix\x00"))
	gl.UniformMatrix4fv(s.ModelMatrixLoc, 1, false, &modelMatrix[0])

	textureUniform := gl.GetUniformLocation(s.Programs[progID], gl.Str("ColorMap\x00"))
	gl.Uniform1i(textureUniform, 0)

	gl.BindFragDataLocation(s.Programs[progID], 0, gl.Str("FragColor\x00"))

	// Load the texture
	/*
		colorMap, err := os.Open("assets/textures/marble.png")
		defer colorMap.Close()
		if err != nil {
			log.Fatalln("failed to open tex:", err)
		}

			if err := shader.LoadTex(colorMap, gl.TEXTURE0); err != nil {
				log.Fatalln(err)
			}


		normalMap, err := os.Open("assets/textures/marble.normal.png")
		defer normalMap.Close()
		if err != nil {
			log.Fatalln("failed to open tex:", err)
		}
		if err := shader.LoadTex(normalMap, gl.TEXTURE1); err != nil {
			log.Fatalln(err)
		}
	*/

	mdlReader, err := os.Open(s.ModelFile)
	defer mdlReader.Close()
	if err != nil {
		log.Fatalln("could not open model:", err)
	}
	s.Model = model.New()
	if err := s.Model.Load(mdlReader); err != nil {
		log.Fatalln("could not load model:", err)
	}

	// Configure the vertex data
	gl.GenVertexArrays(numVAOs, &s.VAOs[0])
	gl.BindVertexArray(s.VAOs[triangleName])

	gl.GenBuffers(numBuffers, &s.Buffers[0])
	gl.BindBuffer(gl.ARRAY_BUFFER, s.Buffers[aBufferName])
	gl.BufferData(gl.ARRAY_BUFFER, len(s.Model.VertexData)*8, gl.Ptr(s.Model.VertexData), gl.STATIC_DRAW)

	mcVertexLoc := uint32(gl.GetAttribLocation(s.Programs[progID], gl.Str("MCVertex\x00")))
	gl.EnableVertexAttribArray(mcVertexLoc)
	gl.VertexAttribPointer(mcVertexLoc, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(0))

	mcNormalLoc := uint32(gl.GetAttribLocation(s.Programs[progID], gl.Str("MCNormal\x00")))
	gl.EnableVertexAttribArray(mcNormalLoc)
	gl.VertexAttribPointer(mcNormalLoc, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(3*4))

	texCoordLoc := uint32(gl.GetAttribLocation(s.Programs[progID], gl.Str("TexCoord0\x00")))
	gl.EnableVertexAttribArray(texCoordLoc)
	gl.VertexAttribPointer(texCoordLoc, 2, gl.FLOAT, false, 8*4, gl.PtrOffset(6*4))

	var fbo uint32
	gl.GenBuffers(1, &fbo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, fbo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(s.Model.FaceData)*4, gl.Ptr(s.Model.FaceData), gl.STATIC_DRAW)

	s.LightPosLoc = gl.GetUniformLocation(s.Programs[progID], gl.Str("LightPos\x00"))
	gl.Uniform3f(s.LightPosLoc, s.LightPos[0], s.LightPos[1], s.LightPos[2])

	s.AmbientColorLoc = gl.GetUniformLocation(s.Programs[progID], gl.Str("AmbientColor\x00"))
	gl.Uniform4f(s.AmbientColorLoc, s.AmbientColor[0], s.AmbientColor[1], s.AmbientColor[2], s.AmbientColor[3])

	s.LightColorLoc = gl.GetUniformLocation(s.Programs[progID], gl.Str("LightColor\x00"))
	gl.Uniform4f(s.LightColorLoc, s.LightColor[0], s.LightColor[1], s.LightColor[2], s.LightColor[3])

	s.LightPowerLoc = gl.GetUniformLocation(s.Programs[progID], gl.Str("LightPower\x00"))
	gl.Uniform1f(s.LightPowerLoc, s.LightPower)

	return nil
}

// Update the state of your scene.
func (s *Scene) Update(dt float32) {
	s.Angle += dt * 141
	s.ModelMatrix = mgl32.HomogRotate3D(float32(s.Angle), mgl32.Vec3{0, 1, 0})
	gl.ClearColor(s.AmbientColor[0], s.AmbientColor[1], s.AmbientColor[2], 1.0)
}

// Display the scene.
func (s *Scene) Display() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(s.Programs[progID])
	gl.UniformMatrix4fv(s.ModelMatrixLoc, 1, false, &s.ModelMatrix[0])
	gl.BindVertexArray(s.VAOs[triangleName])

	/*
		if _, ok := shader.Tex[gl.TEXTURE0]; ok {
			gl.ActiveTexture(gl.TEXTURE0)
			gl.BindTexture(gl.TEXTURE_2D, shader.Tex[gl.TEXTURE0])
		}

		if _, ok := shader.Tex[gl.TEXTURE1]; ok {
			gl.ActiveTexture(gl.TEXTURE1)
			gl.BindTexture(gl.TEXTURE_2D, shader.Tex[gl.TEXTURE1])
		}
	*/

	gl.DrawElements(gl.TRIANGLES, int32(s.Model.FaceCount)*3, gl.UNSIGNED_INT, nil)
}

// Cleanup any resources allocated in Setup.
func (s *Scene) Cleanup() {
	var id uint32
	for i := 0; i < numPrograms; i++ {
		id = s.Programs[i]
		gl.UseProgram(id)
		gl.DeleteProgram(id)
	}
}

func cursorPositionCallback(w *glfw.Window, x, y float64) {
	/*
		MouseX = float32(x)
		_, h := w.GetSize()
		MouseY = float32(h) - float32(y)
		LightPos[0] = MouseX
		LightPos[1] = MouseY
	*/
}

func KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Release && key == glfw.KeyEscape {
		w.SetShouldClose(true)
	}
	/*
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
	*/
}

func mouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	/*
		if (action == glfw.Press || action == glfw.Repeat) && button == glfw.MouseButtonLeft {
			MouseLeft = true
		}
		if action == glfw.Release && button == glfw.MouseButtonLeft {
			MouseLeft = false
		}
	*/
}

/*
func (s *Shader) LoadTex(r io.Reader, id uint32) error {
	img, _, err := image.Decode(r)
	if err != nil {
		return err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var tex uint32
	gl.GenTextures(1, &tex)
	gl.ActiveTexture(id)
	gl.BindTexture(gl.TEXTURE_2D, tex)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))
	s.Tex[id] = tex
	return nil
}
*/
