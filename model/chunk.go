package model

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
	"log"
	"math/rand"
	"time"
)

type ChunkType string

const Nothing ChunkType = "Nothing"
const SingleIsland ChunkType = "SingleIsland"
const TwoIslands ChunkType = "TwoIslands"
const SmallIslands ChunkType = "SmallIslands"

const ChunkSize = 50
const MapWidth = 1000

type Chunk []*Tile

func GenerateWaterChunks(x1, x2, z1, z2, chunkSide int) (c chan Chunk) {
	c = make(chan Chunk, 2) // Cache up to N chunks
	go func() {
		for x := x1; x < x2; x += chunkSide {
			for z := z1; z < z2; z += chunkSide {
				c <- GenerateWaterChunk(x, x+chunkSide-1, z, z+chunkSide-1)
			}
		}
		close(c)
	}()
	return
}

func GenerateWaterChunk(x1, x2, z1, z2 int) (c Chunk) {
	for x := x1; x <= x2; x++ {
		for z := z1; z <= z2; z++ {
			c = append(c, &Tile{
				ID:       fmt.Sprint(x, ".", z),
				X:        x,
				Z:        z,
				TileType: TileTypeWater,
				TileSkin: "",
			})
		}
	}
	return
}

func (c Chunk) ToInterface() (result [][]interface{}) {
	for _, t := range c {
		result = append(result, t.ToInterface())
	}
	return
}

func (c Chunk) Size() int {
	return len(c)
}

func (c Chunk) GetXMin() (r int) {
	r = 99999999
	for _, t := range c {
		if t.X < r {
			r = t.X
		}
	}
	return
}

func (c Chunk) GetXMax() (r int) {
	r = 0
	for _, t := range c {
		if t.X > r {
			r = t.X
		}
	}
	return
}

func (c Chunk) GetXMid() (r int) {
	x1 := c.GetXMin()
	x2 := c.GetXMax()
	m := x2 - x1
	r = x1 + m/2
	return
}

func (c Chunk) GetYMin() (r int) {
	r = 99999999
	for _, t := range c {
		if t.Z < r {
			r = t.Z
		}
	}
	return
}

func (c Chunk) GetYMax() (r int) {

	//Possibly the least efficient impl possible
	r = 0
	for _, t := range c {
		if t.Z > r {
			r = t.Z
		}
	}
	return
}

func (c Chunk) GetYMid() (r int) {
	x1 := c.GetYMin()
	x2 := c.GetYMax()
	m := x2 - x1
	r = x1 + m/2
	return
}

func (c Chunk) CoverInWater() {
	for _, t := range c {
		t.TileType = TileTypeWater
	}
}

func (c Chunk) FindFirstWater(x, y int, r *rand.Rand) (ok bool, t *Tile) {

	for {
		var found bool
		found, t = c.GetAt(x, y)
		if !found {
			log.Fatalln("No tile(1):", x, y)
		}
		if t.TileType == TileTypeWater {
			//fmt.Println("Found one:", t.X, t.Y)
			ok = true
			return
		}

		start := time.Now()
		for {
			d := time.Since(start)
			if d.Seconds() > 1 {
				ok = false
				fmt.Println("Took too long..")
				return
			}
			//Pick a dir randomly
			dir := r.Intn(100)
			if dir < 25 {
				x--
				if x <= c.GetXMin() {
					x++
					//fmt.Println("Bumped <x")
					continue
				}
			} else if dir < 50 {
				x++
				if x >= c.GetXMax() {
					x--
					//fmt.Println("Bumped >x")
					continue
				}
			} else if dir < 75 {
				y++
				if y >= c.GetYMax() {
					y--
					//fmt.Println("Bumped >y")
					continue
				}
			} else {
				y--
				if y <= c.GetYMax() {
					y++
					//fmt.Println("Bumped <y")
					continue
				}
			}

			//fmt.Print("Checking:", x, y)
			found, t = c.GetAt(x, y)
			//fmt.Println(found)
			if !found {
				log.Fatalln("No tile(2):", x, y)
			}
			if t.TileType == TileTypeWater {
				//fmt.Println("Found one(1):", t.X, t.Y)
				ok = true
				return
			}
			if t.TileType == "" {
				//fmt.Println("Found one(2):", t.X, t.Y)
				ok = true
				return
			}

		}

	}

}

func (c Chunk) GetAt(x, y int) (found bool, tile *Tile) {
	for _, tile = range c {
		//fmt.Println(x, tile.X, y, tile.Y)
		if tile.X == x && tile.Z == y {
			found = true
			return
		}
	}
	found = false
	return
}

func (c Chunk) GetAsInterfaceSlSl() (result [][]interface{}) {
	result = make([][]interface{}, 0)
	for _, tile := range c {
		row := make([]interface{}, 0)
		row = append(row, tile.ID, tile.X, tile.Z, tile.TileType, tile.TileSkin)
		result = append(result, row)
	}
	return
}

func (c Chunk) RandXAndY(rnd *rand.Rand) (randX, randY int) {
	randX = c.GetXMid() + rnd.Intn(10) - 5
	randY = c.GetYMid() + rnd.Intn(10) - 5
	return

}

func GetChunkNumber(x, y int) int {
	result := 0

	//Work out xdiff
	result += x / ChunkSize

	//Work out ydiff
	result += y / MapWidth

	return result

}

func (cp *Chunk) InsertUsingCopy(conn *pgx.Conn) (copyCount int) {
	var chunk Chunk
	chunk = *cp
	l := len(chunk)
	if l == 0 {
		return
	}
	rows := chunk.ToInterface()
	copyCount, err := conn.CopyFrom(
		pgx.Identifier{"world", "tiles"},
		TileSqlCols,
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		logrus.Fatalln(err)
	}
	if l != copyCount {
		logrus.Fatalln(l, " is not equal to ", copyCount)
	}
	return
}
