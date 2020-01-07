package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/cosandr/go-beat-playlist/input"
	mt "github.com/cosandr/go-beat-playlist/types"
)

func decideRun(playlists *[]mt.Playlist, installed *mt.Playlist) {
	const helpText = `Beat Saber Playlist editor written in Go
1: Show all read playlists and their songs
2: Show all installed song data
3: Show songs not in any playlists
4: Show songs missing from playlists
9: Exit`
	for {
		fmt.Printf("\n%s\n", helpText)
		fmt.Printf("Loaded %d songs and %d playlists.\n", len((*installed).Songs), len(*playlists))
		in := input.GetInputNumber()
		fmt.Println()
		switch in {
		case 1:
			printAllPlaylists(playlists)
		case 2:
			fmt.Println(installed.String())
		case 3:
			orphansPlaylist := songsWithoutPlaylists(playlists, installed)
			fmt.Print(orphansPlaylist.String())
		case 4:
			missingFromPlaylists(playlists, installed)
		case 9:
			return
		default:
			fmt.Println("Invalid option")
		}
	}
}

// songsWithoutPlaylists returns a Playlist of songs not already in any playlists
func songsWithoutPlaylists(playlists *[]mt.Playlist, installed *mt.Playlist) mt.Playlist {
	var orphans []mt.Song
	var isOrphan bool
	for _, s := range (*installed).Songs {
		isOrphan = true
		for _, p := range *playlists {
			if p.Contains(s) {
				isOrphan = false
				break
			}
		}
		if isOrphan {
			orphans = append(orphans, s)
		}
	}
	return mt.Playlist{Title: "Orphans", Songs: orphans}
}

func missingFromPlaylists(playlists *[]mt.Playlist, installed *mt.Playlist) {
	var missing = make(map[string][]mt.Song)
	// Populate song paths
	for _, p := range *playlists {
		songs := []mt.Song{}
		for _, s := range p.Songs {
			if s.Path == "" {
				songs = append(songs, s)
			}
		}
		if len(songs) > 0 {
			missing[p.Title] = songs
		}
	}
	for k, v := range missing {
		fmt.Printf("\t --- %d missing from %s ---\n", len(v), k)
		for _, s := range v {
			fmt.Println(s.String())
		}
	}
}

func main() {
	c, err := readCfg("./config.json")
	if err != nil {
		panic(err)
	}
	var customSongs = c.Game + "/Beat Saber_Data/CustomLevels"

	var timing bool
	var startTimes = make(map[string]time.Time)
	var endTimes = make(map[string]time.Time)

	// Parse arguments
	flag.BoolVar(&timing, "timing", false, "Enable timing")
	flag.Parse()

	startTimes["Read installed"] = time.Now()
	installed, err := readInstalledSongs(customSongs)
	if err != nil {
		panic(err)
	}
	endTimes["Read installed"] = time.Now()

	startTimes["Read playlists"] = time.Now()
	playlists, err := readAllPlaylists(c.Playlist, &installed)
	if err != nil {
		panic(err)
	}
	endTimes["Read playlists"] = time.Now()

	decideRun(&playlists, &installed)

	if timing {
		for k, v := range endTimes {
			fmt.Printf("%s in: %s\n", k, (v.Sub(startTimes[k]).String()))
		}
	}
}

func writePlaylist(path string, playlist *mt.Playlist) {
	err := ioutil.WriteFile(path, playlist.ToJSON(), 0755)
	if err != nil {
		err = fmt.Errorf("JSON write error: %v ", err)
		return
	}
}

func readCfg(path string) (c mt.ConfigJSON, err error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("Config file error: %v ", err)
		return
	}
	err = json.Unmarshal(file, &c)
	if err != nil {
		err = fmt.Errorf("Cannot parse %s: %v", path, err)
		return
	}
	return
}

func readAllPlaylists(path string, installed *mt.Playlist) (playlists []mt.Playlist, err error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	for _, file := range files {
		p, readErr := mt.MakePlaylist(path + "/" + file.Name())
		if readErr != nil {
			fmt.Println(readErr)
			continue
		}
		p.Installed(installed)
		playlists = append(playlists, p)
	}
	return
}

func readInstalledSongs(path string) (p mt.Playlist, err error) {
	var songs []mt.Song
	err = filepath.Walk(path, func(subpath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == "info.dat" {
			s, makeErr := mt.MakeSong(subpath)
			if makeErr != nil {
				fmt.Printf("Cannot create song: %v\n", makeErr)
				return nil
			}
			songs = append(songs, s)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path: %v\n", err)
		return
	}
	p = mt.Playlist{Title: "Installed Songs", Songs: songs}
	return
}

func printAllPlaylists(playlists *[]mt.Playlist) {
	for _, p := range *playlists {
		fmt.Println(p.String())
	}
}
