package frontend

import (
	"time"

	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/shazow/audioscopic/frontend/control"
	"github.com/shazow/audioscopic/frontend/loader"
)

type Vector interface {
	Position() mgl.Vec3
	Direction() mgl.Vec3
}

type World interface {
	Scene

	Reset()
	Tick(time.Duration) error
	Focus() Vector

	Start(control.Bindings, loader.Shaders, loader.Textures) error
}

func FixedVector(position mgl.Vec3, direction mgl.Vec3) Vector {
	return &fixedVector{
		position:  position,
		direction: direction,
	}
}

type fixedVector struct {
	position  mgl.Vec3
	direction mgl.Vec3
}

func (v fixedVector) Position() mgl.Vec3  { return v.position }
func (v fixedVector) Direction() mgl.Vec3 { return v.direction }

type stubWorld struct {
	Scene
}

func (w stubWorld) Reset()                                                              {}
func (w stubWorld) Tick(d time.Duration) error                                          { return nil }
func (w stubWorld) Focus() Vector                                                       { return fixedVector{} }
func (w stubWorld) Start(_ control.Bindings, _ loader.Shaders, _ loader.Textures) error { return nil }

func StubWorld(scene Scene) World {
	return stubWorld{scene}
}
