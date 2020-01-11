package main

import (
	"regexp"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/cosandr/go-beat-playlist/input"
	mt "github.com/cosandr/go-beat-playlist/types"
)

var conf mt.Config
var installedSongs mt.Playlist
var allPlaylists map[string]mt.Playlist

var rePlayExt *regexp.Regexp = regexp.MustCompile(`(\.json$|\.bplist$)`)

func init() {
	c, err := mt.NewConfig("./config.json")
	if err != nil {
		panic(err)
	}
	conf = c
	loadAll()
}

func loadAll() {
	newInstalled, err := readInstalledSongs(conf.Songs)
	if err != nil {
		panic(err)
	}
	installedSongs = newInstalled
	newPlaylists, err := readAllPlaylists(conf.Playlists)
	if err != nil {
		panic(err)
	}
	allPlaylists = make(map[string]mt.Playlist)
	for _, p := range newPlaylists {
		allPlaylists[p.Title] = p
	}
}

func mainMenu() {
	const helpText = `Beat Saber playlist editor written in Go
1: Show all read playlists and their songs
2: Show all installed song data
3: Show songs not in any playlists
4: Show songs missing from playlists
0: Exit`
	for {
		fmt.Printf("\n%s\n", helpText)
		fmt.Printf("Loaded %d songs and %d playlists.\n", len(installedSongs.Songs), len(allPlaylists))
		in := input.GetInputNumber()
		fmt.Println()
		switch in {
		case 0:
			return
		case 1:
			printAllPlaylists()
		case 2:
			fmt.Println(installedSongs.String())
		case 3:
			songsWithoutPlaylists()
			// Reload
			loadAll()
		case 4:
			missingFromPlaylists()
			// Reload
			loadAll()
		default:
			fmt.Println("Invalid option")
		}
	}
}

// songsWithoutPlaylists provides the UX for handling songs without playlists
func songsWithoutPlaylists() {
	var helpText = `## %d songs without playlists ##
1: Show songs
2: Add to playlist
3: Move or delete all
0: Back to main menu`
	for {
		orphansPlaylist := getSongsWithoutPlaylists()
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
			var outBytes []byte
			// Ask for playlist path
			path, exists := input.GetInputPlaylist(conf.Playlists)
			// Confirm override
			if exists {
				merging := input.GetConfirm("File already exists, merge? (Y/n) ")
				if merging {
					// Read existing playlist
					existing, err := mt.MakePlaylist(path)
					if err != nil {
						fmt.Printf("Cannot read playlist: %v\n", err)
						continue
					}
					// Merge with orphans
					writePlaylist := existing.Merge(&orphansPlaylist)
					outBytes = writePlaylist.ToJSON()
					fmt.Println("Merging orphans with playlist")
				} else {
					outBytes = orphansPlaylist.ToJSON()
					fmt.Println("Writing new playlist")
				}
			}
			err := ioutil.WriteFile(path, outBytes, 0755)
			if err != nil {
				fmt.Printf("Cannot write playlist: %v\n", err)
				continue
			}
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

func missingFromPlaylists() {
	var helpText = `## %d songs missing from all playlists ##
%s
1: Show songs
2: Remove from playlists
0: Back to main menu`
	for {
		missingPlaylists := getMissingFromPlaylists()
		var missingTotal int
		var missingSummary string
		for _, p := range missingPlaylists {
			missingTotal += len(p.Songs)
			missingSummary += fmt.Sprintf("-> %d from %s\n", len(p.Songs), p.Title)
		}
		fmt.Printf(helpText, missingTotal, missingSummary)
		fmt.Println()
		in := input.GetInputNumber()
		fmt.Println()
		switch in {
		case 0:
			return
		case 1:
			for _, p := range missingPlaylists {
				fmt.Println(p.String())
			}
		case 2:
			var writePlaylist mt.Playlist
			for name, p := range missingPlaylists {
				songs := []mt.Song{}
				for _, s := range allPlaylists[name].Songs {
					if s.Path != "" {
						songs = append(songs, s)
					}
				}
				if len(songs) == 0 {
					continue
				}
				writePlaylist = mt.Playlist{
					Title: p.Title,
					Author: p.Author,
					Image: p.Image,
					File: p.File,
					Songs: songs,
				}
				outBytes := writePlaylist.ToJSON()
				path := p.File
				backup := input.GetConfirm(fmt.Sprintf("Backup %s? (Y/n) ", p.Title))
				if backup {
					err := os.Rename(path, rePlayExt.ReplaceAllString(path, ".bak"))
					if err != nil {
						fmt.Printf("Cannot backup %s: %v\n", p.Title, err)
						continue
					}
				}
				err := ioutil.WriteFile(path, outBytes, 0755)
				if err != nil {
					fmt.Printf("Cannot write playlist: %v\n", err)
					continue
				}
			}
			return
		}
	}
}

func getMissingFromPlaylists() map[string]mt.Playlist {
	var missing = make(map[string]mt.Playlist)
	// Populate song paths
	for _, p := range allPlaylists {
		songs := []mt.Song{}
		for _, s := range p.Songs {
			if s.Path == "" {
				songs = append(songs, s)
			}
		}
		if len(songs) > 0 {
			missing[p.Title] = mt.Playlist{
				Title: p.Title,
				Author: p.Author,
				Image: p.Image,
				File: p.File,
				Songs: songs,
			}
		}
	}
	return missing
}


func main() {
	var timing bool
	var startTimes = make(map[string]time.Time)
	var endTimes = make(map[string]time.Time)

	// Parse arguments
	flag.BoolVar(&timing, "timing", false, "Enable timing")
	flag.Parse()

	mainMenu()

	if timing {
		for k, v := range endTimes {
			fmt.Printf("%s in: %s\n", k, (v.Sub(startTimes[k]).String()))
		}
	}
}
