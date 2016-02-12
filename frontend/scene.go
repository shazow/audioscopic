package frontend

import (
	"fmt"

	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/shazow/audioscopic/frontend/camera"
)

type Light struct {
	color    mgl.Vec3
	position mgl.Vec3
}

func (light *Light) MoveTo(position mgl.Vec3) {
	light.position = position
}

type Drawable interface {
	Draw(camera.Camera)
	Transform(*mgl.Mat4) mgl.Mat4
	UseShader(Shader) (Shader, bool)
}

// TODO: node tree with transforms
type Node struct {
	Shape
	transform *mgl.Mat4
	shader    Shader
}

func (node *Node) Draw(cam camera.Camera) {
	node.Shape.Draw(node.shader, cam)
}

func (node *Node) UseShader(parent Shader) (Shader, bool) {
	if node.shader == nil || node.shader == parent {
		node.shader = parent
		return parent, false
	}
	node.shader.Use()
	return node.shader, true
}

func (node *Node) Transform(parent *mgl.Mat4) mgl.Mat4 {
	return MultiMul(node.transform, parent)
}

func (node *Node) String() string {
	return fmt.Sprintf("<Shape of %d vertices; transform: %v>", node.Len(), node.transform)
}

type Scene interface {
	Add(interface{})
	Draw(camera.Camera)
	String() string
}

func NewScene() Scene {
	return &sliceScene{
		nodes: []Drawable{},
	}
}

type sliceScene struct {
	nodes     []Drawable
	transform *mgl.Mat4
}

func (scene *sliceScene) String() string {
	return fmt.Sprintf("%d nodes", len(scene.nodes))
}

func (scene *sliceScene) Add(item interface{}) {
	scene.nodes = append(scene.nodes, item.(Drawable))
}

func (scene *sliceScene) Draw(cam camera.Camera) {
	// Setup MVP
	projection, view, position := cam.Projection(), cam.View(), cam.Position()

	var parentShader Shader
	for _, node := range scene.nodes {
		shader, changed := node.UseShader(parentShader)
		glctx := shader.Context()

		if changed {
			// TODO: Pre-load these into relevant shaders?
			glctx.UniformMatrix4fv(shader.Uniform("cameraPos"), position[:])
			glctx.UniformMatrix4fv(shader.Uniform("view"), view[:])
			glctx.UniformMatrix4fv(shader.Uniform("projection"), projection[:])
		}

		// TODO: Move these into node.Draw?
		model := node.Transform(scene.transform)
		normal := model.Mul4(view).Inv().Transpose()

		// Camera space
		glctx.UniformMatrix4fv(shader.Uniform("model"), model[:])
		glctx.UniformMatrix4fv(shader.Uniform("normalMatrix"), normal[:])

		node.Draw(cam)
	}
}
