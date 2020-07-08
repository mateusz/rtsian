package main

import (
	"container/list"
	"math"
	"strconv"

	"github.com/beefsack/go-astar"
	"github.com/faiface/pixel"
)

var (
	patherMap []*patherNode
)

type patherNode struct {
	X    int
	Y    int
	Cost float64
}

func PatherInit() {
	patherMap = make([]*patherNode, len(gameWorld.tiles.Layers[0].Tiles))
}

func GetPatherNode(x, y int) *patherNode {
	return patherMap[y*gameWorld.tiles.Width+x]
}

func GetPatherNodeFromVec(v pixel.Vec) *patherNode {
	x, y := gameWorld.vecToTile(v)
	return GetPatherNode(x, y)
}

func PatherBuildState() {
	for y := 0; y < gameWorld.tiles.Height; y++ {
		for x := 0; x < gameWorld.tiles.Width; x++ {
			c, err := strconv.ParseFloat(gameWorld.tilesetTileAt(x, y).Properties.GetString("movd"), 64)
			if err != nil {
				c = 1.0
			}
			patherMap[y*gameWorld.tiles.Width+x] = &patherNode{
				X:    x,
				Y:    y,
				Cost: c,
			}
		}
	}
	for _, m := range gameMobiles.List {
		GetPatherNodeFromVec(pixel.Vec{X: m.GetX(), Y: m.GetY()}).Cost = 1000000000.0
		u, ok := m.(*unit)
		if !ok {
			continue
		}
		// Forbid movement-target tile too
		GetPatherNodeFromVec(u.target).Cost = 1000000000.0
	}
}

func FindPath(m mobile, target pixel.Vec) (l *list.List) {
	l = list.New()
	path, _, found := astar.Path(
		GetPatherNodeFromVec(pixel.Vec{X: m.GetX(), Y: m.GetY()}),
		GetPatherNodeFromVec(target),
	)
	if !found {
		return
	}
	for _, n := range path {
		l.PushFront(n)
	}
	// Remove starting tile
	l.Remove(l.Front())

	return
}

func (n *patherNode) PathNeighbors() []astar.Pather {
	ns := []astar.Pather{}
	if n.X > 0 {
		pn := GetPatherNode(n.X-1, n.Y)
		if pn.Cost < 1000.0 {
			ns = append(ns, pn)
		}
	}
	if n.X < gameWorld.tiles.Width-1 {
		pn := GetPatherNode(n.X+1, n.Y)
		if pn.Cost < 1000.0 {
			ns = append(ns, pn)
		}
	}
	if n.Y > 0 {
		pn := GetPatherNode(n.X, n.Y-1)
		if pn.Cost < 1000.0 {
			ns = append(ns, pn)
		}
	}
	if n.Y < gameWorld.tiles.Height-1 {
		pn := GetPatherNode(n.X, n.Y+1)
		if pn.Cost < 1000.0 {
			ns = append(ns, pn)
		}
	}

	return ns
}

func (n *patherNode) PathNeighborCost(to astar.Pather) float64 {
	tn, ok := to.(*patherNode)
	if !ok {
		return 10000000.0
	}

	return tn.Cost
}

func (n *patherNode) PathEstimatedCost(to astar.Pather) float64 {
	tn, ok := to.(*patherNode)
	if !ok {
		return 10000000.0
	}

	return math.Abs(float64(tn.X-n.X)) + math.Abs(float64(tn.Y-n.Y))
}
