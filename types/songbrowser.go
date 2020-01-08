package types

import (
	"encoding/json"
	"strconv"
	"strings"
)

// SongBrowserSong represents a song in Beat Saber Song Browser's API response
type SongBrowserSong struct {
	Diffs  []SongBrowserDiff `json:"diffs"`
	Key    string            `json:"key"`
	Mapper string            `json:"mapper"`
	Name   string            `json:"song"`
}

// SongBrowserDiff represents a song's difficulty entries in Beat Saber Song Browser's API response
type SongBrowserDiff struct {
	PP     string `json:"pp"`
	Star   string `json:"star"`
	Scores string `json:"scores"`
	Diff   string `json:"diff"`
}

// MakeSongBrowserPlaylist returns a Playlist from a byte array (API response data)
func MakeSongBrowserPlaylist(file *[]byte) (p Playlist, err error) {
	var resp map[string]SongBrowserSong
	err = json.Unmarshal(*file, &resp)
	if err != nil {
		return
	}
	var songs []Song
	for k, v := range resp {
		pp, _ := strconv.ParseFloat(v.Diffs[0].PP, 64)
		stars, _ := strconv.ParseFloat(v.Diffs[0].Star, 64)
		s := Song{
			Name:   v.Name,
			Key:    strings.ToLower(v.Key),
			Hash:   strings.ToLower(k),
			Mapper: v.Mapper,
			PP:     pp,
			Stars:  stars,
		}
		maps := []Beatmap{}
		for _, diff := range v.Diffs {
			maps = append(maps, Beatmap{Type: "Standard", Difficulty: diff.Diff})
		}
		s.Maps = maps
		songs = append(songs, s)
	}
	p = Playlist{
		Title: "SongBrowser Response",
		Songs: songs,
	}
	return
}
