package graphics

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"log"
	"strings"
)

type ShaderType uint32

const (
	VERTEX   ShaderType = gl.VERTEX_SHADER
	GEOMETRY ShaderType = gl.GEOMETRY_SHADER
	FRAGMENT ShaderType = gl.FRAGMENT_SHADER
)

type ShaderProgram struct {
	id       uint32
	uniforms map[string]int32
}

func NewDefaultShaderProgram() *ShaderProgram {
	s := ShaderProgram{}
	s.id = gl.CreateProgram()

	s.AttachShader(VertexShaderBase, VERTEX)
	s.AttachShader(FragmentShaderSolidColor, FRAGMENT)

	s.Link()
	return &s
}

func NewShaderProgram(vertSource string, geomSource string, fragSource string) *ShaderProgram {
	s := ShaderProgram{}
	s.id = gl.CreateProgram()

	if vertSource != "" {
		s.AttachShader(vertSource, VERTEX)
	}
	if geomSource != "" {
		s.AttachShader(geomSource, GEOMETRY)
	}
	if fragSource != "" {
		s.AttachShader(fragSource, FRAGMENT)
	}

	s.Link()
	return &s
}

func (s *ShaderProgram) Release() {
	if s.id == 0 {
		log.Panicf("Trying to release a non initialized shader program")
	}
	// TODO
	//var shadersId [8]uint32
	//shaders_id := gl.GetAttachedShaders(s.id, 8, 8, &shadersId )
	//for id in  shaders_id:
	//	gl.DetachShader(self._program_id, shader_id)
	//	gl.DeleteShader(shader_id)

	gl.DeleteProgram(s.id)
}

func (s *ShaderProgram) AttachShader(source string, shaderType ShaderType) {
	shaderId := gl.CreateShader(uint32(shaderType))
	cSource, free := gl.Strs(source)
	gl.ShaderSource(shaderId, 1, cSource, nil)
	free()
	gl.CompileShader(shaderId)

	var status int32
	gl.GetShaderiv(shaderId, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shaderId, gl.INFO_LOG_LENGTH, &logLength)

		logStr := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shaderId, logLength, nil, gl.Str(logStr))

		log.Panicf("failed to compile %v: %v", source, logStr)
	}
	gl.AttachShader(s.id, shaderId)
}

func (s *ShaderProgram) Link() {
	gl.LinkProgram(s.id)
	var status int32
	gl.GetProgramiv(s.id, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(s.id, gl.INFO_LOG_LENGTH, &logLength)

		logStr := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(s.id, logLength, nil, gl.Str(logStr))

		log.Panicf("failed to link program: %v", logStr)
	}
}

func (s *ShaderProgram) Id() uint32 {
	return s.id
}

func (s *ShaderProgram) GetUniform(name string) int32 {
	if s.uniforms == nil {
		s.uniforms = make(map[string]int32)
	}

	uniform, ok := s.uniforms[name]
	if !ok {
		cname := gl.Str(name + "\x00")
		uniform = gl.GetUniformLocation(s.id, cname)
		s.uniforms[name] = uniform
	}

	return uniform
}

func (s *ShaderProgram) SetUniform(name string, val interface{}) {
	uniform := s.GetUniform(name)
	switch v := val.(type) {
	case *float32:
		gl.Uniform1fv(uniform, 1, v)
	case *mgl32.Vec2:
		gl.Uniform2fv(uniform, 1, &(*v)[0])
	case *mgl32.Vec3:
		gl.Uniform3fv(uniform, 1, &(*v)[0])
	case *mgl32.Vec4:
		gl.Uniform4fv(uniform, 1, &(*v)[0])
	case *mgl32.Mat2:
		gl.UniformMatrix2fv(uniform, 1, false, &(*v)[0])
	case *mgl32.Mat3:
		gl.UniformMatrix3fv(uniform, 1, false, &(*v)[0])
	case *mgl32.Mat4:
		gl.UniformMatrix4fv(uniform, 1, false, &(*v)[0])
	case *Color:
		gl.Uniform4fv(uniform, 1, &(*v)[0])
	default:
		log.Panicf("unknown value type: %T %+v", val, val)
	}
}

const (
	VertexShaderBase = `
        #version 410 core

        uniform mat4 model;
        uniform mat4 projection;

        layout(location=0) in vec2 vertex;
        layout(location=1) in vec2 uv;

        out vec2 uv_out;

        void main() {
            vec4 vertex_world = model * vec4(vertex, 0, 1);
            gl_Position = projection * vertex_world;
            uv_out = uv;
        }
        ` + "\x00"

	FragmentShaderSolidColor = `
        #version 410 core

        in vec2 uv_out;
        out vec4 out_color;
        uniform vec4 color;

        uniform sampler2D tex;

        void main() {
            out_color = color;
        }
        ` + "\x00"

	FragmentShaderTexture = `
        #version 410 core

        in vec2 uv_out;
        out vec4 color;

        uniform sampler2D tex;

        void main() {
			if(texture(tex, uv_out).a != 1.0f)
			{
        		discard;
    		}
            color = texture(tex, uv_out);
        }
        ` + "\x00"
)
