package maps

import "math/rand"

func RandomizeChunckType() ChunkType {

	r := rand.Intn(100)

	if r >= 0 && r < 40 {
		// 40% chance
		return Nothing
	} else if r >= 40 && r < 60 {
		// 20% chance
		return SingleIsland
	} else if r >= 60 && r < 80 {
		// 20% chance
		return TwoIslands
	} else if r >= 80 && r < 100 {
		// 20% chance
		return SmallIslands
	}
	return Nothing

}

func GenerateChunk(chunk Chunk, chunkType ChunkType) {
	if chunkType == Nothing {
		generateNothing(chunk)
	} else if chunkType == SingleIsland {
		generateSingleIsland(chunk)
	} else if chunkType == TwoIslands {
		generateTwoIslands(chunk)
	} else if chunkType == SmallIslands {
		generateSmallIslands(chunk)
	}
}
