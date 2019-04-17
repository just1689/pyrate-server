package maps

import (
	"github.com/just1689/pyrate-server/model"
	"math/rand"
	"time"
)

const largeIslandPercentage = 15

func generateSingleIsland(chunk model.Chunk) {

	//Pick a starting point
	randSource := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(randSource)
	randX, randY := chunk.RandXAndY(rnd)
	//fmt.Println("For chunk starting a ", chunk[0].X, chunk[0].Y, " point for single island will be ", randX, randY)

	//Start with just water
	chunk.CoverInWater()
	//Number of tiles to become land
	size := len(chunk) * largeIslandPercentage / 100
	var t *model.Tile
	var ok bool
	for size > 0 {
		ok, t = chunk.FindFirstWater(randX, randY, rnd)
		if !ok {
			randX, randY = chunk.RandXAndY(rnd)
			continue
		}
		t.TileType = model.TileTypeLand
		size--
	}

	//count := 1
	//for _, t := range chunk {
	//	if count == 51 {
	//		fmt.Print("\n")
	//		count = 1
	//	}
	//	if t.TileType == TileTypeWater {
	//		fmt.Print("W")
	//	} else if t.TileType == TileTypeLand {
	//		fmt.Print("L")
	//	} else {
	//		fmt.Print(" ")
	//	}
	//	count++
	//}

	//fmt.Println(len(chunk))

}
