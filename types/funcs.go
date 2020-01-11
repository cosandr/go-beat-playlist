package types

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"strings"
)

// MakePlaylist returns a Playlist from a json file path
func MakePlaylist(path string) (p Playlist, err error) {
	var j PlaylistJSON
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	err = json.Unmarshal(file, &j)
	if err != nil {
		return
	}
	var songs []Song
	for _, s := range j.Songs {
		key := ""
		// key might be read as float64 instead of string
		switch vv := s.Key.(type) {
		case float64:
			key = fmt.Sprintf("%.0f", vv)
		case string:
			key = vv
		}
		songs = append(songs, Song{
			Key:  strings.ToLower(key),
			Name: s.Name,
			Hash: strings.ToLower(s.Hash),
		})
	}
	p = Playlist{
		Title: j.Title,
		Author: j.Author,
		Image: j.Image,
		File: path,
		Songs: songs,
	}
	return
}

// MakeSong returns a Song from a info.dat file path
func MakeSong(path string) (s Song, err error) {
	var j InfoJSON
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	err = json.Unmarshal(file, &j)
	if err != nil {
		return
	}
	var maps []Beatmap
	for _, set := range j.Beatmaps {
		for _, m := range set.Maps {
			bm := Beatmap{Type: set.Type}
			bm.File = m.File
			bm.Difficulty = m.Difficulty
			maps = append(maps, bm)
		}
	}
	s = Song{
		Path: strings.TrimSuffix(path, "info.dat"),
		Name: j.SongName,
		Author: j.SongAuthor,
		Mapper: j.Mapper,
		Maps: maps,
	}
	s.CalcHash()
	return
}

// NewPath does nothing on Windows, replaces C: with /mnt/c and all \ with / on Linux
func NewPath(path string) string {
	ret := path
	if runtime.GOOS == "linux" {
		re := regexp.MustCompile(`(\w):\\`)
		m := re.FindStringSubmatchIndex(ret)
		if len(m) == 4 {
			// Find index of C:\ and replace it with /mnt/c/
			// Works for drive other letters
			ret = fmt.Sprintf("/mnt/%s/%s", strings.ToLower(ret[m[2]:m[3]]), ret[m[1]:])
		}
		// Strip all forward slashes
		ret = strings.ReplaceAll(ret, "\\", "/")
	}
	return ret
}

// FileExists returns true if `path` exists and is a file
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// DirExists returns true if `path` exists and is a directory
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
