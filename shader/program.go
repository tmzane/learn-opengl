package shader

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Program struct {
	id uint32
}

func NewProgram(shaders ...*Shader) (*Program, error) {
	id := gl.CreateProgram()
	for _, s := range shaders {
		gl.AttachShader(id, s.id)
	}

	gl.LinkProgram(id)
	if err := linkProgramError(id); err != nil {
		return nil, fmt.Errorf("linking program: %w", err)
	}

	for _, s := range shaders {
		s.Delete()
	}

	return &Program{id: id}, nil
}

func (p *Program) Use()    { gl.UseProgram(p.id) }
func (p *Program) Delete() { gl.DeleteProgram(p.id) }

func (p *Program) SetUniform(name string, values ...float32) error {
	location := gl.GetUniformLocation(p.id, gl.Str(terminate(name)))
	if location == -1 {
		return fmt.Errorf("uniform `%s` not found", name)
	}

	switch len(values) {
	case 1:
		gl.ProgramUniform1f(p.id, location, values[0])
	case 2:
		gl.ProgramUniform2f(p.id, location, values[0], values[1])
	case 3:
		gl.ProgramUniform3f(p.id, location, values[0], values[1], values[2])
	case 4:
		gl.ProgramUniform4f(p.id, location, values[0], values[1], values[2], values[3])
	default:
		panic("values count must not exceed 4")
	}

	return nil
}

func linkProgramError(id uint32) error {
	var status int32
	gl.GetProgramiv(id, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var size int32
		gl.GetProgramiv(id, gl.INFO_LOG_LENGTH, &size)
		log := make([]byte, size)
		gl.GetProgramInfoLog(id, size, nil, &log[0])
		return fmt.Errorf("%s", log)
	}
	return nil
}
