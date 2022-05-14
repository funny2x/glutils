package glutils

///
// 文理数据加载包
//  材质加载
// ? 日志
///
import (
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/funny2x/glutils/image/bmp"
	"github.com/funny2x/glutils/image/tga"
)

// 纹理变量
const (
	// 纹理单元
	TEXTURE   = 0x1702
	TEXTURE0  = 0x84C0
	TEXTURE1  = 0x84C1
	TEXTURE10 = 0x84CA
	TEXTURE11 = 0x84CB
	TEXTURE12 = 0x84CC
	TEXTURE13 = 0x84CD
	TEXTURE14 = 0x84CE
	TEXTURE15 = 0x84CF
	TEXTURE16 = 0x84D0
	TEXTURE17 = 0x84D1
	TEXTURE18 = 0x84D2
	TEXTURE19 = 0x84D3
	TEXTURE2  = 0x84C2
	TEXTURE20 = 0x84D4
	TEXTURE21 = 0x84D5
	TEXTURE22 = 0x84D6
	TEXTURE23 = 0x84D7
	TEXTURE24 = 0x84D8
	TEXTURE25 = 0x84D9
	TEXTURE26 = 0x84DA
	TEXTURE27 = 0x84DB
	TEXTURE28 = 0x84DC
	TEXTURE29 = 0x84DD
	TEXTURE3  = 0x84C3
	TEXTURE30 = 0x84DE
	TEXTURE31 = 0x84DF
	TEXTURE4  = 0x84C4
	TEXTURE5  = 0x84C5
	TEXTURE6  = 0x84C6
	TEXTURE7  = 0x84C7
	TEXTURE8  = 0x84C8
	TEXTURE9  = 0x84C9
	// 纹理类型
	TEXTURE1D                 = 0x0DE0
	TEXTURE1DARRAY            = 0x8C18
	TEXTURE2D                 = 0x0DE1
	TEXTURE2DARRAY            = 0x8C1A
	TEXTURE2DMULTISAMPLE      = 0x9100
	TEXTURE2DMULTISAMPLEARRAY = 0x9102
	TEXTURE3D                 = 0x806F
)

// ImageDecode 图片解码
var ImageDecode map[string]func(r io.Reader) (image.Image, error)

func init() {
	ImageDecode = make(map[string]func(r io.Reader) (image.Image, error))
	ImageDecode[".png"] = png.Decode
	ImageDecode[".jpg"] = jpeg.Decode
	ImageDecode[".tga"] = tga.Decode
	ImageDecode[".bmp"] = bmp.Decode
}
