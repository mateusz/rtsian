package main

import (
	"sort"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// Any mixture of below interfaces.
type entity interface {
}

type updateable interface {
	Update(float64)
}

type inputtable interface {
	Input(*pixelgl.Window, pixel.Matrix)
}

type drawable interface {
	Draw(pixel.Target)
}

type positionable interface {
	GetZ() float64
	GetX() float64
	GetY() float64
}

type mouseHittable interface {
	positionable
	MouseHit(*pixelgl.Window)
	MouseClear()
}

type explodable interface {
	startExploding()
}

type entities struct {
	List []entity
}

func (e *entities) Add(ent entity) {
	e.List = append(e.List, ent)
}

func (e *entities) Remove(ent entity) {
	for i, ml := range e.List {
		if ent == ml {
			copy(e.List[i:], e.List[i+1:])
			e.List = e.List[:len(e.List)-1]
			break
		}
	}
}

func (e *entities) Len() int {
	return len(e.List)
}

type drawables struct {
	List []drawable
}

func (p *drawables) ByDrawOrder() []drawable {
	ents := p.List
	sort.SliceStable(ents, func(i, j int) bool {
		pi, oki := ents[i].(positionable)
		pj, okj := ents[j].(positionable)

		orderi := 0.0
		orderj := 0.0
		if oki {
			orderi = pi.GetZ()
		}
		if okj {
			orderj = pj.GetZ()
		}

		return orderi < orderj
	})
	return ents
}

func (d *drawables) Add(ent drawable) {
	d.List = append(d.List, ent)
}

func (d *drawables) Remove(ent drawable) {
	for i, ml := range d.List {
		if ent == ml {
			copy(d.List[i:], d.List[i+1:])
			d.List = d.List[:len(d.List)-1]
			break
		}
	}
}

func (d *drawables) Len() int {
	return len(d.List)
}

type positionables struct {
	List []positionable
}

func (p *positionables) ByZ() []positionable {
	ents := p.List
	sort.SliceStable(ents, func(i, j int) bool {
		return ents[i].GetZ() < ents[j].GetZ()
	})
	return ents
}

func (p *positionables) ByReverseZ() []positionable {
	ents := p.List
	sort.SliceStable(ents, func(i, j int) bool {
		return ents[i].GetZ() > ents[j].GetZ()
	})
	return ents
}

func (p *positionables) Add(ent positionable) {
	p.List = append(p.List, ent)
}

func (p *positionables) Remove(ent positionable) {
	for i, ml := range p.List {
		if ent == ml {
			copy(p.List[i:], p.List[i+1:])
			p.List = p.List[:len(p.List)-1]
			break
		}
	}
}

func (p *positionables) Len() int {
	return len(p.List)
}
