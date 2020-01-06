package types

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
			default:
				continue
			}
		}
		p.Songs = append(p.Songs, tmp)
		tmp = Song{}
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

// Beatmap holds information about a song's map, its difficulty, path to the map file and type (standard, 360, lightshow)
type Beatmap struct {
	Difficulty string
	File string
	Type string
}

// String returns a pretty type: difficulty string
func (bm *Beatmap) String() string {
	return fmt.Sprintf("\n%s: %s", bm.Type, bm.Difficulty)
}

// Debug returns a string with all of this map's values
func (bm *Beatmap) Debug() string {
	return fmt.Sprintf("Type %s, %s\n%s", bm.Type, bm.Difficulty, bm.File)
}

// Song holds information about each song
type Song struct {
	Path  string
	key   string
	hash  string
	Name  string
	PP    float64
	Stars float64
	raw map[string]interface{}
	Maps []Beatmap
}

// CalcHash calculates this song's hash using its Path
func (s *Song) CalcHash() {
	// SHA1 hash of info.dat + all difficulties
	// f, err := os.Open(s.Path + "/info.dat")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// h := sha1.New()
	// if _, err := io.Copy(h, f); err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// f.Close()
	file, err := ioutil.ReadFile(s.Path + "/info.dat")
	if err != nil {
		fmt.Println(err)
		return
	}
	h := sha1.New()
	h.Write(file)

	s.hash = fmt.Sprintf("%x", h.Sum(nil))
}

// ParseRaw parses this song's info.dat
func (s *Song) ParseRaw() {
	file, err := ioutil.ReadFile(s.Path + "/info.dat")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(file, &s.raw)
	if err != nil {
		return
	}
	var bm Beatmap
	for _, sets := range s.raw["_difficultyBeatmapSets"].([]interface{}) {
		for k, v := range sets.(map[string]interface{}) {
			switch k {
			case "_beatmapCharacteristicName":
				bm.Type = v.(string)
			case "_difficultyBeatmaps":
				for _, diff := range v.([]interface{}) {
					for kk, vv := range diff.(map[string]interface{}) {
						switch kk {
						case "_difficulty":
							bm.Difficulty = vv.(string)
						case "_beatmapFilename":
							bm.File = s.Path + "/" + vv.(string)
						}
					}
				}
			default:
				continue
			}
		}
		s.Maps = append(s.Maps, bm)
		bm = Beatmap{}
	}
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
