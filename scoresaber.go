package main

import (
	"encoding/json"
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
	// Keep track of added hashes
	var songSet = make(StringSet)
	var empty struct{}
	for _, s := range resp.Songs {
		if songSet.Contains(s.ID) {
			continue
		}
		songSet[s.ID] = empty
		songs = append(songs, s.ToInternal())
	}
	p = Playlist{
		Title: "ScoreSaber Response",
		Songs: songs,
	}
	return
}
