package glutils

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"path"
	"strconv"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// TextureObject 纹理数据结构体
type TextureObject struct {
	// 内部用
	ID     uint32 //? 文理ID
	Name   *uint8 //? 参数名称(材质)
	IsName *uint8 //? 参数名称(标记)
	UnitID uint32 //? 纹理单元
	Target uint32 //? 纹理类型
	// 外部信息
	File      string //? 纹理路径
	IsTexture bool   //? 是否生效

	wrap_s, wrap_t, min_f, mag_f int32
}

// New 创建
// * File   材质文件
// * UnitID 文理单元, 0~31 号可用
// * Target 纹理类型
func (T *TextureObject) New(File string, UnitID uint32, Target uint32) error {
	if UnitID > 32 {
		return fmt.Errorf("材质索引错误")
	}
	// 记录变量
	T.Target = Target
	T.UnitID = UnitID
	T.File = File

	/*
		gl.REPEAT	对纹理的默认行为。重复纹理图像。
		gl.MIRRORED_REPEAT	和GL_REPEAT一样，但每次重复图片是镜像放置的。
		gl.CLAMP_TO_EDGE	纹理坐标会被约束在0到1之间，超出的部分会重复纹理坐标的边缘，产生一种边缘被拉伸的效果。
		gl.CLAMP_TO_BORDER	超出的坐标为用户指定的边缘颜色。
	*/
	T.wrap_s = gl.CLAMP_TO_EDGE
	T.wrap_t = gl.CLAMP_TO_EDGE

	// gl.NEAREST（也叫邻近过滤，Nearest Neighbor Filtering）是OpenGL默认的纹理过滤方式。
	// gl.LINEAR（也叫线性过滤，(Bi)linear Filtering）它会基于纹理坐标附近的纹理像素，计算出一个插值，近似出这些纹理像素之间的颜色。
	T.min_f = gl.LINEAR
	T.mag_f = gl.LINEAR
	return nil
}

// 为当前绑定的纹理对象设置环绕、过滤方式
func (T *TextureObject) SetTexParamOptions(wrap_s, wrap_t, min_f, mag_f int32) {
	T.wrap_s, T.wrap_t, T.min_f, T.mag_f = wrap_s, wrap_t, min_f, mag_f
}

// * File   材质文件
// * UnitID 文理单元, 0~31 号可用
// * Target 纹理类型:gl.TEXTURE_2D
func NewTextureFromFile(file string, UnitID uint32, Target uint32) (*TextureObject, error) {

	textureXXX := &TextureObject{}
	var err = textureXXX.New(file, UnitID, Target)
	if err != nil {
		return nil, err
	}
	err = textureXXX.Init()
	return textureXXX, err
}

// Init 初始化
func (T *TextureObject) Init() error {
	// 初始化变量
	File, UnitID, Target := T.File, T.UnitID, T.Target
	T.Name = gl.Str("Texture[" + strconv.FormatUint(uint64(UnitID), 10) + "]\x00")
	T.IsName = gl.Str("IsTexture[" + strconv.FormatUint(uint64(UnitID), 10) + "]\x00")
	// 读文件
	imgdata, err := ioutil.ReadFile(File)
	if err != nil {
		return fmt.Errorf("材质文件 %q 打开失败: %v", File, err)
	}
	// 解码图片
	var img image.Image
	stamp := path.Ext(File) //? 得到文件名后戳
	if imgD, ok := ImageDecode[stamp]; ok {
		img, err = imgD(bytes.NewReader(imgdata))
		if err != nil {
			return fmt.Errorf("图片解码失败: " + err.Error())
		}
	} else {
		return fmt.Errorf("未知图片格式")
	}
	// 得到图片通道信息
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return fmt.Errorf("材质大小不支持")
	}
	// 转换格式
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	// 创建纹理
	gl.GenTextures(1, &T.ID)
	gl.BindTexture(Target, T.ID)
	// 纹理参数
	gl.TexParameteri(Target, gl.TEXTURE_MIN_FILTER, T.min_f)
	gl.TexParameteri(Target, gl.TEXTURE_MAG_FILTER, T.mag_f)
	gl.TexParameteri(Target, gl.TEXTURE_WRAP_S, T.wrap_s)
	gl.TexParameteri(Target, gl.TEXTURE_WRAP_T, T.wrap_t)
	// 添加纹理
	gl.TexImage2D(Target, 0, gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0, gl.RGBA, gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	gl.GenerateMipmap(Target)
	// 解除纹理
	gl.BindTexture(Target, 0)
	// 设置生效
	T.IsTexture = true
	return nil
}

// Delete 销毁
func (T *TextureObject) Delete() {
	gl.DeleteTextures(1, &T.ID)
}

// Update 更新
//? Program 着色器
func (T *TextureObject) Update(Program uint32) {
	if T.IsTexture {
		// 设置材质
		gl.ActiveTexture(TEXTURE0 + T.UnitID)                                 //? 激活纹理单元
		gl.BindTexture(T.Target, T.ID)                                        //? 绑定纹理
		gl.Uniform1i(gl.GetUniformLocation(Program, T.Name), int32(T.UnitID)) //? 设置材质
		// 设置材质参数
		gl.Uniform1i(gl.GetUniformLocation(Program, T.IsName), 1) //? 设置材质是否生效
	} else {
		gl.Uniform1i(gl.GetUniformLocation(Program, T.IsName), 0) //? 设置材质是否生效
	}
}

func (tex *TextureObject) Bind(texUnit uint32) {
	gl.ActiveTexture(texUnit)
	gl.BindTexture(tex.Target, tex.ID)
	tex.UnitID = texUnit
}

func (tex *TextureObject) UnBind() {
	tex.UnitID = 0
	gl.BindTexture(tex.Target, 0)
}

func (tex *TextureObject) SetUniform(uniformLoc int32) error {
	if tex.UnitID == 0 {
		return fmt.Errorf("texture not bound")
	}
	gl.Uniform1i(uniformLoc, int32(tex.UnitID-gl.TEXTURE0))
	return nil
}
