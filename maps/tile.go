package maps

import (
	"github.com/jackc/pgx"
	"log"
)

type Tile struct {
	ID       string
	X, Y     int
	TileType TileType
	TileSkin string
}

func UpdateTile(conn *pgx.Conn, t *Tile) (ok bool) {
	_, err := conn.Exec("update pyrates.tiles set x=$1, y=$2, tile_type=$3, tile_skin=$4 where id=$5", t.X, t.Y, t.TileType, t.TileSkin, t.ID)
	if err != nil {
		log.Println(err)
		ok = false
		return
	}
	ok = true
	return

}

func InsertTile(conn *pgx.Conn, t *Tile) (ok bool) {
	_, err := conn.Exec("insert into pyrates.tiles values($1, $2, $3, $4, $5)", t.ID, t.X, t.Y, t.TileType, t.TileSkin)
	if err != nil {
		log.Println(err)
		ok = false
		return
	}
	ok = true
	return

}

func GetTilesChunk(x1, x2, y1, y2 int) (tiles Chunk, err error) {
	rows, _ := conn.Query("select * from pyrates.tiles where x>=$1 and x<=$2 and y>=$3 and y<=$4", x1, x2, y1, y2)
	for rows.Next() {
		tile := &Tile{}
		err = rows.Scan(&tile.ID, &tile.X, &tile.Y, &tile.TileType, &tile.TileSkin)
		if err != nil {
			log.Println(err)
			return
		}
		tiles = append(tiles, tile)

	}
	return
}

func GetTileAt(id string) (tile *Tile, err error) {
	rows, _ := conn.Query("select * from pyrates.tiles where id=$1", id)
	for rows.Next() {
		tile = &Tile{}
		err = rows.Scan(&tile.ID, &tile.X, &tile.Y, &tile.TileType, &tile.TileSkin)
	}
	return
}
