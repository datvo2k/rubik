#version 410 core

in vec2 TexCoord;

out vec4 color;

void main()
{
    // mix the two textures together 
    color = vec4(1.0, 0.0, 0.0, 1.0);
}