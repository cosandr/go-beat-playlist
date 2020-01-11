package main

import (
	"fmt"
	"io/ioutil"

	mt "github.com/cosandr/go-beat-playlist/types"
)

func main() {
	// manualHash()
	// path := "C:/Program Files (x86)/Steam/steamapps/common/Beat Saber/Playlists/Top100Star.json"
	path := "parse/beatsaver-api.json"
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var _ = file
	// p, err := mt.MakeBeatSaverPlaylist(&file)
	// fmt.Printf("%s\n", p.Debug())
	// fmt.Println(p.Songs[110].Debug())
	path = "parse/playlist.bplist"
	p, err := mt.MakePlaylist(path)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", p.String())
}
