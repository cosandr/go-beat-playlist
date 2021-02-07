package main

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Matches all invalid NTFS characters
var reInvalid *regexp.Regexp = regexp.MustCompile(`[<>:"\/\\|?*\n]+`)

// Playlist holds the filename, raw JSON content and list of songs
type Playlist struct {
	Author string
	File   string
	Image  string
	Songs  []Song
	Title  string
}

// String returns playlist title and its songs
func (p *Playlist) String() string {
	var ret string
	ret += fmt.Sprintf("%s\n--- %d SONGS ---\n", p.Title, len(p.Songs))
	for _, s := range p.Songs {
		ret += s.String() + "\n"
	}
	return ret
}

// Debug returns all playlist and song fields
func (p *Playlist) Debug() string {
	var ret string
	ret += fmt.Sprintf("Title: %s, Author: %s\nFile: %s\n--- %d SONGS ---\n", p.Title, p.Author, p.File, len(p.Songs))
	for _, s := range p.Songs {
		ret += s.Debug()
	}
	return ret
}

// Contains returns true if Song is in it
func (p *Playlist) Contains(comp Song) bool {
	for _, s := range p.Songs {
		if s.Equals(&comp) {
			return true
		}
	}
	return false
}

// SongPath returns the requested song's path, if it exists
func (p *Playlist) SongPath(comp Song) string {
	for _, s := range p.Songs {
		// Ignore if we have no path or if both key and hash are missing
		if s.Path == "" || (s.Hash == "" && s.Key == "") {
			continue
		}
		if s.Equals(&comp) {
			return s.Path
		}
	}
	return ""
}

// Installed sets the file path for all its songs, if they are present
func (p *Playlist) Installed(installed *Playlist) {
	var newSongs []Song
	for _, s := range p.Songs {
		newSong := s
		newSong.Path = installed.SongPath(newSong)
		newSongs = append(newSongs, newSong)
	}
	p.Songs = newSongs
}

// ToJSON returns a JSON representation of the Playlist
func (p *Playlist) ToJSON() []byte {
	var jSongs []SongJSON
	for _, s := range p.Songs {
		sj := SongJSON{
			Key:  s.Key,
			Hash: s.Hash,
			Name: s.Name,
		}
		jSongs = append(jSongs, sj)
	}
	j := PlaylistJSON{
		Title:  p.Title,
		Author: p.Author,
		Image:  p.Image,
		Count:  len(p.Songs),
		Songs:  jSongs,
	}
	var bytes bytes.Buffer
	json := json.NewEncoder(&bytes)
	json.SetEscapeHTML(false)
	json.SetIndent("", " ")
	err := json.Encode(j)
	if err != nil {
		fmt.Println(err)
	}
	return bytes.Bytes()
}

// Merge returns a new playlist merged with the argument playlist
func (p *Playlist) Merge(op *Playlist) Playlist {
	var songs []Song
	// Keep all songs in this playlist
	songs = append(songs, p.Songs...)
	// Add songs only in other playlist
	for _, s := range op.Songs {
		if !p.Contains(s) {
			songs = append(songs, s)
		}
	}
	return Playlist{
		Title:  p.Title,
		Author: p.Author,
		Image:  p.Image,
		File:   p.File,
		Songs:  songs,
	}
}

// SortByPP sorts this playlist by PP in descending order
func (p *Playlist) SortByPP() {
	sort.Slice(p.Songs, func(i, j int) bool {
		return p.Songs[i].PP > p.Songs[j].PP
	})
}

// Song holds information about each song
type Song struct {
	Author string
	Hash   string
	Key    string
	Mapper string
	Maps   []Beatmap
	Name   string
	Path   string
	PP     float64
	Stars  float64
	URL    string
}

// Equals returns true if the key or hash matches
func (s *Song) Equals(other *Song) bool {
	if s.Hash != "" && s.Hash == other.Hash {
		return true
	} else if s.Key != "" && s.Key == other.Key {
		return true
	}
	return false
}

// Merge returns a new song merged with the argument song
//
// Prioritizes self, that is, only adds missing fields
//
// BTW, I know this is awful
func (s *Song) Merge(os *Song) Song {
	var retSong Song
	if len(s.Path) == 0 {
		retSong.Path = os.Path
	} else {
		retSong.Path = s.Path
	}

	if len(s.Key) == 0 {
		retSong.Key = os.Key
	} else {
		retSong.Key = s.Key
	}

	if len(s.Hash) == 0 {
		retSong.Hash = os.Hash
	} else {
		retSong.Hash = s.Hash
	}

	if len(s.Name) == 0 {
		retSong.Name = os.Name
	} else {
		retSong.Name = s.Name
	}

	if len(s.Author) == 0 {
		retSong.Author = os.Author
	} else {
		retSong.Author = s.Author
	}

	if len(s.Mapper) == 0 {
		retSong.Mapper = os.Mapper
	} else {
		retSong.Mapper = s.Mapper
	}

	if s.PP == 0 {
		retSong.PP = os.PP
	} else {
		retSong.PP = s.PP
	}

	if s.Stars == 0 {
		retSong.Stars = os.Stars
	} else {
		retSong.Stars = s.Stars
	}

	if len(s.Maps) == 0 {
		retSong.Maps = os.Maps
	} else {
		retSong.Maps = s.Maps
	}

	if len(s.URL) == 0 {
		retSong.URL = os.URL
	} else {
		retSong.URL = s.URL
	}

	return retSong
}

