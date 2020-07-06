package main

import (
	"encoding/xml"
	"errors"
	"image"
	"os"
	"path/filepath"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/lafriks/go-tiled"
)

type spriteset struct {
	sprites  map[uint32]*pixel.Sprite
	tileset  *tiled.Tileset
	basePath string
}

func newSpriteset() spriteset {
	t := spriteset{}
	t.sprites = make(map[uint32]*pixel.Sprite)
	return t
}

// TMX library does not load images. Help it out.
func fillMissingMapPieces(m *tiled.Map) (spriteset, error) {
	spr := spriteset{}
	var err error
	for _, ts := range m.Tilesets {
		if ts.Source == "" {
			return spr, errors.New("Tileset has no source")
		}
		spr, err = newSpritesetFromTileset(m.GetFileFullPath(""), ts)
		if err != nil {
			return spr, err
		}

		// Only one permitted at the moment.
		break
	}

	return spr, nil
}

func newSpritesetFromTsx(basePath, path string) (spriteset, error) {
	spr := spriteset{}
	ts := &tiled.Tileset{Source: path}

	f, err := os.Open(filepath.Join(basePath, ts.Source))
	if err != nil {
		return spr, err
	}
	defer f.Close()

	d := xml.NewDecoder(f)
	if err := d.Decode(ts); err != nil {
		return spr, err
	}

	spr, err = newSpritesetFromTileset(basePath, ts)
	if err != nil {
		return spr, err
	}

	spr.tileset = ts
	return spr, nil
}

// Load actual sprite files and associate with tileset.
func newSpritesetFromTileset(basePath string, ts *tiled.Tileset) (spriteset, error) {
	spr := newSpriteset()
	spr.tileset = ts
	spr.basePath = basePath

	f, err := os.Open(filepath.Join(basePath, ts.Source))

	if err != nil {
		return spr, err
	}
	defer f.Close()

	d := xml.NewDecoder(f)

	if err := d.Decode(ts); err != nil {
		return spr, err
	}

	for _, t := range ts.Tiles {
		if t.Image.Source == "" {
			continue
		}

		file, err := os.Open(filepath.Join(basePath, t.Image.Source))
		if err != nil {
			return spr, err
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			return spr, err
		}

		pic := pixel.PictureDataFromImage(img)
		spr.sprites[t.ID] = pixel.NewSprite(pic, pic.Bounds())
	}

	return spr, nil
}
