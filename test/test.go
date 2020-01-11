package main

import (
	"fmt"
	"io/ioutil"

	mt "github.com/cosandr/go-beat-playlist/types"
)

func main() {
	// path := "C:/Program Files (x86)/Steam/steamapps/common/Beat Saber/Playlists/Top100Star.json"
	path := "test/parse/songbrowser-ranked.json"
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var _ = file
	p, err := mt.MakeSongBrowserPlaylist(&file)
	fmt.Println(p.Debug())
	p.SortByPP()
	fmt.Println(p.Debug())

}
