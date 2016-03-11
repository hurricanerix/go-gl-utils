#version 410 core

layout (location = 0) in vec4 mcVertex;

void main() {
	gl_Position = mcVertex;
}
