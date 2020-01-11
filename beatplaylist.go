package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/cosandr/go-beat-playlist/input"
	mt "github.com/cosandr/go-beat-playlist/types"
)

var conf mt.Config

func init() {
	c, err := mt.NewConfig("./config.json")
	if err != nil {
		panic(err)
	}
	conf = c
}

func decideRun(playlists *[]mt.Playlist, installed *mt.Playlist) {
	const helpText = `Beat Saber Playlist editor written in Go
1: Show all read playlists and their songs
2: Show all installed song data
3: Show songs not in any playlists
4: Show songs missing from playlists
0: Exit`
	for {
		fmt.Printf("\n%s\n", helpText)
		fmt.Printf("Loaded %d songs and %d playlists.\n", len((*installed).Songs), len(*playlists))
		in := input.GetInputNumber()
		fmt.Println()
		switch in {
		case 0:
			return
		case 1:
			printAllPlaylists(playlists)
		case 2:
			fmt.Println(installed.String())
		case 3:
			songsWithoutPlaylists(playlists, installed)
			// Reload
			newInstalled, err := readInstalledSongs(conf.Songs)
			if err != nil {
				panic(err)
			}
			newPlaylists, err := readAllPlaylists(conf.Playlists, &newInstalled)
			if err != nil {
				panic(err)
			}
			installed = &newInstalled
			playlists = &newPlaylists
		case 4:
			missingFromPlaylists(playlists, installed)
		default:
			fmt.Println("Invalid option")
		}
	}
}

// songsWithoutPlaylists provides the UX for handling songs without playlists
func songsWithoutPlaylists(playlists *[]mt.Playlist, installed *mt.Playlist) {
	var helpText = `## %d songs without playlists ##
1: Show songs
2: Add to playlist
3: Delete or move all
0: Back to main menu`
	for {
		orphansPlaylist := getSongsWithoutPlaylists(playlists, installed)
		fmt.Printf(helpText, len((orphansPlaylist).Songs))
		fmt.Println()
		in := input.GetInputNumber()
		fmt.Println()
		switch in {
		case 0:
			return
		case 1:
			fmt.Print(orphansPlaylist.String())
		case 2:
			// Ask for playlist path
			path, exists := input.GetInputPlaylist(conf.Playlists)
			// Confirm override
			if exists {
				ok := input.GetConfirm("File already exists, override? (Y/n) ")
				if !ok {
					continue
				}
			}
			outBytes := orphansPlaylist.ToJSON()
			err := ioutil.WriteFile(path, outBytes, 0755)
			if err != nil {
				fmt.Printf("Cannot write playlist: %v\n", err)
				continue
			}
			fmt.Println("New playlist created")
			return
		case 3:
			move := input.GetConfirm("Move to DeletedSongs instead of deleting? (Y/n) ")
			for _, s := range orphansPlaylist.Songs {
				if !move {
					err := os.Remove(s.Path)
					if err != nil {
						fmt.Printf("Cannot delete %s: %v\n", s.String(), err)
						continue
					}
					fmt.Printf("Deleted %s\n", s.String())
				} else {
					err := os.Rename(s.Path, fmt.Sprintf("%s/%s", conf.DeletedSongs, s.DirName()))
					if err != nil {
						fmt.Printf("Cannot move %s: %v\n", s.String(), err)
						continue
					}
					fmt.Printf("Moved %s\n", s.String())
				}
			}
			return
		default:
			fmt.Println("Invalid option")
		}
	}
}

// getSongsWithoutPlaylists returns a Playlist of songs not already in any playlists
func getSongsWithoutPlaylists(playlists *[]mt.Playlist, installed *mt.Playlist) mt.Playlist {
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
	var timing bool
	var startTimes = make(map[string]time.Time)
	var endTimes = make(map[string]time.Time)

	// Parse arguments
	flag.BoolVar(&timing, "timing", false, "Enable timing")
	flag.Parse()

	startTimes["Read installed"] = time.Now()
	installed, err := readInstalledSongs(conf.Songs)
	if err != nil {
		panic(err)
	}
	endTimes["Read installed"] = time.Now()

	startTimes["Read playlists"] = time.Now()
	playlists, err := readAllPlaylists(conf.Playlists, &installed)
	if err != nil {
		panic(err)
	}
	endTimes["Read playlists"] = time.Now()
	fmt.Println(installed.Songs[0].Debug())
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
