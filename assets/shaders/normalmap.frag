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

uniform vec4 AmbientColor;
uniform sampler2D ColorMap;
uniform sampler2D NormalMap;

uniform vec3 LightPos;
uniform vec4 LightColor;
uniform float LightPower;

in vec3 Pos;
in vec3 LightDir;
in vec3 EyeDir;
in vec2 TexCoord;

out vec4 FragColor;

void main() {
  float alpha = 1.0;
  vec3 diffuse = vec3(0.0, 0.0, 0.2);
  //alpha = texture(ColorMap, TexCoord.st).a;
  //diffuse = texture(ColorMap, TexCoord.st).rgb;
  vec3 ambient = AmbientColor.rgb * diffuse;
  vec3 specular = diffuse/8;

  vec3 normal = texture(NormalMap, TexCoord.st).rgb * 2 - 1;
  float distance = length(LightPos - Pos);

  vec3 n = normalize(normal);
  vec3 l = normalize(LightDir);

  float cosTheta = clamp(dot(n, l), 0.0, 1.0);

  vec3 e = normalize(EyeDir);
  vec3 r = reflect(-l, n);

  float cosAlpha = clamp(dot(e, r), 0.0, 1.0);

  FragColor = vec4(
    ambient +
    diffuse * LightColor.rgb * LightPower * cosTheta /
      (distance * distance) +
    specular * LightColor.rgb * LightPower * pow(cosAlpha, 5) /
      (distance * distance), alpha);
}
