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

	"github.com/hurricanerix/shader-tool/app"
)

var (
	config app.Config
)

func init() {
	flag.StringVar(&config.ModelName, "model", "assets/models/cube.ply", "Name of 3D model to render.")
	flag.StringVar(&config.TextureGroup, "texture", "assets/textures/marble", "Name of texture group to use.")
	flag.StringVar(&config.ShaderGroup, "shader", "assets/shaders/normalmap", "Name of shader group to use.")
	flag.IntVar(&config.MajorVer, "major", 4, "OpenGL Major Version.")
	flag.IntVar(&config.MinorVer, "minor", 1, "OpenGL Minor Version.")
	flag.BoolVar(&config.UseCore, "use-core", true, "Use OpenGL Core Profile.")
	flag.IntVar(&config.WinWidth, "width", 800, "Window width in pixels.")
	flag.IntVar(&config.WinHeight, "height", 600, "Window height in pixels.")
	flag.BoolVar(&config.Fullscreen, "fullscreen", false, "Run in fullscreen window.")
}

func main() {
	flag.Parse()

	app.Run(config)
}
