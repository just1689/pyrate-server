package main

import (
	"fmt"
	"github.com/just1689/pyrate-server/maps"
)

func main() {

	var chunk maps.Chunk
	chunk = make([]*maps.Tile, 1000*1000)

	//l := len(chunk)

	n := 0
	for x := 0; x < 1000; x++ {
		for y := 0; y < 1000; y++ {

			t := maps.Tile{}
			t.X = x
			t.Y = y
			t.ID = fmt.Sprint(t.X, ".", t.Y)
			chunk[n] = &t
			n++
		}
	}

	chunk.CoverInWater()

	//maps.GenerateChunk(chunk, maps.SingleIsland)

	conn, err := maps.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}

	for n, row := range chunk {
		if (n == 2000) || (n == 4000) || (n == 6000) || (n == 8000) {
			fmt.Println("At", n)
		}
		maps.InsertTile(conn, row)

	}

}
