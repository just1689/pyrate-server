package main

import (
	"github.com/just1689/pyrate-server/db"
	"github.com/just1689/pyrate-server/maps"
	"github.com/just1689/pyrate-server/model"
	"github.com/sirupsen/logrus"
	"sync"
	"sync/atomic"
	"time"
)

const workers = 64

var ops uint64
var rowsWritten uint64

var wStop []chan bool

var mapWidth = 5000
var mapHeight = 5000
var chunkDiff = 50

func main() {
	start := time.Now()

	//Figure out when all work is done
	var wg sync.WaitGroup

	c := model.GenerateWaterChunks(0, mapWidth, 0, mapHeight, chunkDiff)

	//Start workers go generate chunks
	wg.Add(1)
	go func() {
		for i := 0; i < workers; i++ {

			//Run worker in go routine
			startWorker(c, &wg)

			//Give the DB some breathing room
			time.Sleep(10 * time.Millisecond)
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

func startWorker(in chan model.Chunk, wg *sync.WaitGroup) {
	wg.Add(1)
	conn, err := db.Connect()
	if err != nil {
		logrus.Fatalln(err)
	}
	go func() {
		for chunk := range in {
			wg.Add(1)
			chunkID := atomic.AddUint64(&ops, 1)
			start := time.Now()
			ct := maps.RandomizeChunckType()
			maps.GenerateChunk(chunk, ct)
			chunk.InsertUsingCopy(conn)
			logrus.Infoln("Inserted chunk", chunkID, "took", time.Since(start), "(", ct, ")")
			wg.Done()
			atomic.AddUint64(&rowsWritten, uint64(len(chunk)))
		}
		conn.Close()
		wg.Done()
	}()

}
