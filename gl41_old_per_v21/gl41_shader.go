package gl41

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func CompileShader(src string, shaderType uint32) (uint32, error) {
	sh := gl.CreateShader(shaderType)

	csources, free := gl.Strs(src + "\x00")
	gl.ShaderSource(sh, 1, csources, nil)
	free()

	gl.CompileShader(sh)

	var status int32
	gl.GetShaderiv(sh, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(sh, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(sh, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("shader compile error: %s", log)
	}

	return sh, nil
}

func CreateProgram(vsrc, fsrc string) (uint32, error) {
	vs, err := CompileShader(vsrc, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fs, err := CompileShader(fsrc, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vs)
	gl.AttachShader(prog, fs)
	gl.LinkProgram(prog)

	var status int32
	gl.GetProgramiv(prog, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(prog, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(prog, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("program link error: %s", log)
	}

	return prog, nil
}
