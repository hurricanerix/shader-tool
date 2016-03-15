Shader Tool
===========

Simple tool to assist me in writing GLSL shaders.

CLI Args
--------

- **color:** Filename of texture to use for color map.
- **frag:** List of fragment shaders filenames to compile (separated by commas). (default "assets/shaders/normalmap.frag")
- **height:** Set screen height in pixels.
- **model:** Filename of 3D model to render. (default "assets/models/cube.ply")
- **normal:** Filename of texture to use for normal map.
- **screen:** Set screen to display on. If set to 0, will run in windowed mode, otherwise will run in fullscreen mode.
- **vert:** List of vertex shader filenames to compile (separated by commas). (default "assets/shaders/normalmap.vert")
- **width:** Set screen width in pixels.

Example
-------

```
$ go run main.go -normal assets/textures/marble.normal.png -model assets/models/cube.ply
```

![Alt text](https://github.com/hurricanerix/shader-tool/raw/master/screenshot.png "Screenshot")
