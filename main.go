package main

import (
	"fmt"
	"github.com/just1689/pyrate-server/maps"
)

func main() {

	var chunk maps.Chunk
	chunk = make([]*maps.Tile, 50*50)

	//l := len(chunk)

	n := 0
	for x := 0; x < 50; x++ {
		for y := 0; y < 50; y++ {

			t := maps.Tile{}
			t.X = x
			t.Y = y
			t.ID = fmt.Sprint(t.X, ".", t.Y)
			chunk[n] = &t
			n++
		}
	}

	chunk.CoverInWater()

	maps.GenerateChunk(chunk, maps.TwoIslands)

}
