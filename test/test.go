package main

import (
	"fmt"

	mt "github.com/cosandr/go-beat-playlist/types"
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
	path := "test/parse/playlist.bplist"
	p, _ := mt.MakePlaylist(path)
	path = "C:/Program Files (x86)/Steam/steamapps/common/Beat Saber/Playlists/Top100Star.json"
	existing, _ := mt.MakePlaylist(path)
	merged := existing.Merge(&p)
	fmt.Printf("Existing: T: %s, A: %s, F: %s, S: %d\nNew: T: %s, A: %s, F: %s, S: %d\nMerged: T: %s, A: %s, F: %s, S: %d\n",
		existing.Title, existing.Author, existing.File, len(existing.Songs),
		p.Title, p.Author, p.File, len(p.Songs),
		merged.Title, merged.Author, merged.File, len(merged.Songs))
}
