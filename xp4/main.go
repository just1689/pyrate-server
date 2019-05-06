package main

import (
	"fmt"
	"github.com/just1689/pyrate-server/db"
	"github.com/just1689/pyrate-server/maps"
	"github.com/just1689/pyrate-server/model"
	"sync"
	"time"
)

const workers = 32

func main() {

	//Figure out when all work is done
	var wg sync.WaitGroup

	//Produces x1, x2, y1, y2 for each chunk to be generated
	work := giveMeWork(&wg)

	//Start workers go generate chunks
	for i := 0; i < workers; i++ {

		//Run worker in go routine
		StartWorker(fmt.Sprint(i), work, &wg)

		//Give the DB some breathing room
		time.Sleep(1 * time.Second)
	}

	wg.Wait()

}

func giveMeWork(wg *sync.WaitGroup) chan *Work {
	wg.Add(1)
	result := make(chan *Work)
	go func() {
		for x := 0; x < 1000; x += 50 {
			for y := 0; y < 1000; y += 50 {
				result <- &Work{
					X1: x,
					X2: x + 49,
					Y1: y,
					Y2: y + 49,
				}
			}
		}
		close(result)
		wg.Done()
	}()
	return result

}

type Work struct {
	X1, X2, Y1, Y2 int
}

func StartWorker(name string, in chan *Work, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		fmt.Println("Worker", name, "STARTING")
		conn, err := db.Connect()
		if err != nil {
			fmt.Println(err)
			return
		}
		for work := range in {
			start := time.Now()
			chunk := getTilesChunk(work.X1, work.X2, work.Y1, work.Y2)
			ct := maps.RandomizeChunckType()
			maps.GenerateChunk(chunk, ct)
			for _, tile := range chunk {
				model.InsertTile(conn, tile)
			}
			d := time.Since(start)
			fmt.Println("Chunk took", d, "(", ct, ")")
		}
		wg.Done()

	}()

}

func getTilesChunk(x1 int, x2 int, z1 int, z2 int) (c model.Chunk) {
	for x := x1; x <= x2; x++ {
		for z := z1; z <= z2; z++ {
			t := model.Tile{
				ID:       fmt.Sprint(x, ".", z),
				X:        x,
				Z:        z,
				TileType: model.TileTypeWater,
				TileSkin: "",
			}
			c = append(c, &t)
		}
	}
	return
}
