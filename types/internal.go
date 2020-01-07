package types

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"strings"
)

// Playlist holds the filename, raw JSON content and list of songs
type Playlist struct {
	Title string
	File  string
	Songs []Song
	json  *PlaylistJSON
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

// Song holds information about each song
type Song struct {
	Path  string
	key   string
	hash  string
	Name  string
	PP    float64
	Stars float64
	Maps  []Beatmap
	json  *InfoJSON
}

// CalcHash calculates this song's hash using its Path
func (s *Song) CalcHash() {
	// sha1 hash of (info.dat contents + contents of diff.dat files in order listed in info.dat)
	var files = []string{s.Path+"info.dat"}
	for _, bm := range s.Maps {
		files = append(files, s.Path+bm.File)
	}
	var buf bytes.Buffer
	for _, f := range files {
		file, err := ioutil.ReadFile(f)
		if err != nil {
			fmt.Printf("%s hash failed: %v\n", s.Name, err)
			return
		}
		buf.Write(file)
	}
	s.hash = fmt.Sprintf("%x", sha1.Sum(buf.Bytes()))
}

// Hash returns the hash in lower-case
func (s *Song) Hash() string {
	return strings.ToLower(s.hash)
}

// Key returns the key in lower case
func (s *Song) Key() string {
	return strings.ToLower(s.key)
}

// String returns a string representation of the song
func (s *Song) String() string {
	var ret string
	if len(s.Name) > 0 {
		ret += s.Name
	} else {
		ret += "MISSING"
	}
	if len(s.key) > 0 {
		ret += fmt.Sprintf(" [%s]", s.Key())
	} else if len(s.hash) > 0 {
		ret += fmt.Sprintf(" [%s]", s.Hash())
	}
	for _, m := range s.Maps {
		ret += m.String()
	}
	return ret
}

// Debug returns a string with all values in song
func (s *Song) Debug() string {
	var ret string
	ret += fmt.Sprintf("N: %s, K: %s, H: %s\nPP: %.2f, S: %.2f\n", s.Name, s.Key(), s.Hash(), s.PP, s.Stars)
	for _, m := range s.Maps {
		ret += m.Debug()
	}
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
	return fmt.Sprintf("Type %s, %s\n%s", bm.Type, bm.Difficulty, bm.File)
}
