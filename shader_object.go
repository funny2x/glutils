/*
创建着色器对象
将源码字符串赋予着色器对象
编译着色器

创建着色器程序对象
将编译好的着色器附加到程序对象上
链接生成程序
*/

package glutils

import (
	"github.com/funny2x/mgl32"
	"github.com/go-gl/gl/v4.1-core/gl"
)

// Shader 着色器对象
type ShaderObject struct {
	ID uint32
}

//NewShader 结构体构造函数
func NewShaderObject(vertShaderPath, fragShaderPath string) (*ShaderObject, error) {
	//读出源码
	vertShaderCode, err := getShaderFromFile(vertShaderPath)
	if err != nil {
		return nil, err
	}
	fragShaderCode, err := getShaderFromFile(fragShaderPath)
	if err != nil {
		return nil, err
	}
	//生成编译着色器对象
	vertexShader, err := generateCompileShader(vertShaderCode, gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}
	fragmentShader, err := generateCompileShader(fragShaderCode, gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, err
	}
	// 删除着色器
	defer func() {
		gl.DeleteShader(vertexShader)
		gl.DeleteShader(fragmentShader)
	}()
	//链接生成着色器程序
	shaderProgram, err := linkShader(vertexShader, fragmentShader)
	return &ShaderObject{
		ID: shaderProgram,
	}, err
}

//Use 激活着色器
func (s *ShaderObject) Use() {
	gl.UseProgram(s.ID)
}

//SetBool 赋 bool 类型值给着色器程序中的uniform
func (s *ShaderObject) SetBool(name string, value bool) {
	if value {
		s.SetInt(name, 1)
	} else {
		s.SetInt(name, 0)
	}
}

//SetInt 赋 int 类型值给着色器程序中的uniform
func (s *ShaderObject) SetInt(name string, value int32) {
	gl.Uniform1i(s.GetUniform(name), value)
}

//SetFloat 赋 float 类型值给着色器程序中的uniform
func (s *ShaderObject) SetFloat(name string, value float32) {
	gl.Uniform1f(s.GetUniform(name), value)
}

//SetVec2XY 赋 Vec2(X,Y) 类型值给着色器程序中的uniform
func (s *ShaderObject) SetVec2XY(name string, x, y float32) {
	gl.Uniform2f(s.GetUniform(name), x, y)
}

//SetVec2 赋 Vec2 类型值给着色器程序中的uniform
func (s *ShaderObject) SetVec2(name string, value mgl32.Vec2) {
	s.SetVec2XY(name, value.X(), value.Y())
}

//SetVec3XYZ 赋 Vec3(X,Y,Z) 类型值给着色器程序中的uniform
func (s *ShaderObject) SetVec3XYZ(name string, x, y, z float32) {
	gl.Uniform3f(s.GetUniform(name), x, y, z)
}

//SetVec3 赋 Vec3 类型值给着色器程序中的uniform
func (s *ShaderObject) SetVec3(name string, value mgl32.Vec3) {
	s.SetVec3XYZ(name, value.X(), value.Y(), value.Z())
}

//SetVec4XYZW 赋 Vec4(X,Y,Z,W) 类型值给着色器程序中的uniform
func (s *ShaderObject) SetVec4XYZW(name string, x, y, z, w float32) {
	gl.Uniform4f(s.GetUniform(name), x, y, z, w)
}

//SetVec4 赋 Vec4 类型值给着色器程序中的uniform
func (s *ShaderObject) SetVec4(name string, value mgl32.Vec4) {
	s.SetVec4XYZW(name, value.X(), value.Y(), value.Z(), value.W())
}

//SetMat2 赋 Mat2 类型值给着色器程序中的uniform
func (s *ShaderObject) SetMat2(name string, value mgl32.Mat2) {
	gl.UniformMatrix2fv(s.GetUniform(name), 1, false, &value[0])
}

//SetMat3 赋 Mat3 类型值给着色器程序中的uniform
func (s *ShaderObject) SetMat3(name string, value mgl32.Mat3) {
	gl.UniformMatrix3fv(s.GetUniform(name), 1, false, &value[0])
}

//SetMat4 赋 Mat4 类型值给着色器程序中的uniform
func (s *ShaderObject) SetMat4(name string, value mgl32.Mat4) {
	gl.UniformMatrix4fv(s.GetUniform(name), 1, false, &value[0])
}

//Delete 删除着色器程序
func (s *ShaderObject) Delete() {
	gl.DeleteShader(s.ID)
}

//GetUniform 赋值给着色器程序中的uniform
func (s *ShaderObject) GetUniform(name string) int32 {
	return s.getUniform(name)
}
