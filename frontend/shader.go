package frontend

import (
	"fmt"

	"golang.org/x/mobile/gl"
)

// TODO: Need a ShaderRegistry of somekind, ideally with support for default
// scene values vs per-shape values and attribute checking.
// TODO: Should each NewShader be a struct embedding a Program?

type Shader interface {
	Use()
	Close() error
	Attrib(string) gl.Attrib
	Uniform(string) gl.Uniform
	Context() gl.Context
}

func NewShader(glctx gl.Context, vertAsset, fragAsset string) (Shader, error) {
	program, err := LoadProgram(glctx, vertAsset, fragAsset)
	if err != nil {
		return nil, err
	}

	return &shader{
		glctx:    glctx,
		program:  program,
		attribs:  map[string]gl.Attrib{},
		uniforms: map[string]gl.Uniform{},
	}, nil
}

type shader struct {
	glctx   gl.Context
	program gl.Program

	attribs  map[string]gl.Attrib
	uniforms map[string]gl.Uniform
}

func (shader *shader) Context() gl.Context {
	return shader.glctx
}

func (shader *shader) Attrib(name string) gl.Attrib {
	v, ok := shader.attribs[name]
	if !ok {
		v = shader.glctx.GetAttribLocation(shader.program, name)
		shader.attribs[name] = v
	}
	return v
}

func (shader *shader) Uniform(name string) gl.Uniform {
	v, ok := shader.uniforms[name]
	if !ok {
		v = shader.glctx.GetUniformLocation(shader.program, name)
		shader.uniforms[name] = v
	}
	return v
}

func (shader *shader) Use() {
	shader.glctx.UseProgram(shader.program)
}

func (shader *shader) Close() error {
	shader.glctx.DeleteProgram(shader.program)
	return nil
}

type Shaders interface {
	Load(...string) error
	Get(string) Shader
	Reload() error
	Close() error
}

func ShaderLoader(glctx gl.Context) *shaderLoader {
	return &shaderLoader{
		glctx:   glctx,
		shaders: map[string]*shader{},
	}
}

type shaderLoader struct {
	glctx   gl.Context
	shaders map[string]*shader
}

func (loader *shaderLoader) Load(names ...string) error {
	for _, name := range names {
		s, err := NewShader(
			loader.glctx,
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
			loader.glctx,
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
