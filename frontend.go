package main

import (
	"time"

	"github.com/shazow/go-gameblocks"

	mgl "github.com/go-gl/mathgl/mgl32"
)

func startFrontend() {
	world := newWorld()
	engine := gameblocks.NewEngine(world)

	gameblocks.StartMobile(engine)
}

func newWorld() gameblocks.World {
	scene := gameblocks.NewScene()
	return &world{
		Scene: scene,
	}
}

type world struct {
	gameblocks.Scene

	particles gameblocks.Emitter
}

func (w *world) Start(ctx gameblocks.WorldContext) error {
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
	skybox := gameblocks.NewSkybox(ctx.Shaders.Get("skybox"), ctx.Textures.GetCube("square.png"))
	w.Add(skybox)

	floor := gameblocks.NewFloor(ctx.Shaders.Get("main"))
	w.Add(floor)
	//emitter := gameblocks.ParticleEmitter
	return nil
}

func (w *world) Tick(d time.Duration) error {
	return nil
}

func (w *world) Reset() {}
func (w *world) Focus() gameblocks.Vector {
	return gameblocks.FixedVector(mgl.Vec3{0, 0, 0}, mgl.Vec3{0, 0, 1})
}
