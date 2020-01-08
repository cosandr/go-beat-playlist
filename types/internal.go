package types

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Playlist holds the filename, raw JSON content and list of songs
type Playlist struct {
	Title  string
	Author string
	File   string
	Songs  []Song
	json   *PlaylistJSON
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
		if s.Hash != "" && s.Hash == comp.Hash {
			return true
		} else if s.Key != "" && s.Key == comp.Key {
			return true
		}
	}
	return false
}

// SongPath returns the requested song's path, if it exists
func (p *Playlist) SongPath(comp Song) string {
	for _, s := range p.Songs {
		if s.Path == "" || (s.Hash == "" && s.Key == "") {
			continue
		}
		if s.Hash == comp.Hash {
			return s.Path
		} else if s.Key == comp.Key {
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
		newSong.Path = installed.SongPath(s)
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

// Song holds information about each song
type Song struct {
	Path   string
	Key    string
	Hash   string
	Name   string
	Author string
	Mapper string
	PP     float64
	Stars  float64
	Maps   []Beatmap
	URL    string
	json   *InfoJSON
}

// CalcHash calculates this song's hash using its Path
func (s *Song) CalcHash() {
	// sha1 hash of (info.dat contents + contents of diff.dat files in order listed in info.dat)
	var files = []string{s.Path + "info.dat"}
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
	s.Hash = fmt.Sprintf("%x", sha1.Sum(buf.Bytes()))
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
	ret += fmt.Sprintf("N: %s, K: %s, H: %s\nPP: %.2f, S: %.2f\n", s.Name, s.Key, s.Hash, s.PP, s.Stars)
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
