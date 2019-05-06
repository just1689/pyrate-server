package main

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/just1689/pyrate-server/db"
	"github.com/just1689/pyrate-server/maps"
	"github.com/just1689/pyrate-server/model"
	"github.com/sirupsen/logrus"
	"sync"
	"sync/atomic"
	"time"
)

const workers = 32
const writers = 4

var ops uint64
var rowsWritten uint64

var wStop []chan bool

var mapWidth = 100
var mapHeight = 100
var chunkDiff = 50

func main() {
	start := time.Now()

	//Figure out when all work is done
	var wg sync.WaitGroup

	//Have a channel that can hold enough work for every worker to be io busy some where else - then its not the blocker
	dbInChan := make(chan *model.Tile, chunkDiff*chunkDiff*writers)

	c := model.GenerateWaterChunks(0, mapWidth, 0, mapHeight, chunkDiff)

	//Start workers go generate chunks
	wg.Add(1)
	go func() {
		for i := 0; i < workers; i++ {

			//Run worker in go routine
			startWorker(fmt.Sprint(i), c, &wg, dbInChan)

			//Give the DB some breathing room
			time.Sleep(10 * time.Millisecond)
		}
		wg.Done()
	}()

	go func() {
		wg.Add(1)
		for i := 0; i < writers; i++ {
			s := make(chan bool)
			wStop = append(wStop, s)
			buildWriter(s, chunkDiff*chunkDiff, dbInChan)
		}
		wg.Done()
	}()

	wg.Wait()

	//Stop the DB writers
	for _, c := range wStop {
		c <- true
		close(c)
	}

	logrus.Infoln("Took: ", time.Since(start), "to write", atomic.AddUint64(&rowsWritten, 0), "rows")

}

func startWorker(name string, in chan model.Chunk, wg *sync.WaitGroup, dbChan chan *model.Tile) {
	wg.Add(1)
	go func() {
		for chunk := range in {
			chunkID := atomic.AddUint64(&ops, 1)
			start := time.Now()
			ct := maps.RandomizeChunckType()
			maps.GenerateChunk(chunk, ct)
			for _, tile := range chunk {
				dbChan <- tile
			}
			logrus.Infoln("Chunk", chunkID, "took", time.Since(start), "(", ct, ")")
		}
		wg.Done()
	}()

}

func buildWriter(stop chan bool, max int, in chan *model.Tile) {
	conn, err := db.Connect()
	if err != nil {
		logrus.Fatal(err)
	}
	go func() {
		var cache model.Chunk
		for {
			select {
			case <-stop:
				copyToDB(conn, cache)
				conn.Close()
				return
			case t := <-in:
				cache = append(cache, t)
				if cache.Size() >= max {
					copyToDB(conn, cache)
					cache = model.Chunk{}
				}
			}
		}
	}()
	return
}

func copyToDB(conn *pgx.Conn, chunk model.Chunk) {
	l := len(chunk)
	if l == 0 {
		return
	}

	if true {
		//logrus.Infoln("Wrote", len(chunk))
		atomic.AddUint64(&rowsWritten, uint64(l))
		return
	}

	rows := chunk.ToInterface()
	copyCount, err := conn.CopyFrom(
		pgx.Identifier{"world", "tiles"},
		model.TileSqlCols,
		pgx.CopyFromRows(rows),
	)
	logrus.Infoln("Wrote", l, "rows")
	atomic.AddUint64(&rowsWritten, uint64(l))
	if err != nil {
		logrus.Fatalln(err)
	}
	if l != copyCount {
		logrus.Fatalln(l, " is not equal to ", copyCount)
	}

}
