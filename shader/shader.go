package shader

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Shader struct {
	id uint32
}

func New(kind uint32, r io.Reader) (*Shader, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading shader: %w", err)
	}

	src, free := gl.Strs(terminate(string(b)))
	defer free()

	id := gl.CreateShader(kind)
	gl.ShaderSource(id, 1, src, nil)

	gl.CompileShader(id)
	if err := compileShaderError(id); err != nil {
		return nil, fmt.Errorf("compiling shader: %w", err)
	}

	return &Shader{id: id}, nil
}

func NewFromFile(kind uint32, name string) (*Shader, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	return New(kind, f)
}

func (s *Shader) Delete() { gl.DeleteShader(s.id) }

func terminate(s string) string {
	if !strings.HasSuffix(s, "\x00") {
		return s + "\x00"
	}
	return s
}

func compileShaderError(id uint32) error {
	var status int32
	gl.GetShaderiv(id, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var size int32
		gl.GetShaderiv(id, gl.INFO_LOG_LENGTH, &size)
		log := make([]byte, size)
		gl.GetShaderInfoLog(id, size, nil, &log[0])
		return fmt.Errorf("%s", log)
	}
	return nil
}
