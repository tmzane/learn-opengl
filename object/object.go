package object

import (
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Vertex []float32

type Object struct {
	vao uint32
	vbo uint32
	ebo uint32

	index  uint32
	stride int32
	offset int

	verticesCount int32
	indicesCount  int32
}

func New(vertices []Vertex, opts ...Option) *Object {
	obj := Object{
		stride:        int32(len(vertices[0]) * sizeOf[float32]()),
		verticesCount: int32(len(vertices)),
	}

	gl.GenVertexArrays(1, &obj.vao)
	gl.BindVertexArray(obj.vao)
	defer gl.BindVertexArray(0)

	var data []float32
	for _, v := range vertices {
		data = append(data, v...)
	}

	gl.GenBuffers(1, &obj.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, obj.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*sizeOf[float32](), gl.Ptr(data), gl.STATIC_DRAW)

	WithAttribute(3, false)(&obj) // +coordinates
	for _, opt := range opts {
		opt(&obj)
	}

	return &obj
}

func (o *Object) Draw() {
	gl.BindVertexArray(o.vao)
	defer gl.BindVertexArray(0)

	if o.ebo != 0 {
		gl.DrawElements(gl.TRIANGLES, o.indicesCount, gl.UNSIGNED_INT, gl.PtrOffset(0))
	} else {
		gl.DrawArrays(gl.TRIANGLES, 0, o.verticesCount)
	}
}

func (o *Object) Delete() {
	gl.DeleteVertexArrays(1, &o.vao)
	gl.DeleteBuffers(1, &o.vbo)
	gl.DeleteBuffers(1, &o.ebo)
}

type Option func(*Object)

func WithAttribute(size int, normalized bool) Option {
	return func(obj *Object) {
		gl.EnableVertexAttribArray(obj.index)
		gl.VertexAttribPointer(obj.index, int32(size), gl.FLOAT, normalized, obj.stride, gl.PtrOffset(obj.offset*sizeOf[float32]()))
		obj.index++
		obj.offset += size
	}
}

func WithIndices(indices []uint32) Option {
	return func(obj *Object) {
		obj.indicesCount = int32(len(indices))
		gl.GenBuffers(1, &obj.ebo)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, obj.ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*sizeOf[uint32](), gl.Ptr(indices), gl.STATIC_DRAW)
	}
}

func sizeOf[T any]() int {
	var t T
	return int(unsafe.Sizeof(t))
}
