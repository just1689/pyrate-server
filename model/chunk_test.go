package model

import (
	"testing"
)

func TestGenerateWaterChunks(t *testing.T) {

	width := 100
	chunkWidth := 10
	x1 := 0
	x2 := width - 1
	z1 := 0
	z2 := width - 1

	c := GenerateWaterChunks(x1, x2, z1, z2, chunkWidth)
	var chunks []Chunk
	for chunk := range c {
		chunks = append(chunks, chunk)
	}

	if len(chunks) != 100 {
		t.Error("Unexpected number of chunks! Not being ", 100, " and not ", len(chunks))
	}

}

func TestGenerateWaterChunk(t *testing.T) {

	width := 100
	x1 := 0
	x2 := width - 1
	z1 := 0
	z2 := width - 1

	chunk := GenerateWaterChunk(x1, x2, z1, z2)

	if len(chunk) != width*width {
		t.Error("Chunk not the right size")
	}

	last := chunk[width*width-1]
	if last.X != width-1 {
		t.Error("Last tile unexpected X")
	}
	if last.Z != width-1 {
		t.Error("Last tile unexpected Z")
	}

}

func TestChunk_CoverInWater(t *testing.T) {
	var c Chunk
	c = make([]*Tile, 4)
	c[0] = &Tile{X: 10}
	c[1] = &Tile{X: 11}
	c[2] = &Tile{X: 12}
	c[3] = &Tile{X: 13}

	c.CoverInWater()

	for _, tile := range c {
		if tile.TileType != TileTypeWater {
			t.Error("Unexpected tile type - ", tile.TileType, "Not", TileTypeWater)
		}
	}

}

func TestChunk_GetXMid(t *testing.T) {

	var c Chunk
	c = make([]*Tile, 4)
	c[0] = &Tile{X: 10}
	c[1] = &Tile{X: 11}
	c[2] = &Tile{X: 12}
	c[3] = &Tile{X: 13}

	correct := 11
	attempt := c.GetXMid()

	if attempt != correct {
		t.Error("Mid should be ", correct, " not ", attempt, "Min was:", c.GetXMin(), "Max was", c.GetXMax())
	}

}

func TestChunk_GetXMin(t *testing.T) {

	m := 0

	var c Chunk
	c = make([]*Tile, 3)
	c[0] = &Tile{X: 1}
	c[1] = &Tile{X: m}
	c[2] = &Tile{X: 3}

	end := c.GetXMin()

	if end != m {
		t.Error("Max should be ", m, " not ", end)
	}

}

func TestChunk_GetXMax(t *testing.T) {

	m := 20

	var c Chunk
	c = make([]*Tile, 3)
	c[0] = &Tile{X: 1}
	c[1] = &Tile{X: 20}
	c[2] = &Tile{X: 3}

	end := c.GetXMax()

	if end != m {
		t.Error("Max should be ", m)
	}

}
