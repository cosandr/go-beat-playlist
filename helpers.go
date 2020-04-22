package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// getSongsWithoutPlaylists returns a Playlist of songs not already in any playlists
func getSongsWithoutPlaylists() Playlist {
	var orphans []Song
	var isOrphan bool
	for _, s := range installedSongs.Songs {
		isOrphan = true
		for _, p := range allPlaylists {
			if p.Contains(s) {
				isOrphan = false
				break
			}
		}
		if isOrphan {
			orphans = append(orphans, s)
		}
	}
	return Playlist{Title: "Orphans", Songs: orphans}
}

func readAllPlaylists(path string) (playlists []Playlist, err error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	for _, file := range files {
		if !rePlayExt.MatchString(file.Name()) {
			if !strings.HasSuffix(file.Name(), ".bak") {
				fmt.Printf("%s is not a valid playlist, skipping.", file.Name())
			}
			continue
		}
		p, readErr := MakePlaylist(path + "/" + file.Name())
		if readErr != nil {
			fmt.Println(readErr)
			continue
		}
		p.Installed(&installedSongs)
		playlists = append(playlists, p)
	}
	return
}

func readInstalledSongs(path string) (p Playlist, err error) {
	var songs []Song
	err = filepath.Walk(path, func(subpath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == "info.dat" {
			s, makeErr := MakeSong(subpath)
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
	p = Playlist{Title: "Installed Songs", Songs: songs}
	return
}

func printAllPlaylists() {
	for _, p := range allPlaylists {
		fmt.Println(p.String())
	}
}
