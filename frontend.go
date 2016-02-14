package main

import (
	"log"
	"time"

	"github.com/shazow/audioscopic/frontend"

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

	particles frontend.Emitter
}

func (w *world) Start(ctx frontend.WorldContext) error {
	// Load shaders
	err := ctx.Shaders.Load("main", "skybox", "particle")
	if err != nil {
		return err
	}

	// Load textures
	err = ctx.Textures.Load("square.png")
	if err != nil {
		return err
	}

	// Make skybox
	skybox := frontend.NewSkybox(ctx.Shaders.Get("skybox"), ctx.Textures.GetCube("square.png"))
	w.Add(skybox)

	floor := frontend.NewFloor(ctx.Shaders.Get("main"))
	w.Add(floor)

	log.Println("scene", w)
	//emitter := frontend.ParticleEmitter
	return nil
}

func (w *world) Tick(d time.Duration) error {
	return nil
}

func (w *world) Reset() {}
func (w *world) Focus() frontend.Vector {
	return frontend.FixedVector(mgl.Vec3{0, 0, 0}, mgl.Vec3{0, 0, 1})
}
