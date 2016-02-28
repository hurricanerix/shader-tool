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
#version 330

uniform mat4 ProjMatrix;
uniform mat4 ViewMatrix;
uniform mat4 ModelMatrix;

in vec3 MCVertex;
in vec3 MCNormal;
in vec2 TexCoord0;

out vec2 TexCoord;

void main() {
    TexCoord = TexCoord0;
    gl_Position = ProjMatrix * ViewMatrix * ModelMatrix * vec4(MCVertex, 1);
}
