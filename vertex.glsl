#version 410 core

layout (location = 0) in vec3 pos;
layout (location = 1) in vec3 in_color;

uniform vec2 shift;

out vec3 color;

void main()
{
    gl_Position = vec4(pos.x+shift.x, pos.y+shift.y, pos.z, 1.0);
    color = in_color;
}