// CalcHash calculates this song's hash using its Path
//
// song.Hash must end with a trailing slash
func (s *Song) CalcHash() (err error) {
	// sha1 hash of (info.dat contents + contents of diff.dat files in order listed in info.dat)
	infoPath, err := FindInfo(s.Path)
	if err != nil {
		log.Debugf("base: %s, info: %s, err: %v", s.Path, infoPath, err)
		return
	}
	var files = []string{infoPath}
	for _, bm := range s.Maps {
		files = append(files, s.Path+"/"+bm.File)
	}
	var buf bytes.Buffer
	for _, f := range files {
		file, errF := ioutil.ReadFile(f)
		if errF != nil {
			err = fmt.Errorf("%s hash failed: %v", s.Name, errF)
			return
		}
		buf.Write(file)
	}
	s.Hash = fmt.Sprintf("%x", sha1.Sum(buf.Bytes()))
	return
}

// String returns a string representation of the song
func (s *Song) String() string {
	var ret string
	if len(s.Name) > 0 {
		ret += s.Name
	} else {
		ret += "MISSING"
	}
	if len(s.Key) > 0 {
		ret += fmt.Sprintf(" [%s]", s.Key)
	} else if len(s.Hash) > 0 {
		ret += fmt.Sprintf(" [%s]", s.Hash)
	}
	return ret
}

// Debug returns a string with all values in song
func (s *Song) Debug() string {
	var ret string
	ret += fmt.Sprintf("Path: %s, URL: %s\n", s.Path, s.URL)
	ret += fmt.Sprintf("Name: %s, Aut: %s, Mapper: %s\n", s.Name, s.Author, s.Mapper)
	ret += fmt.Sprintf("Key: %s, Hash: %s, PP: %.2f, Stars: %.2f\n", s.Key, s.Hash, s.PP, s.Stars)
	ret += "Beatmaps: "
	for _, m := range s.Maps {
		ret += "[" + m.Debug() + "] "
	}
	return ret
}

// DirName returns this song's directory name
func (s *Song) DirName() string {
	var ret string
	var paran bool
	if len(s.Key) > 0 {
		ret = fmt.Sprintf("%s (%s", s.Key, s.Name)
		paran = true
	} else {
		ret = s.Name
	}
	if len(s.Author) > 0 {
		ret += fmt.Sprintf(" - %s", s.Author)
	} else if len(s.Mapper) > 0 {
		ret += fmt.Sprintf(" - %s", s.Mapper)
	}
	if paran {
		ret += ")"
	}
	// Strip invalid NTFS characters
	ret = reInvalid.ReplaceAllString(ret, "")
	// Remove trailing space
	ret = strings.TrimSuffix(ret, " ")
	return ret
}

// Beatmap holds information about a song's map, its difficulty, path to the map file and type (standard, 360, lightshow)
type Beatmap struct {
	Difficulty string
	File       string
	Type       string
}

// String returns a pretty type: difficulty string
func (bm *Beatmap) String() string {
	return fmt.Sprintf("\n%s: %s", bm.Type, bm.Difficulty)
}

// Debug returns a string with all of this map's values
func (bm *Beatmap) Debug() string {
	if len(bm.File) > 0 {
		return fmt.Sprintf("%s, %s (%s)", bm.Type, bm.Difficulty, bm.File)
	}
	return fmt.Sprintf("%s, %s", bm.Type, bm.Difficulty)
}

// Config is the internal config, storing various game paths
type Config struct {
	Base         string
	DeletedSongs string
	Playlists    string
	Songs        string
}

// NewConfig reads the config at `path` and returns a `Config` object
//
// Will check for valid game path, creating missing directories (Playlists/CustomLevels)
func NewConfig(path string) (c Config, err error) {
	var jc ConfigJSON
	file, err := ioutil.ReadFile(path)
	if err == nil {
		errJ := json.Unmarshal(file, &jc)
		if errJ != nil {
			err = fmt.Errorf("Cannot parse %s: %v", path, errJ)
			return
		}
	} else {
		// Try to run without config file
		err = nil
	}
	// Default to C Steam folder
	if len(jc.Game) == 0 {
		c.Base = "C:/Program Files (x86)/Steam/steamapps/common/Beat Saber"
	} else {
		c.Base = NewPath(jc.Game)
	}
	// Check for valid game path
	for {
		if !FileExists(c.Base + "/Beat Saber.exe") {
			fmt.Printf("game not found at %s, enter game path: ", c.Base)
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				c.Base = scanner.Text()
				break
			}
		} else {
			break
		}
	}
	// Write to config
	if jc.Game != c.Base {
		jc.Game = c.Base
		file, errJ := json.MarshalIndent(&jc, "", " ")
		if errJ != nil {
			fmt.Printf("cannot marshal config file: %v", errJ)
		} else {
			if errJ = ioutil.WriteFile(path, file, 0644); errJ != nil {
				fmt.Printf("cannot marshal config file: %v", errJ)
			} else {
				fmt.Printf("updated config file %s\n", path)
			}
		}
	}
	mkdirMap := map[string]string{
		"Playlists":     c.Base + "/Playlists",
		"Custom songs":  c.Base + "/Beat Saber_Data/CustomLevels",
		"Deleted songs": c.Base + "/DeletedSongs",
	}
	for k, v := range mkdirMap {
		if !DirExists(v) {
			err = os.MkdirAll(v, 0755)
			if err != nil {
				return
			}
			fmt.Printf("%s folder %s created\n", k, v)
		}
	}
	c.Playlists = mkdirMap["Playlists"]
	c.Songs = mkdirMap["Custom songs"]
	c.DeletedSongs = mkdirMap["Deleted songs"]
	return
}

// StringSet a set for strings, useful for keeping track of elements
type StringSet map[string]struct{}

// Contains returns true if `v` is in the set
func (s StringSet) Contains(v string) bool {
	_, ok := s[v]
	return ok
}
