package main

import (
	"fmt"

	"github.com/cosandr/go-beat-playlist/download"
	// mt "github.com/cosandr/go-beat-playlist/types"
)

func main() {
	// manualHash()
	// path := "C:/Program Files (x86)/Steam/steamapps/common/Beat Saber/Playlists/Top100Star.json"
	// path := "parse/beatsaver-api.json"
	// file, err := ioutil.ReadFile(path)
	// if err != nil {
	// 	panic(err)
	// }
	// var _ = file
	// p, err := mt.MakeBeatSaverPlaylist(&file)
	// fmt.Printf("%s\n", p.Debug())
	// fmt.Println(p.Songs[110].Debug())
	p, err := (download.FetchByStars(10))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p.Debug())
}
