package types

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ScoreSaberResp represents a list of songs in ScoreSaber's API response
type ScoreSaberResp struct {
	Songs []ScoreSaberSong `json:"songs"`
}

// ScoreSaberSong represents a song in ScoreSaber's API response
type ScoreSaberSong struct {
	UID    int     `json:"uid"`
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Author string  `json:"songAuthorName"`
	Mapper string  `json:"levelAuthorName"`
	Stars  float64 `json:"stars"`
}

// ToInternal returns a Song from this API response
func (s *ScoreSaberSong) ToInternal() Song {
	return Song{
		Key:    fmt.Sprintf("%d", s.UID),
		Hash:   strings.ToLower(s.ID),
		Name:   s.Name,
		Author: s.Author,
		Mapper: s.Mapper,
		Stars:  s.Stars,
	}
}

// MakeScoreSaberPlaylist returns a Playlist from a byte array (API response data)
func MakeScoreSaberPlaylist(file *[]byte) (p Playlist, err error) {
	var resp ScoreSaberResp
	err = json.Unmarshal(*file, &resp)
	if err != nil {
		return
	}
	var songs []Song
	for _, s := range resp.Songs {
		songs = append(songs, s.ToInternal())
	}
	p = Playlist{
		Title: "ScoreSaber Response",
		Songs: songs,
	}
	return
}
