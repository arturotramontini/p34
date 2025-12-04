
#version 410 core

layout (location = 0) in vec2 position;
out vec2 vUV;

void main() {
    vUV = (position + 1.0) * 0.5;   // da [-1,+1] a [0,1]
    gl_Position = vec4(position, 0.0, 1.0);
}
