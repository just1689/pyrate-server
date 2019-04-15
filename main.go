package main

import (
	"github.com/just1689/pyrate-server/maps"
)

func main() {

	var chunk maps.Chunk
	chunk = make([]*maps.Tile, 50*50)

	x := 0
	y := 0
	l := len(chunk)

	for i := 0; i < l; i++ {
		t := maps.Tile{}
		t.X = x
		t.Y = y
		x++
		if x >= 49 {
			x = 0
			y++
		}
		chunk[i] = &t
	}

	maps.GenerateChunk(chunk, maps.SingleIsland)

}
