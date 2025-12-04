#version 410 core

layout(location = 0) in vec2 inPos;
layout(location = 1) in vec4 inColor;

uniform vec2 uWinSize;

out vec4 vColor;

void main() {
    float x = (inPos.x / uWinSize.x) * 2.0 - 1.0;
    float y = 1.0 - (inPos.y / uWinSize.y) * 2.0;
    gl_Position = vec4(x, y, 0, 1);
    vColor = inColor;
}
