package types

import (
	"fmt"
	"strings"
)

// ConfigJSON is the structure of the config.json file
type ConfigJSON struct {
	Playlist string `json:"playlists"`
	Game     string `json:"game"`
}

// Playlist holds the filename, raw JSON content and list of songs
type Playlist struct {
	Title string
	File  string
	Raw   map[string]interface{}
	Songs []Song
}

// ParseRaw reads the raw data and fils in Title and Songs
func (p *Playlist) ParseRaw() {
	var tmp Song
	p.Title = p.Raw["playlistTitle"].(string)
	for _, s := range p.Raw["songs"].([]interface{}) {
		for k, v := range s.(map[string]interface{}) {
			str := ""
			// Convert to string if it is something else
			switch vv := v.(type) {
			case float64:
				str = fmt.Sprintf("%.0f", vv)
			case string:
				str = vv
			default:
				continue
			}
			switch k {
			case "songName":
				tmp.Name = str
			case "hash":
				tmp.hash = str
			case "key":
				tmp.key = str
			}
		}
		if tmp != (Song{}) {
			p.Songs = append(p.Songs, tmp)
			tmp = Song{}
		}
	}
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
	return ret
}

// Debug returns a string with all values in song
func (s *Song) Debug() string {
	return fmt.Sprintf("N: %s, K: %s, H: %s\nPP: %.2f, S: %.2f", s.Name, s.Key(), s.Hash(), s.PP, s.Stars)
}
