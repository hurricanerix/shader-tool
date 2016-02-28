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
uniform vec3 LightPos;

in vec3 MCVertex;
in vec3 MCNormal;
in vec2 TexCoord0;

out vec2 TexCoord;
out vec3 Pos;
out vec3 LightDir;
out vec3 EyeDir;

vec3 tangent(vec3 n) {
  vec3 t = vec3(0.0, 0.0, 0.0);
  vec3 c1 = cross(n, vec3(0.0, 0.0, 1.0)); 
  vec3 c2 = cross(n, vec3(0.0, 1.0, 0.0)); 
  if(length(c1) > length(c2)) {
    t = c1;	
  } else {
    t = c2;	
  }
  normalize(t);
  return t;
}


void main() {
  mat4 mvMatrix = ViewMatrix * ModelMatrix;
  vec4 ccVertex = mvMatrix * vec4(MCVertex, 1.0);
  gl_Position = ProjMatrix * ccVertex;
  Pos = vec4(ModelMatrix * vec4(MCVertex, 1.0)).xyz;

  TexCoord = vec3(TexCoord0, 1.0).st;

  mat3 normalMatrix = mat3x3(mvMatrix);
  normalMatrix = inverse(normalMatrix);
  normalMatrix = transpose(normalMatrix);

  // TODO: tangent should point in direction of increasing U
  // TODO: bi-tangent should point in direction of increasing V 
  vec3 MCTangent = tangent(MCNormal);
  mat3 mv3Matrix = mat3x3(mvMatrix);
  vec3 n = normalize(MCNormal); // TODO: fix normalize(mv3Matrix * MCNormal);
  vec3 t = normalize(mv3Matrix * MCTangent);
  vec3 b = normalize(mv3Matrix * cross(n, t));

  LightDir = vec3(ViewMatrix * vec4(LightPos, 0.0)) - vec3(ccVertex);
  vec3 v;
  v.x = dot(LightDir, t);
  v.y = dot(LightDir, b);
  v.z = dot(LightDir, n);
  LightDir = v;

  EyeDir = vec3(-ccVertex);
  v.x = dot(EyeDir, t);
  v.y = dot(EyeDir, b);
  v.z = dot(EyeDir, n);
  EyeDir = v;
}
