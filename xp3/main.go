package main

import (
	"fmt"
	"github.com/just1689/pyrate-server/maps"
	"sync"
	"time"
)

const workers = 8

func main() {

	//Figure out when all work is done
	var wg sync.WaitGroup

	//Produces x1, x2, y1, y2 for each chunk to be generated
	work := giveMeWork(&wg)

	for i := 0; i < workers; i++ {
		go Worker(fmt.Sprint(i), work, &wg)
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

func Worker(name string, in chan *Work, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		fmt.Println("Worker", name, "STARTING")
		conn, err := maps.Connect()
		if err != nil {
			fmt.Println(err)
			return
		}
		for work := range in {
			start := time.Now()
			chunk, err := maps.GetTilesChunk(conn, work.X1, work.X2, work.Y1, work.Y2)
			//fmt.Println("Worker", name, "received chunk starting at", work.X1, work.Y1)
			if err != nil {
				fmt.Println(err)
				return
			}
			maps.GenerateChunk(chunk, maps.SingleIsland)
			for _, tile := range chunk {
				maps.UpdateTile(conn, tile)
			}
			d := time.Since(start)
			fmt.Println("Chunk took", d)
		}
		wg.Done()

	}()

}
