#version 330

uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

in vec3 vert;

out vec4 fragColor;

void main() {
    fragColor = vec4(0, 0.3, 0.9, 1.0);

    vec4 pos = projection * view * model * vec4(vert, 1);

    gl_Position = pos;
}
