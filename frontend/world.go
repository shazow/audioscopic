package frontend

import (
	"time"

	mgl "github.com/go-gl/mathgl/mgl32"
)

type Vector interface {
	Position() mgl.Vec3
	Direction() mgl.Vec3
}

type fixedVector struct {
	position  mgl.Vec3
	direction mgl.Vec3
}

func (v fixedVector) Position() mgl.Vec3  { return v.position }
func (v fixedVector) Direction() mgl.Vec3 { return v.direction }

type World interface {
	Reset()
	Tick(time.Duration) error
	Focus() Vector
}

type stubWorld struct{}

func (w *stubWorld) Reset()                     {}
func (w *stubWorld) Tick(d time.Duration) error { return nil }
func (w *stubWorld) Focus() Vector              { return fixedVector{} }
