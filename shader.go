package glutils

// import "C"
import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func BasicProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	s, e := NewShaderObject(vertexShaderSource, fragmentShaderSource)
	return s.ID, e
}

func NewShader(vertFile, fragFile, geomFile string) (*Shader, error) {
	vertSrc, err := readFile(vertFile)
	if err != nil {
		return nil, err
	}

	fragSrc, err := readFile(fragFile)
	if err != nil {
		return nil, err
	}

	var geomSrc []byte
	if geomFile != "" {
		geomSrc, err = readFile(geomFile)
		if err != nil {
			return nil, err
		}
	}

	p, err := createProgram(vertSrc, fragSrc, geomSrc)
	if err != nil {
		return nil, err
	}

	return setupShader(p), nil
}

func setupShader(program uint32) *Shader {
	var (
		c, b, s int32
		i       uint32
		n       uint8
	)
	b = 255
	uniforms := make(map[string]int32)
	attributes := make(map[string]uint32)

	gl.GetProgramiv(program, gl.ACTIVE_UNIFORMS, &c)
	for i = 0; i < uint32(c); i++ {
		gl.GetActiveUniform(program, i, b, nil, &s, nil, &n)
		loc := gl.GetUniformLocation(program, &n)
		name := gl.GoStr(&n)
		fmt.Println(name, loc)
		uniforms[name] = loc
	}
	fmt.Println("---")
	gl.GetProgramiv(program, gl.ACTIVE_ATTRIBUTES, &c)
	for i = 0; i < uint32(c); i++ {
		gl.GetActiveAttrib(program, i, b, nil, nil, nil, &n)
		loc := gl.GetAttribLocation(program, &n)
		name := gl.GoStr(&n)
		fmt.Println(name, loc)
		attributes[name] = uint32(loc)
	}

	return &Shader{
		Program:    program,
		Uniforms:   uniforms,
		Attributes: attributes,
	}
}

func createProgram(v, f, g []byte) (uint32, error) {
	vertex, err := generateCompileShader(string(v), gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	frag, err := generateCompileShader(string(f), gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	var geom uint32
	use_geom := false
	if len(g) > 0 {
		geom, err = generateCompileShader(string(g), gl.GEOMETRY_SHADER)
		use_geom = true
		if err != nil {
			return 0, err
		}
	}

	p, err := linkProgram(vertex, frag, geom, use_geom)
	if err != nil {
		return 0, err
	}

	defer deleteShader(p, vertex)
	defer deleteShader(p, frag)
	if use_geom {
		defer deleteShader(p, geom)
	}

	return p, nil
}

func deleteShader(p, s uint32) {
	gl.DetachShader(p, s)
	gl.DeleteShader(s)
}

func linkProgram(v, f, g uint32, use_geom bool) (uint32, error) {
	program := gl.CreateProgram()
	gl.AttachShader(program, v)
	gl.AttachShader(program, f)
	if use_geom {
		gl.AttachShader(program, g)
	}

	gl.LinkProgram(program)
	// check for program linking errors
	return program, checkCompileErrors(program)
}

type Shader struct {
	Program    uint32
	Uniforms   map[string]int32
	Attributes map[string]uint32
}

type VertexArray struct {
	Data          []float32
	Indices       []uint32
	Stride        int32
	Normalized    bool
	DrawMode      uint32
	Attributes    map[uint32]int32 //map attrib loc to size
	Vao, vbo, ebo uint32
}

func (v *VertexArray) Setup() {
	gl.GenVertexArrays(1, &v.Vao)
	gl.GenBuffers(1, &v.vbo)

	gl.BindVertexArray(v.Vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, v.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(v.Data)*GL_FLOAT32_SIZE, gl.Ptr(v.Data), v.DrawMode)

	if len(v.Indices) > 0 {
		gl.GenBuffers(1, &v.ebo)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, v.ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(v.Indices)*GL_FLOAT32_SIZE, gl.Ptr(v.Indices), v.DrawMode)
	}

	i := 0
	for loc, size := range v.Attributes {
		gl.VertexAttribPointer(loc, size, gl.FLOAT, v.Normalized, v.Stride*GL_FLOAT32_SIZE, gl.PtrOffset(i*GL_FLOAT32_SIZE))
		gl.EnableVertexAttribArray(loc)
		i += int(size)
	}
	gl.BindVertexArray(0)
}
