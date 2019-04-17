package main

import (
	"fmt"
	"github.com/just1689/pyrate-server/db"
	"github.com/just1689/pyrate-server/maps"
	"github.com/just1689/pyrate-server/model"
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
			maps.GenerateChunk(l, model.SingleIsland)
			for _, tile := range l {
				(rote.Next()) <- tile
				wg.Done()
			}
		}
	}
	wg.Wait()

}

func ReadChunks(inc *sync.WaitGroup) (out chan model.Chunk) {
	out = make(chan model.Chunk, 40) //40 chunks
	go func() {
		//fmt.Println("Fetching x: ", x1, "to", x2, " and y: ", y1, "to", y2)
		conn, err := db.Connect()
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

				chunk, err := model.GetTilesChunk(conn, x1, x2, y1, y2)
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

func BuildUpdater() (in chan *model.Tile) {
	in = make(chan *model.Tile, 1000)
	c, err := db.Connect()
	if err != nil {
		fmt.Println(err)
	}
	go func() {
		for item := range in {
			model.UpdateTile(c, item)
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
	in chan chan *model.Tile
}

func (r Rote) Start(size int) {
	r.in = make(chan chan *model.Tile, size)
}

func (r Rote) Next() chan *model.Tile {
	res := <-r.in
	r.in <- res
	return res
}
