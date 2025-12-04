package gl41

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

var (
	fsProg uint32
	fsVBO  uint32
	fsVAO  uint32
)

var fullscreenVert = `
#version 410 core

layout(location = 0) in vec2 inPos;
layout(location = 1) in vec2 inUV;

out vec2 fragUV;

void main() {
    fragUV = inUV;
    gl_Position = vec4(inPos, 0.0, 1.0);
}
`

var fullscreenFrag = `
#version 410 core

in vec2 fragUV;
out vec4 fragColor;

uniform sampler2D uTex;

void main() {
    fragColor = texture(uTex, fragUV);
}
`

func InitFullscreen() error {

	prog, err := CreateProgram(fullscreenVert, fullscreenFrag)
	if err != nil {
		return err
	}
	fsProg = prog

	quad := []float32{
		// pos      // uv
		-1, -1, 0, 0,
		1, -1, 1, 0,
		1, 1, 1, 1,
		-1, 1, 0, 1,
	}

	indices := []uint32{
		0, 1, 2,
		2, 3, 0,
	}

	// VAO
	gl.GenVertexArrays(1, &fsVAO)
	gl.BindVertexArray(fsVAO)

	// VBO
	gl.GenBuffers(1, &fsVBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, fsVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(quad)*4, gl.Ptr(quad), gl.STATIC_DRAW)

	// EBO
	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(0))

	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(8))

	return nil
}

func DrawFullscreen(tex uint32) {
	gl.UseProgram(fsProg)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, tex)
	gl.Uniform1i(gl.GetUniformLocation(fsProg, gl.Str("uTex\x00")), 0)

	gl.BindVertexArray(fsVAO)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
}
