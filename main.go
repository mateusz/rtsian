package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/lafriks/go-tiled"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

var (
	workDir    string
	terra      spriteset
	mobSprites spriteset
	mobs       []mobile
	tmx        *tiled.Map
	p1         player
	monW float64
	monH float64
)

func main() {
	rand.Seed(time.Now().UnixNano())

	var err error
	workDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Printf("Error checking working dir: %s\n", err)
		os.Exit(2)
	}

	tmx, err = tiled.LoadFromFile(fmt.Sprintf("%s/../assets/map1.tmx", workDir))
	if err != nil {
		fmt.Printf("Error parsing map: %s\n", err)
		os.Exit(2)
	}
	terra, err = fillMissingMapPieces(tmx)
	if err != nil {
		fmt.Printf("Error loading aux tilesets: %s\n", err)
		os.Exit(2)
	}

	mobSprites, err = newSpritesetFromTsx(fmt.Sprintf("%s/../assets", workDir), "mobs.tsx")
	if err != nil {
		fmt.Printf("Error loading mobs: %s\n", err)
		os.Exit(2)
	}

	p1.wp = pixel.Vec{
		X: float64(tmx.Width*tmx.TileWidth) / 2.0,
		Y: float64(tmx.Height*tmx.TileHeight) / 2.0,
	}

	scrollSpeed = 200.0
	scrollHotZone = 10.0

	pixelgl.Run(run)
}

func run() {
	monitor := pixelgl.PrimaryMonitor()

	monW, monH = monitor.Size()
	pixSize := 4.0

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
	win.SetMousePosition(pixel.Vec{X:monW/2.0,Y:monH/2.0})

	worldMap := pixelgl.NewCanvas(pixel.R(0, 0, float64(tmx.Width*tmx.TileWidth), float64(tmx.Height*tmx.TileHeight)))
	drawMap(worldMap)

	p1view := pixelgl.NewCanvas(pixel.R(0, 0, monW/pixSize, monH/pixSize))
	hud := pixelgl.NewCanvas(pixel.R(0, 0, monW/pixSize, monH/pixSize))

	staticHud := imdraw.New(nil)
	staticHud.Color = colornames.Black
	fps := text.New(pixel.ZV, text.NewAtlas(basicfont.Face7x13, text.ASCII))

	last := time.Now()
	fpsAvg := 60.0
	for !win.Closed() {
		if win.Pressed(pixelgl.KeyEscape) {
			break
		}

		dt := time.Since(last).Seconds()
		last = time.Now()

		fpsAvg -= fpsAvg / 50.0
		fpsAvg += 1.0 / dt / 50.0

		p1.Update(dt, win)

		// Center views on players
		cam1 := pixel.IM.Moved(pixel.Vec{
			X: -p1.wp.X + p1view.Bounds().W()/2,
			Y: -p1.wp.Y + p1view.Bounds().H()/2,
		})
		p1view.SetMatrix(cam1)

		// Draw
		win.Clear(colornames.Black)
		hud.Clear(pixel.RGBA{})
		p1view.Clear(colornames.Green)

		worldMap.Draw(p1view, pixel.IM.Moved(pixel.Vec{
			X: worldMap.Bounds().W() / 2.0,
			Y: worldMap.Bounds().H() / 2.0,
		}))

		if win.Pressed(pixelgl.KeyG) {
			fps.Clear()
			fmt.Fprintf(fps, "%.0f", fpsAvg)
			fps.Draw(hud, pixel.IM)
		}

		sort.Slice(mobs, func(i, j int) bool {
			return mobs[i].GetZ() > mobs[j].GetZ()
		})
		for _, mob := range mobs {
			mob.Update(dt)
			mob.Draw(p1view)
		}

		// Draw  views onto respective halves of the screen
		p1view.Draw(win, pixel.IM.Moved(pixel.Vec{
			X: p1view.Bounds().W() / 2,
			Y: p1view.Bounds().H() / 2,
		}))

		staticHud.Draw(hud)
		hud.Draw(win, pixel.IM.Moved(pixel.V(hud.Bounds().W()/2, hud.Bounds().H()/2)))
		win.Update()
	}
}

func drawMap(c *pixelgl.Canvas) {
	l := tmx.Layers[0]
	for y := 0; y < tmx.Height; y++ {
		for x := 0; x < tmx.Width; x++ {
			lt := l.Tiles[y*tmx.Width+x]
			terra.sprites[lt.ID].Draw(c, pixel.IM.Moved(tileVec(x, tmx.Height-y-1)))
		}
	}
}
