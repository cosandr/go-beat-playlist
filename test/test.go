package main

import (
	"fmt"

	"github.com/cosandr/go-beat-playlist/download"
	mt "github.com/cosandr/go-beat-playlist/types"
)

func main() {
	// path := "C:/Program Files (x86)/Steam/steamapps/common/Beat Saber/Playlists/Top100Star.json"
	// path := "test/parse/songbrowser-ranked.json"
	// file, err := ioutil.ReadFile(path)
	// if err != nil {
	// 	panic(err)
	// }
	// var _ = file
	// p, err := mt.MakeSongBrowserPlaylist(&file)
	// fmt.Println(p.Debug())
	// p.SortByPP()
	// fmt.Println(p.Debug())
	s := mt.Song{Hash: "9bf202f68c333421c69ca6aa15c648d65d4a1e0f", Name: "Night Raid"}
	out, err := download.Song("test-nightraid", &s)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(s.Debug())
	fmt.Println()
	fmt.Println(out.Debug())
}
