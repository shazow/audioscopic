package loader

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"golang.org/x/mobile/asset"
)

// TODO: Need a ShaderRegistry of somekind, ideally with support for default
// scene values vs per-shape values and attribute checking.
// TODO: Should each NewShader be a struct embedding a Program?

type Shader interface {
	Use()
	Close() error
	Attrib(string) int32
	Uniform(string) int32
	Program() uint32
}

func NewShader(vertAsset, fragAsset string) (Shader, error) {
	program, err := LoadProgram(vertAsset, fragAsset)
	if err != nil {
		return nil, err
	}

	return &shader{
		program:  program,
		attribs:  map[string]int32{},
		uniforms: map[string]int32{},
	}, nil
}

type shader struct {
	program uint32

	attribs  map[string]int32
	uniforms map[string]int32
}

func (shader *shader) Attrib(name string) int32 {
	v, ok := shader.attribs[name]
	if !ok {
		v = gl.GetAttribLocation(shader.program, gl.Str(name+"\x00"))
		shader.attribs[name] = v
		log.Println(name, "->", v)
	}
	return v
}

func (shader *shader) Uniform(name string) int32 {
	v, ok := shader.uniforms[name]
	if !ok {
		v = gl.GetUniformLocation(shader.program, gl.Str(name+"\x00"))
		shader.uniforms[name] = v
		log.Println(name, "->", v)
	}
	return v
}

func (shader *shader) Use() {
	gl.UseProgram(shader.program)
}

func (shader *shader) Close() error {
	gl.DeleteProgram(shader.program)
	return nil
}

func (shader *shader) Program() uint32 {
	return shader.program
}

type Shaders interface {
	Load(...string) error
	Get(string) Shader
	Reload() error
	Close() error
}

func ShaderLoader() *shaderLoader {
	return &shaderLoader{
		shaders: map[string]*shader{},
	}
}

type shaderLoader struct {
	shaders map[string]*shader
}

func (loader *shaderLoader) Load(names ...string) error {
	for _, name := range names {
		s, err := NewShader(
			fmt.Sprintf("%s.v.glsl", name),
			fmt.Sprintf("%s.f.glsl", name),
		)
		if err != nil {
			return err
		}
		loader.shaders[name] = s.(*shader)
	}
	return nil
}

func (loader *shaderLoader) Get(name string) Shader {
	return loader.shaders[name]
}

func (loader *shaderLoader) Reload() error {
	for k, shader := range loader.shaders {
		err := LoadShaders(
			shader.program,
			fmt.Sprintf("%s.v.glsl", k),
			fmt.Sprintf("%s.f.glsl", k),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (loader *shaderLoader) Close() error {
	for _, shader := range loader.shaders {
		shader.Close()
	}
	return nil
}

func loadAsset(name string) ([]byte, error) {
	f, err := asset.Open(name)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}

func loadShader(shaderType uint32, assetName string) (uint32, error) {
	// Borrowed from golang.org/x/mobile/exp/gl/glutil
	src, err := loadAsset(assetName)
	if err != nil {
		return 0, err
	}

	shader := gl.CreateShader(shaderType)
	if shader == 0 {
		return 0, fmt.Errorf("glutil: could not create shader (type %v)", shaderType)
	}

	csources, free := gl.Strs(string(src) + "\x00")
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", src, log)
	}
	log.Println("+shader:", status)

	return shader, nil
}

func LoadShaders(program uint32, vertexAsset, fragmentAsset string) error {
	vertexShader, err := loadShader(gl.VERTEX_SHADER, vertexAsset)
	if err != nil {
		return err
	}
	fragmentShader, err := loadShader(gl.FRAGMENT_SHADER, fragmentAsset)
	if err != nil {
		gl.DeleteShader(vertexShader)
		return err
	}

	/*
		if gl.GetProgramiv(program, gl.ATTACHED_SHADERS) > 0 {
			for _, shader := range gl.GetAttachedShaders(program) {
				gl.DetachShader(program, shader)
			}
		}
	*/

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	// Flag shaders for deletion when program is unlinked.
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		defer gl.DeleteProgram(program)
		var logLength int32
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))
		return fmt.Errorf("failed to link program: %v", log)
	}
	log.Println("LoadShaders:", program, "->", status)

	return nil
}

// LoadProgram reads shader sources from the asset repository, compiles, and
// links them into a program.
func LoadProgram(vertexAsset, fragmentAsset string) (program uint32, err error) {
	program = gl.CreateProgram()
	if program == gl.FALSE {
		return program, fmt.Errorf("gl: no programs available")
	}

	log.Println("LoadProgram:", vertexAsset, fragmentAsset, "->", program)

	err = LoadShaders(program, vertexAsset, fragmentAsset)
	return
}
