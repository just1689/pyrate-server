package main

import (
	"fmt"
	"github.com/just1689/pyrate-server/maps"
	"sync"
)

func main() {

	var wg *sync.WaitGroup

	rote := Rote{}
	rote.Start(4)
	rote.in <- BuildUpdater()
	rote.in <- BuildUpdater()
	rote.in <- BuildUpdater()
	rote.in <- BuildUpdater()

	pool := CreatePool(4, wg)

	chunks := ReadChunks(wg)
	for chunk := range chunks {
		pool <- func() {
			l := chunk
			maps.GenerateChunk(l, maps.SingleIsland)
			for _, tile := range l {
				(rote.Next()) <- tile
				wg.Done()
			}
		}
	}
	wg.Wait()

}

func ReadChunks(inc *sync.WaitGroup) (out chan maps.Chunk) {
	out = make(chan maps.Chunk, 40) //40 chunks
	go func() {
		//fmt.Println("Fetching x: ", x1, "to", x2, " and y: ", y1, "to", y2)
		conn, err := maps.Connect()
		if err != nil {
			fmt.Println(err)
			return
		}

		for x := 0; x < 1000; x += 50 {
			for y := 0; y < 1000; y += 50 {

				x1 := x
				x2 := x + 49

				y1 := y
				y2 := y + 49

				chunk, err := maps.GetTilesChunk(conn, x1, x2, y1, y2)
				if err != nil {
					fmt.Println(err)
					return
				}

				inc.Add(1)
				out <- chunk

			}
		}
		close(out)

	}()
	return
}

func BuildUpdater() (in chan *maps.Tile) {
	in = make(chan *maps.Tile, 1000)
	c, err := maps.Connect()
	if err != nil {
		fmt.Println(err)
	}
	go func() {
		for item := range in {
			maps.UpdateTile(c, item)
		}

	}()
	return
}

type TinyFunction func()

func CreatePool(size int, wg *sync.WaitGroup) (in chan TinyFunction) {
	in = make(chan TinyFunction, 256)
	go func() {

		for i := 0; i < size; i++ {
			go func() {
				for f := range in {
					wg.Add(1)
					f()
					wg.Done()
				}
			}()
		}
	}()
	return

}

type Rote struct {
	in chan chan *maps.Tile
}

func (r Rote) Start(size int) {
	r.in = make(chan chan *maps.Tile, size)
}

func (r Rote) Next() chan *maps.Tile {
	res := <-r.in
	r.in <- res
	return res
}
