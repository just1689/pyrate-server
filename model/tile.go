package model

import (
	"fmt"
	"github.com/jackc/pgx"
	"log"
	"time"
)

type Tile struct {
	ID       string
	X, Z     int
	TileType TileType
	TileSkin string
}

func (tile Tile) ToString() string {
	return fmt.Sprint("ID:", tile.ID, "x:", tile.X, "y:", tile.Z)
}

func UpdateTile(conn *pgx.Conn, t *Tile) (ok bool) {
	_, err := conn.Exec("update world.tiles set x=$1, y=$2, tile_type=$3, tile_skin=$4 where id=$5", t.X, t.Z, t.TileType, t.TileSkin, t.ID)
	if err != nil {
		log.Println(err)
		ok = false
		return
	}
	ok = true
	return

}

func InsertTile(conn *pgx.Conn, t *Tile) (ok bool) {
	_, err := conn.Exec("insert into world.tiles values($1, $2, $3, $4, $5)", t.ID, t.X, t.Z, t.TileType, t.TileSkin)
	if err != nil {
		log.Println(err)
		ok = false
		return
	}
	ok = true
	return

}

func GetTilesChunkAsync(conn *pgx.Conn, x1, x2, z1, z2 int) (c chan *Tile) {
	c = make(chan *Tile, 1024)
	go func() {
		count := 0
		rows, _ := conn.Query("select * from world.tiles where x>=$1 and x<=$2 and y>=$3 and y<=$4", x1, x2, z1, z2)
		for rows.Next() {
			tile := &Tile{}
			err := rows.Scan(&tile.ID, &tile.X, &tile.Z, &tile.TileType, &tile.TileSkin)
			if err != nil {
				log.Println(err)
				return
			}
			c <- tile
			count++
		}
		fmt.Println("GetTilesChunkAsync() returned rows:", count)
		close(c)
	}()
	return
}

func GetAllTilesAsync(conn *pgx.Conn) (c chan *Tile) {
	c = make(chan *Tile, 1024)
	go func() {
		fmt.Println("GetAllTilesAsync()")
		start := time.Now()
		count := 0
		rows, _ := conn.Query("select * from world.tiles order by (x + y)")
		for rows.Next() {
			tile := &Tile{}
			err := rows.Scan(&tile.ID, &tile.X, &tile.Z, &tile.TileType, &tile.TileSkin)
			if err != nil {
				log.Println(err)
				return
			}
			c <- tile
			count++
		}
		fmt.Println("Sent:", count, "rows in", time.Since(start))
		close(c)
	}()

	return
}

func GetTileAtAsync(conn *pgx.Conn, id string) (tile *Tile, err error) {
	rows, _ := conn.Query("select * from world.tiles where id=$1", id)
	for rows.Next() {
		tile = &Tile{}
		err = rows.Scan(&tile.ID, &tile.X, &tile.Z, &tile.TileType, &tile.TileSkin)
	}
	return
}
