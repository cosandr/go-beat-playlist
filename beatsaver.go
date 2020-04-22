package main

import (
	"encoding/json"
	"strings"
)

// BeatSaverSong is a BeatSaver song
type BeatSaverSong struct {
	Metadata BeatSaverMeta `json:"metadata"`
	Key      string        `json:"key"`
	Hash     string        `json:"hash"`
	URL      string        `json:"downloadURL"`
}

// BeatSaverMeta is the metadata from the BeatSaver API
type BeatSaverMeta struct {
	Chars  []BeatSaverMetaChar `json:"characteristics"`
	Mapper string              `json:"levelAuthorName"`
	Author string              `json:"songAuthorName"`
	Name   string              `json:"songName"`
}

// BeatSaverMetaChar is a metadata characteristic
type BeatSaverMetaChar struct {
	Name  string                 `json:"name"`
	Diffs map[string]interface{} `json:"difficulties"`
}

// MakeBeatSaverPlaylist returns a Playlist from a byte array (API response data)
func MakeBeatSaverPlaylist(file *[]byte) (p Playlist, err error) {
	var resp []BeatSaverSong
	err = json.Unmarshal(*file, &resp)
	if err != nil {
		return
	}
	var songs []Song
	for _, r := range resp {
		s := Song{
			Name:   r.Metadata.Name,
			Author: r.Metadata.Author,
			Key:    strings.ToLower(r.Key),
			Hash:   strings.ToLower(r.Hash),
			Mapper: r.Metadata.Mapper,
			URL:    r.URL,
		}
		maps := []Beatmap{}
		for _, diff := range r.Metadata.Chars {
			for k, v := range diff.Diffs {
				if v != nil {
					maps = append(maps, Beatmap{Type: diff.Name, Difficulty: k})
				}
			}
		}
		s.Maps = maps
		songs = append(songs, s)
	}
	p = Playlist{
		Title: "BeatSaver Response",
		Songs: songs,
	}
	return
}

// MakeBeatSaverSong returns a Song from a byte array (API response data)
func MakeBeatSaverSong(file *[]byte) (s Song, err error) {
	var resp BeatSaverSong
	err = json.Unmarshal(*file, &resp)
	if err != nil {
		return
	}
	maps := []Beatmap{}
	for _, diff := range resp.Metadata.Chars {
		for k, v := range diff.Diffs {
			if v != nil {
				maps = append(maps, Beatmap{Type: diff.Name, Difficulty: k})
			}
		}
	}
	s = Song{
		Name:   resp.Metadata.Name,
		Author: resp.Metadata.Author,
		Key:    strings.ToLower(resp.Key),
		Hash:   strings.ToLower(resp.Hash),
		Mapper: resp.Metadata.Mapper,
		URL:    resp.URL,
		Maps:   maps,
	}
	return
}
