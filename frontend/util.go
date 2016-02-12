package frontend

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	image_draw "image/draw"
	_ "image/png"
	"io/ioutil"
	"log"

	mgl "github.com/go-gl/mathgl/mgl32"

	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/gl"
)

type dimslice_float32 struct {
	dim   int
	slice []float32
}

func (o dimslice_float32) Slice(i, j int) interface{} { return o.slice[i:j] }
func (o dimslice_float32) Dim() int                   { return o.dim }
func (o dimslice_float32) String() string {
	return fmt.Sprintf("<float32 slice: len=%d dim=%d>", len(o.slice), o.dim)
}

type dimslice_uint8 struct {
	dim   int
	slice []uint8
}

func (o dimslice_uint8) Slice(i, j int) interface{} { return o.slice[i:j] }
func (o dimslice_uint8) Dim() int                   { return o.dim }
func (o dimslice_uint8) String() string {
	return fmt.Sprintf("<uint8 slice: len=%d dim=%d>", len(o.slice), o.dim)
}

func NewDimSlice(dim int, slice interface{}) DimSlicer {
	switch slice := slice.(type) {
	case []float32:
		return &dimslice_float32{dim, slice}
	case []uint8:
		return &dimslice_uint8{dim, slice}
	}
	panic(fmt.Sprintf("invalid slice type: %T", slice))
	return nil
}

type DimSlicer interface {
	Slice(int, int) interface{}
	Dim() int
	String() string
}

// EncodeObjects converts float32 vertices into a LittleEndian byte array.
// Offset and length are based on the number of rows per dimension.
// TODO: Replace with https://github.com/lunixbochs/struc?
func EncodeObjects(offset int, length int, objects ...DimSlicer) []byte {
	//log.Println("EncodeObjects:", offset, length, objects)
	// TODO: Pre-allocate? Use a SyncPool?
	/*
		dimSum := 0 // yum!
		for _, obj := range objects {
			dimSum += obj.Dim()
		}
		v := make([]float32, dimSum*length)
	*/

	buf := bytes.Buffer{}

	for i := offset; i < length; i++ {
		for _, obj := range objects {
			data := obj.Slice(i*obj.Dim(), (i+1)*obj.Dim())
			if err := binary.Write(&buf, binary.LittleEndian, data); err != nil {
				panic(fmt.Sprintln("binary.Write failed:", err))
			}
		}
	}
	//fmt.Printf("Wrote %d vertices: %d to %d \t", shape.Len()-n, n, shape.Len())
	//fmt.Println(wrote)

	return buf.Bytes()
}

func loadAsset(name string) ([]byte, error) {
	f, err := asset.Open(name)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}

func loadShader(glctx gl.Context, shaderType gl.Enum, assetName string) (gl.Shader, error) {
	// Borrowed from golang.org/x/mobile/exp/gl/glutil
	src, err := loadAsset(assetName)
	if err != nil {
		return gl.Shader{}, err
	}

	shader := glctx.CreateShader(shaderType)
	if shader.Value == 0 {
		return gl.Shader{}, fmt.Errorf("glutil: could not create shader (type %v)", shaderType)
	}
	glctx.ShaderSource(shader, string(src))
	glctx.CompileShader(shader)
	if glctx.GetShaderi(shader, gl.COMPILE_STATUS) == 0 {
		defer glctx.DeleteShader(shader)
		return gl.Shader{}, fmt.Errorf("shader compile: %s", glctx.GetShaderInfoLog(shader))
	}
	return shader, nil
}

func LoadShaders(glctx gl.Context, program gl.Program, vertexAsset, fragmentAsset string) error {
	vertexShader, err := loadShader(glctx, gl.VERTEX_SHADER, vertexAsset)
	if err != nil {
		return err
	}
	fragmentShader, err := loadShader(glctx, gl.FRAGMENT_SHADER, fragmentAsset)
	if err != nil {
		glctx.DeleteShader(vertexShader)
		return err
	}

	if glctx.GetProgrami(program, gl.ATTACHED_SHADERS) > 0 {
		for _, shader := range glctx.GetAttachedShaders(program) {
			glctx.DetachShader(program, shader)
		}
	}

	glctx.AttachShader(program, vertexShader)
	glctx.AttachShader(program, fragmentShader)
	glctx.LinkProgram(program)

	// Flag shaders for deletion when program is unlinked.
	glctx.DeleteShader(vertexShader)
	glctx.DeleteShader(fragmentShader)

	if glctx.GetProgrami(program, gl.LINK_STATUS) == 0 {
		defer glctx.DeleteProgram(program)
		return fmt.Errorf("LoadShaders: %s", glctx.GetProgramInfoLog(program))
	}
	return nil
}

