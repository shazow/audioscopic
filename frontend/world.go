package frontend

import (
	"time"

	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/shazow/audioscopic/frontend/control"
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

	Start(control.Bindings, Shaders, Textures) error
}

type FixedVector struct {
	position  mgl.Vec3
	direction mgl.Vec3
}

func (v FixedVector) Position() mgl.Vec3  { return v.position }
func (v FixedVector) Direction() mgl.Vec3 { return v.direction }

type stubWorld struct {
	Scene
}

func (w stubWorld) Reset()                                                {}
func (w stubWorld) Tick(d time.Duration) error                            { return nil }
func (w stubWorld) Focus() Vector                                         { return FixedVector{} }
func (w stubWorld) Start(_ control.Bindings, _ Shaders, _ Textures) error { return nil }

func StubWorld(scene Scene) World {
	return stubWorld{scene}
}
