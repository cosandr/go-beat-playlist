package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	mt "github.com/cosandr/go-beat-playlist/types"
)

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

func readAllPlaylists(path string) (allPlaylists []mt.Playlist, err error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	for _, file := range files {
		p, readErr := readPlaylist(path + "/" + file.Name())
		if readErr != nil {
			fmt.Println(readErr)
			continue
		}
		allPlaylists = append(allPlaylists, p)
	}
	return
}

func readPlaylist(path string) (p mt.Playlist, err error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	err = json.Unmarshal(file, &p.Raw)
	if err != nil {
		return
	}
	p.File = path
	p.ParseRaw()
	return
}

func printAllPlaylists(all []mt.Playlist) {
	for _, p := range all {
		fmt.Println(p.String())
	}
}

func readInstalledSongs(path string) (songs []mt.Song, err error) {
	err = filepath.Walk(path, func(subpath string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", subpath, err)
			return err
		}
		if info.Name() == "info.dat" {
			s := mt.Song{Path: strings.TrimSuffix(subpath, "info.dat")}
			s.ParseRaw()
			s.CalcHash()
			songs = append(songs, s)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path: %v\n", err)
		return
	}
	return
}

func main() {
	c, err := readCfg("./config.json")
	if err != nil {
		panic(err)
	}
	// var customSongs = c.Game + "/Beat Saber_Data/CustomLevels"

	var timing bool
	var startTimes = make(map[string]time.Time)
	var endTimes = make(map[string]chan time.Time)
	startTimes["MAIN"] = time.Now()

	// Parse arguments
	flag.BoolVar(&timing, "timing", false, "Enable timing")
	flag.Parse()
	startTimes["Read playlists"] = time.Now()
	allPlaylists, err := readAllPlaylists(c.Playlist)
	endTimes["Read playlists"] = make(chan time.Time, 1)
	endTimes["Read playlists"] <- time.Now()

	// printAllPlaylists(allPlaylists)
	var _ = allPlaylists

	// installed, err := readInstalledSongs(customSongs)
	// if err != nil {
	// 	panic(err)
	// }

	// for _, s := range installed {
	// 	fmt.Println(s.String())
	// }

	testS := mt.Song{Path: "./test/nightraid"}
	testS.ParseRaw()
	testS.CalcHash()
	fmt.Println(testS.Debug())

	endTimes["MAIN"] = make(chan time.Time, 1)
	endTimes["MAIN"] <- time.Now()
	if timing {
		for k, v := range endTimes {
			fmt.Printf("%s ran in: %s\n", k, ((<-v).Sub(startTimes[k]).String()))
		}
	}

}