// LoadProgram reads shader sources from the asset repository, compiles, and
// links them into a program.
func LoadProgram(glctx gl.Context, vertexAsset, fragmentAsset string) (program gl.Program, err error) {
	log.Println("LoadProgram:", vertexAsset, fragmentAsset)

	program = glctx.CreateProgram()
	if program.Value == 0 {
		return gl.Program{}, fmt.Errorf("glutil: no programs available")
	}

	err = LoadShaders(glctx, program, vertexAsset, fragmentAsset)
	return
}

// LoadTexture2D reads and decodes an image from the asset repository and creates
// a texture object based on the full dimensions of the image.
func LoadTexture2D(glctx gl.Context, name string) (tex gl.Texture, err error) {
	imgFile, err := asset.Open(name)
	if err != nil {
		return
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return
	}

	rgba := image.NewRGBA(img.Bounds())
	image_draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, image_draw.Src)

	tex = glctx.CreateTexture()
	glctx.ActiveTexture(gl.TEXTURE0)
	glctx.BindTexture(gl.TEXTURE_2D, tex)
	glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	glctx.TexImage2D(
		gl.TEXTURE_2D,
		0,
		rgba.Rect.Size().X,
		rgba.Rect.Size().Y,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		rgba.Pix)
	return
}

// LoadTextureCube reads and decodes an image from the asset repository and creates
// a texture cube map object based on the full dimensions of the image.
func LoadTextureCube(glctx gl.Context, name string) (tex gl.Texture, err error) {
	imgFile, err := asset.Open(name)
	if err != nil {
		return
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return
	}

	rgba := image.NewRGBA(img.Bounds())
	image_draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, image_draw.Src)

	tex = glctx.CreateTexture()
	glctx.ActiveTexture(gl.TEXTURE0)
	glctx.BindTexture(gl.TEXTURE_CUBE_MAP, tex)

	target := gl.TEXTURE_CUBE_MAP_POSITIVE_X
	for i := 0; i < 6; i++ {
		// TODO: Load atlas, not the same image
		glctx.TexImage2D(
			gl.Enum(target+i),
			0,
			rgba.Rect.Size().X,
			rgba.Rect.Size().Y,
			gl.RGBA,
			gl.UNSIGNED_BYTE,
			rgba.Pix,
		)
	}

	glctx.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	glctx.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	glctx.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	glctx.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	// Not available in GLES 2.0 :(
	//gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)

	return
}

// MultiMul multiplies every non-nil Mat4 reference and returns the result. If
// none are given, then it returns the identity matrix.
func MultiMul(matrices ...*mgl.Mat4) mgl.Mat4 {
	var r mgl.Mat4
	ok := false
	for _, m := range matrices {
		if m == nil {
			continue
		}
		if !ok {
			r = *m
			ok = true
			continue
		}
		r = r.Mul4(*m)
	}
	if ok {
		return r
	}
	return mgl.Ident4()
}

func Quad(a mgl.Vec3, b mgl.Vec3) []float32 {
	return []float32{
		// First triangle
		b[0], b[1], b[2], // Top Right
		a[0], b[1], a[2], // Top Left
		a[0], a[1], a[2], // Bottom Left
		// Second triangle
		a[0], a[1], a[2], // Bottom Left
		b[0], b[1], b[2], // Top Right
		b[0], a[1], b[2], // Bottom Right
	}
}

func Upvote(tip mgl.Vec3, size float32) []float32 {
	a := tip.Add(mgl.Vec3{-size / 2, -size * 2, 0})
	b := tip.Add(mgl.Vec3{size / 2, -size, 0})
	return []float32{
		tip[0], tip[1], tip[2], // Top
		tip[0] - size, tip[1] - size, tip[2], // Bottom left
		tip[0] + size, tip[1] - size, tip[2], // Bottom right

		// Arrow handle
		b[0], b[1], b[2], // Top Right
		a[0], b[1], a[2], // Top Left
		a[0], a[1], a[2], // Bottom Left
		a[0], a[1], a[2], // Bottom Left
		b[0], b[1], b[2], // Top Right
		b[0], a[1], b[2], // Bottom Right
	}
}
