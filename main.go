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

// Package main parses CLI arguments and runs the app.
package main

import (
	"flag"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/go-gl-utils/app"
	"github.com/hurricanerix/go-gl-utils/path"
	"github.com/hurricanerix/shader-tool/scene"
)

var (
	modelFile  string
	colorFile  string
	normalFile string
	vertFiles  string
	fragFiles  string
)

func init() {
	flag.StringVar(&modelFile, "model", "assets/models/cube.ply", "Filename of 3D model to render.")
	flag.StringVar(&colorFile, "color", "", "Filename of texture to use for color map.")
	flag.StringVar(&normalFile, "normal", "", "Filename of texture to use for normal map.")
	flag.StringVar(&vertFiles, "vert", "assets/shaders/normalmap.vert", "List of vertex shader filenames to compile (separated by commas).")
	flag.StringVar(&fragFiles, "frag", "assets/shaders/normalmap.frag", "List of fragment shaders filenames to compile (separated by commas).")

	if err := path.SetWorkingDir("github.com/hurricanerix/shader-tool"); err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()

	// Create a config.  See app.Config for details on supported values.
	c := app.Config{
		Name:                "Shader Tool",
		DefaultScreenWidth:  640,
		DefaultScreenHeight: 480,
		EscapeToQuit:        true,
		SupportedGLVers: []mgl32.Vec2{
			mgl32.Vec2{4, 3}, // Try to load a OpenGL 4.3 context.
			mgl32.Vec2{4, 1}, // If that fails, try to load a 4.1 contex.
			// If all fail, a.Run() will return an error.
		},
		KeyCallback: scene.KeyCallback,
	}

	// Create an instance of your scene.
	// See app.Scene for details on this interface.
	s := &scene.Scene{
		ModelFile:  modelFile,
		ColorFile:  colorFile,
		NormalFile: normalFile,
		VertFiles:  strings.Split(vertFiles, ","),
		FragFiles:  strings.Split(fragFiles, ","),
	}

	// Create a new app, providing a config and scene.
	a := app.New(c, s)
	if err := a.Run(); err != nil {
		panic(err)
	}
}
