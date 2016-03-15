Shader Tool
===========

Simple tool to assist me in writing GLSL shaders.

```
$ go run main.go -h
Usage of /var/folders/20/6fzsltdd42374k31kfnyg8540000gn/T/go-build758895528/command-line-arguments/_obj/exe/main:
  -color string
    	Filename of texture to use for color map.
  -frag string
    	List of fragment shaders filenames to compile (separated by commas). (default "assets/shaders/normalmap.frag")
  -height int
    	Set screen height in pixels.
  -model string
    	Filename of 3D model to render. (default "assets/models/cube.ply")
  -normal string
    	Filename of texture to use for normal map.
  -screen int
    	Set screen to display on. If set to 0, will run in windowed mode, otherwise will run in fullscreen mode.
  -vert string
    	List of vertex shader filenames to compile (separated by commas). (default "assets/shaders/normalmap.vert")
  -width int
    	Set screen width in pixels.
exit status 2
$ go run main.go -normal assets/textures/marble.normal.png -model assets/models/cube.ply
```

![Alt text](https://github.com/hurricanerix/shader-tool/raw/master/screenshot.png "Screenshot")
