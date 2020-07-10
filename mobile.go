package main

import (
	"container/list"
	"log"
	"math"

	"github.com/faiface/pixel"
)

type mobile struct {
	position        pixel.Vec
	target          pixel.Vec // position of current movement target
	d               float64   // distance to current movement target
	v               pixel.Vec // velocity vector
	pathingTarget   pixel.Vec
	pathing         *list.List
	stickyDirOffset uint32
	baseSpeed       float64
}

func (m *mobile) Update(dt float64) {
	if m.d > 0.0 {
		m.position = m.position.Add(m.v)
		m.d -= m.v.Len()
	} else {
		m.applyPath()
	}

	m.updateDirOffset()
}

func (m *mobile) applyPath() {
	if m.pathing == nil || m.pathing.Len() == 0 {
		m.target = m.position
		m.d = 0.0
		m.v = pixel.ZV
		m.pathing = nil
		return
	}

	m.pathing = FindPath(m, m.pathingTarget)
	if m.pathing.Len() == 0 {
		return
	}

	// Next path step
	n, ok := m.pathing.Remove(m.pathing.Front()).(*patherNode)
	if !ok {
		log.Panic("Fatal: pathing list contained non-patherNode!")
	}

	m.target = gameWorld.tileToVec(n.X, n.Y)
	mv := m.target.Sub(m.position)
	m.d = mv.Len()
	m.v = mv.Unit().Scaled(m.baseSpeed / n.Cost)
}

func (m *mobile) updateDirOffset() {
	if math.Abs(m.v.X) > math.Abs(m.v.Y) {
		if m.v.X < 0 {
			m.stickyDirOffset = 0
		}
		if m.v.X > 0 {
			m.stickyDirOffset = 2
		}
	} else {
		if m.v.Y > 0 {
			m.stickyDirOffset = 1
		}
		if m.v.Y < 0 {
			m.stickyDirOffset = 3
		}
	}
}

func (m *mobile) GetZ() float64 {
	return -m.position.Y
}

func (m *mobile) GetX() float64 {
	return m.position.X
}

func (m *mobile) GetY() float64 {
	return m.position.Y
}
