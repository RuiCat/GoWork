#version 310 es

// SPDX-License-Identifier: Unlicense OR MIT

precision mediump float;

layout(binding = 0) uniform sampler2D vTexture;

layout(location=0) in highp vec2 vUV;
layout(location=1) in highp vec4 vertexColor;

layout(location=0) out vec4 fragColor;

void main() {
	fragColor = texture(vTexture, vUV)*vertexColor;
}
