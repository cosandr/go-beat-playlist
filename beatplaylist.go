package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

const configPath = "./config.json"

var conf Config
var installedSongs Playlist
var allPlaylists map[string]Playlist

var rePlayExt *regexp.Regexp = regexp.MustCompile(`(\.json$|\.bplist$)`)

func init() {
	c, err := NewConfig(configPath)
	if err != nil {
		panic(err)
	}
	conf = c
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
	allPlaylists = make(map[string]Playlist)
	for _, p := range newPlaylists {
		allPlaylists[p.Title] = p
	}
}

func mainMenu() {
	const helpText = `Beat Saber playlist editor written in Go

1: Show all read playlists and their songs
2: Show all installed song data
3: Songs not in any playlists
4: Songs missing from playlists
5: Create playlist sorted by ScoreSaber star difficulty
6: Create playlist sorted by PP using Song Browser data
7: Check local song hashes
0: Exit`
	for {
		fmt.Printf("%s\n", helpText)
		fmt.Printf("Loaded %d songs and %d playlists.\n", len(installedSongs.Songs), len(allPlaylists))
		fmt.Print("Select option: ")
		in := GetInputNumber()
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
		case 5:
			songsFromScoreSaber()
			// Reload
			loadAll()
		case 6:
			songsFromSongBrowser()
		case 7:
			// Check hashes
			checkLocalSongs()
			// Reload
			loadAll()
		default:
			fmt.Println("Invalid option")
		}
	}
}

func checkLocalSongs() {
	var ok int
	var fail []Song
	var mismatch []Song
	allSongs, err := DownloadScrapedData(false)
	if err != nil {
		fmt.Printf("Cannot download scraped data: %v", err)
		return
	}
	isOK := func(other Song) rune {
		// Look for hash
		for _, s := range allSongs.Songs {
			if s.Hash == other.Hash {
				return 'o'
			}
		}
		// Check for name matches
		for _, s := range allSongs.Songs {
			// Using Contains because scraped data has an extra ' - ' at the end
			if strings.Contains(other.Name, s.Name) {
				log.Debugf("%s name match with %s\n", s.Name, other.Name)
				// Only if the author matches as well
				if s.Mapper == s.Mapper || s.Mapper == s.Author {
					// Maybe try to download song info from ScoreSaber as well?
					// Scraped data might be out of date
					return 'm'
				}
			}
		}
		return 'f'
	}
	for _, s := range installedSongs.Songs {
		switch isOK(s) {
		case 'o':
			ok++
		case 'm':
			mismatch = append(mismatch, s)
			fmt.Printf("-> Mismatch: %s\n", s.String())
		case 'f':
			fail = append(fail, s)
			fmt.Printf("-> Cannot find: %s\n", s.String())
		}
	}
	var helpText = `## %d OK, %d mismatched, %d failed ##

1: Delete mismatches and failed
0: Back to main menu`
	fmt.Printf(helpText, ok, len(mismatch), len(fail))
	fmt.Println()
	fmt.Print("Select option: ")
	in := GetInputNumber()
	switch in {
	case 0:
		return
	case 1:
		if len(mismatch) > 0 {
			move := GetConfirm("Move mismatches to DeletedSongs instead of deleting? (Y/n) ")
			deleteSongsFromPlaylist(Playlist{Songs: mismatch}, move)
		}
		if len(fail) > 0 {
			move := GetConfirm("Move failed songs to DeletedSongs instead of deleting? (Y/n) ")
			deleteSongsFromPlaylist(Playlist{Songs: fail}, move)
		}
	}
}

func songsFromSongBrowser() {
	var helpText = `## %d songs from Song Browser data ##

1: Show songs
2: Add to playlist
0: Back to main menu`
	fmt.Print("Enter max number of songs to fetch: ")
	numSongs := GetInputNumber()
	ppSongs, err := (DownloadPPPlaylist(numSongs))
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		numSongs = len((ppSongs).Songs)
		fmt.Printf(helpText, numSongs)
		fmt.Println()
		fmt.Print("Select option: ")
		in := GetInputNumber()
		switch in {
		case 0:
			return
		case 1:
			for _, s := range ppSongs.Songs {
				fmt.Printf("-> %.2f PP: %s\n", s.PP, s.Name)
			}
		case 2:
			path := fmt.Sprintf("Top%dPP.bplist", numSongs)
			fmt.Printf("Saving as %s\n", path)
			path = fmt.Sprintf("%s/%s", conf.Playlists, path)
			if FileExists(path) {
				backup := GetConfirm("Backup existing file? (Y/n) ")
				if backup {
					err := os.Rename(path, rePlayExt.ReplaceAllString(path, ".bak"))
					if err != nil {
						fmt.Printf("Cannot backup: %v\n", err)
						continue
					}
				}
			}
			ppSongs.Title = fmt.Sprintf("Top %d PP", numSongs)
			ppSongs.Author = "Dre"
			err := ioutil.WriteFile(path, ppSongs.ToJSON(), 0755)
			if err != nil {
				fmt.Printf("Cannot write playlist: %v\n", err)
				continue
			}
			return
		}
	}
}

