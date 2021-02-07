package main

import "testing"

func TestDownloadSongInfo(t *testing.T) {
	s := Song{Hash: "9bf202f68c333421c69ca6aa15c648d65d4a1e0f", Name: "Night Raid"}
	out, err := DownloadSongInfo(&s)
	if err != nil {
		t.Errorf("Song info download failed: %v", err)
	} else {
		t.Logf("Song info download successful\n%s", out.Debug())
	}
}

func TestDownloadSong(t *testing.T) {
	s := Song{Hash: "9bf202f68c333421c69ca6aa15c648d65d4a1e0f", Name: "Night Raid"}
	out, err := DownloadSong(&s)
	if err != nil {
		t.Errorf("Song download failed: %v", err)
	} else {
		t.Logf("Song download successful\n%s", out.Debug())
	}
}
