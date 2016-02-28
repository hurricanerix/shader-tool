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

// Package shader handles GLSL program compilation.
package shader

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"io"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Shader struct {
	Prog uint32
	Tex  map[uint32]uint32
}

func New() Shader {
	s := Shader{
		Tex: make(map[uint32]uint32),
	}
	return s
}

func (s *Shader) Compile(vert, frag io.Reader) error {
	vertexShader, err := compileShader(vert, gl.VERTEX_SHADER)
	if err != nil {
		return err
	}

	fragmentShader, err := compileShader(frag, gl.FRAGMENT_SHADER)
	if err != nil {
		return err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	s.Prog = program

	return nil
}

func compileShader(r io.Reader, t uint32) (uint32, error) {
	shader := gl.CreateShader(t)

	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	src := fmt.Sprintf("%s\x00", buf.String())

	csrc, free := gl.Strs(src)
	gl.ShaderSource(shader, 1, csrc, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	if gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status); status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("failed to compile %s: %v", r, log)
	}

	return shader, nil
}

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
