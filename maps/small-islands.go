package maps

import (
	"math/rand"
	"time"
)

const smallIslandPercentage = 3
const smallIslandCount = 3

func generateSmallIslands(chunk Chunk) {

	for i := 0; i < smallIslandCount; i++ {
		generateSmallIsland(chunk)
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

}

func generateSmallIsland(chunk Chunk) {
	//Pick a starting point
	randSource := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(randSource)
	randX, randY := randXAndY(chunk, rnd)
	//fmt.Println("For chunk starting a ", chunk[0].X, chunk[0].Y, " point for single island will be ", randX, randY)

	//Start with just water
	chunk.CoverInWater()
	//Number of tiles to become land
	size := len(chunk) * smallIslandPercentage / 100
	var t *Tile
	var ok bool
	for size > 0 {
		ok, t = chunk.FindFirstWater(randX, randY, rnd)
		if !ok {
			randX, randY = randXAndY(chunk, rnd)
			continue
		}
		t.TileType = TileTypeLand
		size--
	}

}
