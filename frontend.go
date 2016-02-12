package main

import (
	"time"

	"github.com/shazow/audioscopic/frontend"
	"github.com/shazow/audioscopic/frontend/control"
	"github.com/shazow/audioscopic/frontend/loader"

	mgl "github.com/go-gl/mathgl/mgl32"
)

func startFrontend() {
	world := newWorld()
	engine := frontend.NewEngine(world)

	frontend.StartMobile(engine)
}

func newWorld() frontend.World {
	scene := frontend.NewScene()
	return &world{
		Scene: scene,
	}
}

type world struct {
	frontend.Scene
}

func (w *world) Start(bindings control.Bindings, shaders loader.Shaders, textures loader.Textures) error {
	// Load shaders
	err := shaders.Load("skybox", "main")
	if err != nil {
		return err
	}

	// Load textures
	err = textures.Load("square.png")
	if err != nil {
		return err
	}

	// Make skybox
	skybox := frontend.NewSkybox(shaders.Get("skybox"), textures.GetCube("square.png"))
	w.Add(skybox)

	floor := frontend.NewFloor(shaders.Get("main"))
	w.Add(floor)
	return nil
}

func (w *world) Reset()                     {}
func (w *world) Tick(d time.Duration) error { return nil }
func (w *world) Focus() frontend.Vector {
	return frontend.FixedVector(mgl.Vec3{0, 1, 1}, mgl.Vec3{0, 1, 1})
}
