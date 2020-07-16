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
	patherMap = make([]*patherNode, len(gameWorld.Tiles.Layers[0].Tiles))
}

func GetPatherNode(x, y int) *patherNode {
	return patherMap[y*gameWorld.Tiles.Width+x]
}

func GetPatherNodeFromVec(v pixel.Vec) *patherNode {
	x, y := gameWorld.VecToTile(v)
	return GetPatherNode(x, y)
}

func PatherBuildState(collisionArea pixel.Rect) {
	for y := 0; y < gameWorld.Tiles.Height; y++ {
		for x := 0; x < gameWorld.Tiles.Width; x++ {
			c, err := strconv.ParseFloat(gameWorld.TilesetTileAt(x, y).Properties.GetString("movd"), 64)
			if err != nil {
				c = 10.0
			}
			patherMap[y*gameWorld.Tiles.Width+x] = &patherNode{
				X:    x,
				Y:    y,
				Cost: c,
			}
		}
	}
	for _, m := range gamePositionables.List {
		mvec := pixel.Vec{X: m.GetX(), Y: m.GetY()}
		if !collisionArea.Contains(mvec) {
			continue
		}
		GetPatherNodeFromVec(mvec).Cost = 1000000000.0
		u, ok := m.(*unit)
		if !ok {
			continue
		}
		// Forbid movement-target tile too
		GetPatherNodeFromVec(u.target).Cost = 1000000000.0
	}
}

func FindPath(p positionable, target pixel.Vec) (l *list.List) {
	apos := gameWorld.AlignToTile(pixel.Vec{X: p.GetX(), Y: p.GetY()})
	collisionRange := 24.0
	PatherBuildState(pixel.Rect{
		Min: pixel.Vec{X: apos.X - collisionRange, Y: apos.Y - collisionRange},
		Max: pixel.Vec{X: apos.X + collisionRange, Y: apos.Y + collisionRange},
	})
	l = list.New()
	path, _, found := astar.Path(
		GetPatherNodeFromVec(pixel.Vec{X: p.GetX(), Y: p.GetY()}),
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
	if n.X < gameWorld.Tiles.Width-1 {
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
	if n.Y < gameWorld.Tiles.Height-1 {
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