func songsFromScoreSaber() {
	var helpText = `## %d songs from ScoreSaber ##

1: Show songs
2: Add to playlist
0: Back to main menu`
	fmt.Print("Enter max number of songs to fetch: ")
	numSongs := GetInputNumber()
	starSongs, err := (DownloadStarsPlaylist(numSongs))
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		numSongs = len((starSongs).Songs)
		fmt.Printf(helpText, numSongs)
		fmt.Println()
		fmt.Print("Select option: ")
		in := GetInputNumber()
		switch in {
		case 0:
			return
		case 1:
			for _, s := range starSongs.Songs {
				fmt.Printf("-> %.2f stars: %s\n", s.Stars, s.Name)
			}
		case 2:
			path := fmt.Sprintf("Top%dStars.bplist", numSongs)
			fmt.Printf("Saving as %s\n", path)
			path = fmt.Sprintf("%s/%s", conf.Playlists, path)
			if FileExists(path) {
				backup := GetConfirm("Backup existing file? (Y/n) ")
				if backup {
					err := os.Rename(path, rePlayExt.ReplaceAllString(path, ".bak"))
					if err != nil {
						fmt.Printf("Cannot backup: %v\n", err)
						continue
					}
				}
			}
			starSongs.Title = fmt.Sprintf("Top %d Stars", numSongs)
			starSongs.Author = "Dre"
			err := ioutil.WriteFile(path, starSongs.ToJSON(), 0755)
			if err != nil {
				fmt.Printf("Cannot write playlist: %v\n", err)
				continue
			}
			return
		}
	}
}

func deleteSongsFromPlaylist(p Playlist, move bool) {
	for _, s := range p.Songs {
		if !move {
			err := os.RemoveAll(s.Path)
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
		fmt.Print("Select option: ")
		in := GetInputNumber()
		fmt.Println()
		switch in {
		case 0:
			return
		case 1:
			fmt.Print(orphansPlaylist.String())
		case 2:
			var outBytes []byte
			// Ask for playlist path
			path, exists := GetInputPlaylist(conf.Playlists)
			// Confirm override
			if exists {
				merging := GetConfirm("File already exists, merge? (Y/n) ")
				if merging {
					// Read existing playlist
					existing, err := MakePlaylist(path)
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
			move := GetConfirm("Move to DeletedSongs instead of deleting? (Y/n) ")
			deleteSongsFromPlaylist(orphansPlaylist, move)
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
3: Download
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
		fmt.Print("Select option: ")
		in := GetInputNumber()
		fmt.Println()
		switch in {
		case 0:
			return
		case 1:
			for _, p := range missingPlaylists {
				fmt.Println(p.String())
			}
		case 2:
			var writePlaylist Playlist
			for name, p := range missingPlaylists {
				songs := []Song{}
				for _, s := range allPlaylists[name].Songs {
					if s.Path != "" {
						songs = append(songs, s)
					}
				}
				if len(songs) == 0 {
					continue
				}
				writePlaylist = Playlist{
					Title:  p.Title,
					Author: p.Author,
					Image:  p.Image,
					File:   p.File,
					Songs:  songs,
				}
				outBytes := writePlaylist.ToJSON()
				path := p.File
				backup := GetConfirm(fmt.Sprintf("Backup %s? (Y/n) ", p.Title))
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
		case 3:
			for name, p := range missingPlaylists {
				fmt.Printf("--> Downloading missing from %s\n", name)
				for _, s := range p.Songs {
					fmt.Printf(" --> Downloading %s\n", s.String())
					_, err := DownloadSong(&s)
					if err != nil {
						fmt.Printf("  -> Failed: %v\n", err)
						continue
					}
					fmt.Println("  -> Success")
				}
			}
			return
		}
	}
}

func getMissingFromPlaylists() map[string]Playlist {
	var missing = make(map[string]Playlist)
	// Populate song paths
	for _, p := range allPlaylists {
		songs := []Song{}
		for _, s := range p.Songs {
			if len(s.Path) == 0 {
				songs = append(songs, s)
			}
		}
		if len(songs) > 0 {
			missing[p.Title] = Playlist{
				Title:  p.Title,
				Author: p.Author,
				Image:  p.Image,
				File:   p.File,
				Songs:  songs,
			}
		}
	}
	return missing
}

func main() {
	var debug bool

	// Parse arguments
	flag.BoolVar(&debug, "debug", false, "Debug logging")
	flag.Parse()
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	loadAll()
	mainMenu()
}
