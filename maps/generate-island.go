package maps

import (
	"github.com/just1689/pyrate-server/model"
	"math/rand"
)

func RandomizeChunckType() model.ChunkType {

	r := rand.Intn(100)

	if r >= 0 && r < 40 {
		// 40% chance
		return model.Nothing
	} else if r >= 40 && r < 60 {
		// 20% chance
		return model.SingleIsland
	} else if r >= 60 && r < 80 {
		// 20% chance
		return model.TwoIslands
	} else if r >= 80 && r < 100 {
		// 20% chance
		return model.SmallIslands
	}
	return model.Nothing

}

func GenerateChunk(chunk model.Chunk, chunkType model.ChunkType) {
	if chunkType == model.Nothing {
		generateNothing(chunk)
	} else if chunkType == model.SingleIsland {
		generateSingleIsland(chunk)
	} else if chunkType == model.TwoIslands {
		generateTwoIslands(chunk)
	} else if chunkType == model.SmallIslands {
		generateSmallIslands(chunk)
	}
}
