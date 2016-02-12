package main

import (
	"time"

	"github.com/shazow/audioscopic/frontend"
	"github.com/shazow/audioscopic/frontend/control"
	"github.com/shazow/audioscopic/frontend/loader"
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
	err := shaders.Load("skybox")
	if err != nil {
		return err
	}

	// Load textures
	err = textures.Load("square.png")
	if err != nil {
		return err
	}

	// Make skybox
	// TODO: Add closer, or use a texture loader
	w.Add(frontend.NewSkybox(shaders.Get("skybox"), textures.GetCube("square.png")))
	return nil
}

func (w *world) Reset()                     {}
func (w *world) Tick(d time.Duration) error { return nil }
func (w *world) Focus() frontend.Vector     { return frontend.FixedVector{} }
