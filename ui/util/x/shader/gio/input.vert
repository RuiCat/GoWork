#version 310 es

// SPDX-License-Identifier: Unlicense OR MIT

#extension GL_GOOGLE_include_directive : enable

precision highp float;

#include "common.h"

layout(std140,binding=0) uniform Block
{
    layout(offset=0) vec4 [4]Matrix;
} _block;

layout(location=0) in vec3 inPos;
layout(location=1) in vec3 inColor;
layout(location=2) in vec2 inUV;


layout(location=0) out vec2 vUV;
layout(location=1) out vec4 vertexColor;

void main() {
    vUV = inUV;
    vertexColor =  vec4(inColor,1);
    gl_Position =  mat4(_block.Matrix[0],_block.Matrix[1],_block.Matrix[2],_block.Matrix[3]) * vec4(inPos,1);
}
