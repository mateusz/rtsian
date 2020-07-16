package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/lafriks/go-tiled"
	"github.com/mateusz/rtsian/piksele"
	"golang.org/x/image/colornames"
)

var (
	workDir           string
	monW              float64
	monH              float64
	pixSize           float64
	mobSprites        piksele.Spriteset
	cursorSprites     piksele.Spriteset
	p1                player
	gameWorld         piksele.World
	gameHud           hud
	gameEntities      entities
	gameDrawables     drawables
	gamePositionables positionables
	gameMouseHits     []mouseHittable
)

func main() {
	rand.Seed(time.Now().UnixNano())

	var err error
	workDir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Printf("Error checking working dir: %s\n", err)
		os.Exit(2)
	}

	gameWorld = piksele.World{}
	gameWorld.Load(fmt.Sprintf("%s/../assets/world.tmx", workDir))

	mobSprites, err = piksele.NewSpritesetFromTsx(fmt.Sprintf("%s/../assets", workDir), "mobs.tsx")
	if err != nil {
		fmt.Printf("Error loading mobs: %s\n", err)
		os.Exit(2)
	}

	loadObjects(gameWorld.Tiles)

	cursorSprites, err = piksele.NewSpritesetFromTsx(fmt.Sprintf("%s/../assets", workDir), "cursors.tsx")
	if err != nil {
		fmt.Printf("Error loading cursors: %s\n", err)
		os.Exit(2)
	}

	p1.position = pixel.Vec{
		X: float64(gameWorld.PixelWidth()) / 2.0,
		Y: float64(gameWorld.PixelHeight()) / 2.0,
	}
	p1.scrollSpeed = 200.0
	p1.scrollHotZone = 10.0
	p1.army = 1

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

	mapCanvas := pixelgl.NewCanvas(pixel.R(0, 0, float64(gameWorld.PixelWidth()), float64(gameWorld.PixelHeight())))
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

		// Register mouse hits
		gameMouseHits = gameMouseHits[:0]
		mp := cam1.Unproject(win.MousePosition().Scaled(1.0 / pixSize))
		for _, ent := range gameEntities.List {
			hit, ok := ent.(mouseHittable)
			if ok {
				if mp.X > hit.GetX()-8 &&
					mp.X < hit.GetX()+8 &&
					mp.Y > hit.GetY()-8 &&
					mp.Y < hit.GetY()+8 {
					hit.MouseHit(win)
					gameMouseHits = append(gameMouseHits, hit)
				} else {
					hit.MouseClear()
				}
			}
		}

		// Update world state
		UnitInput(win, cam1)
		p1.Input(win, cam1)
		p1.Update(dt)
		gameHud.Update(dt)

		for _, ent := range gameEntities.List {
			inp, ok := ent.(inputtable)
			if ok {
				inp.Input(win, cam1)
			}

			upd, ok := ent.(updateable)
			if ok {
				upd.Update(dt)
			}
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
		for _, mob := range gameDrawables.ByDrawOrder() {
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

func loadObjects(m *tiled.Map) {
	for _, o := range m.ObjectGroups[0].Objects {
		lt, err := m.TileGIDToTile(o.GID)
		if err != nil {
			log.Fatal(err)
		}
		army, err := strconv.Atoi(o.Properties.GetString("army"))
		if err != nil {
			army = 0
		}
		if lt.ID >= MOBS_TANK_START_ID && lt.ID < MOBS_TANK_START_ID+4 {
			p := gameWorld.AlignToTile(pixel.Vec{X: o.X + 10.0, Y: piksele.TiledFlipY(m, o.Y) + 10.0})
			u := NewUnit(p, army)
			u.target = u.position
			gameEntities.Add(&u)
			gamePositionables.Add(&u)
			gameDrawables.Add(&u)
		}
	}
}
