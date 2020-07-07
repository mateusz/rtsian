package main

import (
	"sort"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type mobile interface {
	Input(win *pixelgl.Window, cam pixel.Matrix)
	Update(dt float64)
	Draw(pixel.Target)
	GetZ() float64
	GetX() float64
	GetY() float64
}

type mobiles struct {
	List []mobile
}

func (m *mobiles) ByZ() []mobile {
	mobs := m.List
	sort.SliceStable(mobs, func(i, j int) bool {
		return mobs[i].GetZ() < mobs[j].GetZ()
	})
	return mobs
}

func (m *mobiles) ByReverseZ() []mobile {
	mobs := m.List
	sort.SliceStable(mobs, func(i, j int) bool {
		return mobs[i].GetZ() > mobs[j].GetZ()
	})
	return mobs
}

func (m *mobiles) Add(mob mobile) {
	m.List = append(m.List, mob)
}

func (m *mobiles) Remove(mob mobile) {
	for i, ml := range m.List {
		if mob == ml {
			copy(m.List[i:], m.List[i+1:])
			m.List = m.List[:len(m.List)-1]
			break
		}
	}
}

func (m *mobiles) Len() int {
	return len(m.List)
}
