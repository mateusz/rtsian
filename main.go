package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var (
	workDir       string
	monW          float64
	monH          float64
	pixSize       float64
	mobSprites    spriteset
	cursorSprites spriteset
	p1            player
	gameWorld     world
	gameHud       hud
	gameMobiles   mobiles
)

func main() {
	rand.Seed(time.Now().UnixNano())

	var err error
	workDir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Printf("Error checking working dir: %s\n", err)
		os.Exit(2)
	}

	gameWorld = world{}
	gameWorld.load()

	mobSprites, err = newSpritesetFromTsx(fmt.Sprintf("%s/../assets", workDir), "mobs.tsx")
	if err != nil {
		fmt.Printf("Error loading mobs: %s\n", err)
		os.Exit(2)
	}

	cursorSprites, err = newSpritesetFromTsx(fmt.Sprintf("%s/../assets", workDir), "cursors.tsx")
	if err != nil {
		fmt.Printf("Error loading cursors: %s\n", err)
		os.Exit(2)
	}

	p1.position = pixel.Vec{
		X: float64(gameWorld.pixelWidth()) / 2.0,
		Y: float64(gameWorld.pixelHeight()) / 2.0,
	}
	p1.scrollSpeed = 200.0
	p1.scrollHotZone = 10.0

	PatherInit()

	pixelgl.Run(run)
}

func run() {
	monitor := pixelgl.PrimaryMonitor()

	monW, monH = monitor.Size()
	pixSize = 4.0

	cfg := pixelgl.WindowConfig{
		Title:   "Rtsian",
		Bounds:  pixel.R(0, 0, monW, monH),
		VSync:   true,
		Monitor: monitor,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Zoom in to get nice pixels
	win.SetSmooth(false)
	win.SetMatrix(pixel.IM.Scaled(pixel.ZV, pixSize))
	win.SetMousePosition(pixel.Vec{X: monW / 2.0, Y: monH / 2.0})

	mapCanvas := pixelgl.NewCanvas(pixel.R(0, 0, float64(gameWorld.pixelWidth()), float64(gameWorld.pixelHeight())))
	gameWorld.Draw(mapCanvas)

	p1view := pixelgl.NewCanvas(pixel.R(0, 0, monW/pixSize, monH/pixSize))
	hudCanvas := pixelgl.NewCanvas(pixel.R(0, 0, monW/pixSize, monH/pixSize))
	gameHud.bounds = p1view.Bounds()

	last := time.Now()
	for !win.Closed() {
		if win.Pressed(pixelgl.KeyEscape) {
			break
		}

		dt := time.Since(last).Seconds()
		last = time.Now()

		// Move player's view
		cam1 := pixel.IM.Moved(pixel.Vec{
			X: -p1.position.X + p1view.Bounds().W()/2,
			Y: -p1.position.Y + p1view.Bounds().H()/2,
		})
		p1view.SetMatrix(cam1)

		// Update world state
		PatherBuildState()
		UnitInput(win, cam1)
		p1.Input(win, cam1)
		p1.Update(dt)
		gameHud.Update(dt)
		for _, mob := range gameMobiles.List {
			mob.Input(win, cam1)
			mob.Update(dt)
		}

		// Clean up for new frame
		win.Clear(colornames.Black)
		hudCanvas.Clear(pixel.RGBA{})
		p1view.Clear(colornames.Green)

		// Draw transformed map
		mapCanvas.Draw(p1view, pixel.IM.Moved(pixel.Vec{
			X: mapCanvas.Bounds().W() / 2.0,
			Y: mapCanvas.Bounds().H() / 2.0,
		}))

		// Draw transformed mobs
		for _, mob := range gameMobiles.ByZ() {
			mob.Draw(p1view)
		}

		// Render hud
		gameHud.Draw(hudCanvas)

		// Blit player view
		p1view.Draw(win, pixel.IM.Moved(pixel.Vec{
			X: p1view.Bounds().W() / 2,
			Y: p1view.Bounds().H() / 2,
		}))

		// Overlay with hud
		hudCanvas.Draw(win, pixel.IM.Moved(pixel.V(hudCanvas.Bounds().W()/2, hudCanvas.Bounds().H()/2)))

		// Present frame!
		win.Update()
	}
}
